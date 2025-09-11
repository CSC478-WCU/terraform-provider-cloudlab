// internal/experiment/wait.go
package experiment

import (
	"strings"

	"github.com/csc478-wcu/terraform-provider-cloudlab/internal/portalclient"
)

// 1) Canonical ordering (low → high). Adjust if your portal uses different steps.
var statusOrder = []string{
	StatusProvisioning, // "provisioning"
	StatusProvisioned,  // "provisioned"
	StatusSwapped,      // "swapped"
	StatusCreating,     // "creating"
	StatusCreated,      // "created"
	StatusStarting,     // "starting"
	StatusReady,        // "ready"
}

// 2) Fast lookup: status -> rank
var statusRankMap = func() map[string]int {
	m := make(map[string]int, len(statusOrder))
	for i, s := range statusOrder {
		m[strings.ToLower(s)] = i + 1 // ranks start at 1
	}
	// Synonyms / site variants mapped to their canonical
	m["swapped-in"] = m[strings.ToLower(StatusSwapped)]
	m["swappedin"]  = m[strings.ToLower(StatusSwapped)]
	m["swapin complete"] = m[strings.ToLower(StatusSwapped)]
	m["swapin_complete"] = m[strings.ToLower(StatusSwapped)]
	m["booted"] = m[strings.ToLower(StatusProvisioned)]
	m["booted/provisioned"] = m[strings.ToLower(StatusProvisioned)]
	return m
}()

func canon(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

func rankOf(s string) int {
	if r, ok := statusRankMap[canon(s)]; ok {
		return r
	}
	return 0 // unknown/unordered → treat as before the first state
}

// Optional: deeper "ready" test (aggregates ready + nodes have IPv4)
func deepReady(p *portalclient.StatusPayload) bool {
	if p == nil {
		return false
	}
	if canon(p.Status) == "ready" {
		return true
	}
	for _, agg := range p.AggregateStatus {
		if canon(agg.Status) != "ready" {
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

// Predicate: satisfied if current rank >= target rank.
// For waitFor="ready", also require deepReady().
func predicate(waitFor string) func(*portalclient.StatusPayload) bool {
	target := rankOf(waitFor)
	wantReady := canon(waitFor) == "ready"

	return func(p *portalclient.StatusPayload) bool {
		if p == nil {
			return false
		}
		if rankOf(p.Status) < target {
			return false
		}
		if wantReady {
			return deepReady(p)
		}
		return true
	}
}
