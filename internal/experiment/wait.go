package experiment

import "github.com/csc478-wcu/terraform-provider-cloudlab/internal/portalclient"

func predicate(waitFor string) func(*portalclient.StatusPayload) bool {
	return func(p *portalclient.StatusPayload) bool {
		if p == nil {
			return false
		}
		switch waitFor {
		case "ready":
			return p.Status == StatusReady
		default: // "provisioned"
			return p.Status == StatusProvisioned || p.Status == StatusSwapped || p.Status == StatusReady
		}
	}
}
