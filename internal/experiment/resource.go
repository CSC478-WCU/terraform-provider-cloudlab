package experiment

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/csc478-wcu/terraform-provider-cloudlab/internal/portalclient"
)

// Starts the experiment and waits until the requested status.
func resourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*portalclient.Config)

	spec := buildSpec(d)
	if err := validateSpec(spec); err != nil { return diag.FromErr(err) }
	specJSON, err := encodeSpec(spec)
	if err != nil { return diag.FromErr(err) }

	expName := d.Get("name").(string)
	project := d.Get("project").(string)
	if project == "" { project = cfg.Project }

	params := composeParams(project, expName, specJSON)
	tflog.Info(ctx, "starting experiment", map[string]any{"project": project, "experiment": expName})
	if _, err := portalclient.StartExperiment(cfg.Client, params); err != nil {
		return diag.FromErr(err)
	}

	waitFor := canon(d.Get("wait_for_status").(string))
	pred := Predicate(ctx, waitFor)

	// Timeouts/polling
	to := d.Timeout(schema.TimeoutCreate)
	if to == 0 { to = 30 * time.Minute }
	poll := 10 * time.Second
	warmup := 15 * time.Second // allow control plane to register

	stateConf := &retry.StateChangeConf{
		Pending: []string{
			StatusProvisioning, StatusProvisioned,
			StatusCreating, StatusCreated,
			StatusBooting, StatusBooted,
		},
		Target:     []string{waitFor}, // dynamic: exactly what user asked for
		Timeout:    to,
		Delay:      warmup,
		MinTimeout: poll,
		Refresh: func() (interface{}, string, error) {
			resp, err := portalclient.Status(cfg.Client, project, expName, true, false, true)
			if err != nil {
				tflog.Warn(ctx, "status fetch failed; retrying", map[string]any{"error": err})
				return nil, "", nil
			}

			// Always parse JSON out of mixed output (no lock handling needed).
			p, perr := portalclient.ParseStatusJSONLoose(resp.Output)
			if perr != nil {
				tflog.Warn(ctx, "bad status json; retrying", map[string]any{"error": perr})
				return nil, "", nil
			}

			tflog.Debug(ctx, "poll", map[string]any{
				"project": project, "experiment": expName,
				"status": p.Status, "aggs": len(p.AggregateStatus),
				"nodes": len(portalclient.FlattenNodes(p)),
			})

			if pred(p) {
				return p, waitFor, nil // success: return exactly the waited-for state
			}
			return p, p.Status, nil   // keep waiting
		},
	}

	ctxWait, cancel := context.WithTimeout(ctx, to)
	defer cancel()

	out, err := stateConf.WaitForStateContext(ctxWait)
	if err != nil {
		last := ""
		if p, ok := out.(*portalclient.StatusPayload); ok && p != nil {
			last = p.Status
		}
		return diag.Errorf("waiting for %q to reach %q failed (last=%q): %v",
			expName, waitFor, last, err)
	}

	d.SetId(expName)
	return resourceRead(ctx, d, meta)
}

// Reads experiment state into Terraform.
func resourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*portalclient.Config)

	expName := d.Id()
	if expName == "" {
		return nil
	}
	project := d.Get("project").(string)
	if project == "" {
		project = cfg.Project
	}

	resp, err := portalclient.Status(cfg.Client, project, expName, true, false, true)
	if err != nil {
		tflog.Warn(ctx, "experiment not found during read; clearing state",
			map[string]any{"project": project, "experiment": expName, "error": err})
		d.SetId("")
		return nil
	}
	p, err := portalclient.ParseStatusJSONLoose(resp.Output)
	if err != nil {
		return diag.FromErr(err)
	}
	setStatusFields(d, p)
	return nil
}

// Deletes the experiment.
func resourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*portalclient.Config)

	expName := d.Id()
	if expName == "" {
		return nil
	}
	project := d.Get("project").(string)
	if project == "" {
		project = cfg.Project
	}

	tflog.Info(ctx, "terminating experiment", map[string]any{
		"project": project, "experiment": expName,
	})
	if _, err := portalclient.Terminate(cfg.Client, project, expName); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
