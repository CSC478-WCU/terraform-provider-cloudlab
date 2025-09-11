package experiment

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/csc478-wcu/terraform-provider-cloudlab/internal/portalclient"
)

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*portalclient.Config)

	spec := buildSpec(d)
	if err := validateSpec(spec); err != nil {
		return diag.FromErr(err)
	}
	specJSON, err := encodeSpec(spec)
	if err != nil {
		return diag.FromErr(err)
	}

	expName := d.Get("name").(string)
	project := d.Get("project").(string)
	if project == "" {
		project = cfg.Project
	}

	params := composeParams(project, expName, specJSON)
	if _, err := portalclient.StartExperiment(cfg.Client, params); err != nil {
		return diag.FromErr(err)
	}

	waitFor := d.Get("wait_for_status").(string)
	pr := predicate(waitFor)

	to := d.Timeout(schema.TimeoutCreate)
	if to == 0 {
		to = 30 * time.Minute
	}
	ctxWait, cancel := context.WithTimeout(ctx, to)
	defer cancel()

	_, err = portalclient.Wait(ctxWait, cfg.Client, expName, 10*time.Second, to, pr)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(expName)
	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*portalclient.Config)
	exp := d.Id()
	if exp == "" {
		return nil
	}

	resp, err := portalclient.Status(cfg.Client, exp, true, false, true)
	if err != nil {
		// if not found, clear state (best effort)
		d.SetId("")
		return nil
	}
	p, err := portalclient.ParseStatusJSON(resp.Output)
	if err != nil {
		return diag.FromErr(err)
	}
	setStatusFields(d, p)
	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*portalclient.Config)
	exp := d.Id()
	if exp == "" {
		return nil
	}
	if _, err := portalclient.Terminate(cfg.Client, exp); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
