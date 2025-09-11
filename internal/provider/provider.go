package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/csc478-wcu/terraform-provider-cloudlab/internal/experiment"
	"github.com/csc478-wcu/terraform-provider-cloudlab/internal/portalclient"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"project":  {Type: schema.TypeString, Optional: true},
			"pem_path": {Type: schema.TypeString, Optional: true, Default: "~/cloudlab.pem"},
			"server":   {Type: schema.TypeString, Optional: true, Default: "boss.emulab.net"},
			"port":     {Type: schema.TypeInt, Optional: true, Default: 3069},
			"path":     {Type: schema.TypeString, Optional: true, Default: "/usr/testbed"},
			"timeout":  {Type: schema.TypeString, Optional: true, Default: "10m"},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cloudlab_portal_experiment": experiment.Resource(),
		},
	}

	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		to, _ := time.ParseDuration(d.Get("timeout").(string))
		opts := portalclient.Options{
			Server:  d.Get("server").(string),
			Port:    d.Get("port").(int),
			Path:    d.Get("path").(string),
			CertPEM: d.Get("pem_path").(string),
			KeyPEM:  d.Get("pem_path").(string),
			Verify:  false, // always self-signed
			Timeout: to,
		}
		cli, err := portalclient.New(opts)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		cfg := &portalclient.Config{
			Client:  cli,
			Project: d.Get("project").(string),
			PemPath: d.Get("pem_path").(string),
		}
		return cfg, nil
	}

	return p
}
