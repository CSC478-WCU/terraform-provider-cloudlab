package experiment

import (
	"github.com/csc478-wcu/terraform-provider-cloudlab/internal/portalclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setStatusFields(d *schema.ResourceData, p *portalclient.StatusPayload) {
	_ = d.Set("status", p.Status)
	_ = d.Set("uuid", p.UUID)
	_ = d.Set("url", p.URL)
	_ = d.Set("expires", p.Expires)

	nodes := portalclient.FlattenNodes(p)
	out := map[string]string{}
	for id, n := range nodes {
		out[id] = n.IPv4
	}
	_ = d.Set("nodes", out)
}
