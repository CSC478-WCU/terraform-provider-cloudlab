package experiment

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/csc478-wcu/terraform-provider-cloudlab/internal/portalclient"
)

// deepReady: overall "ready", all aggregates "ready", and every node has IPv4.
func deepReady(p *portalclient.StatusPayload) bool {
	if p == nil {
		return false
	}
	if canon(p.Status) == StatusReady {
		return true
	}
	for _, agg := range p.AggregateStatus {
		if canon(agg.Status) != StatusReady {
			return false
		}
	}
	nodes := portalclient.FlattenNodes(p)
	if len(nodes) == 0 {
		return false
	}
	for _, n := range nodes {
		if strings.TrimSpace(n.IPv4) == "" {
			return false
		}
	}
	return true
}

// Predicate:
// - waitFor == "ready"  -> require deepReady() (every node ready)
// - else                -> succeed when current rank >= target rank
func Predicate(ctx context.Context, waitFor string) func(*portalclient.StatusPayload) bool {
	target := rankOf(waitFor)
	wantReady := canon(waitFor) == StatusReady

	return func(p *portalclient.StatusPayload) bool {
		curStatus := "<nil>"
		curRank := 0
		if p != nil {
			curStatus = p.Status
			curRank = rankOf(p.Status)
		}
		tflog.Debug(ctx, "wait tick", map[string]any{
			"target": waitFor, "target_rank": target,
			"current": curStatus, "current_rank": curRank,
		})

		if p == nil {
			return false
		}
		if wantReady {
			return deepReady(p)
		}
		return curRank >= target
	}
}
