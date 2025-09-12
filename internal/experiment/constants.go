// internal/experiment/status.go
package experiment

import "strings"

const (
	profileName          = "cloud-edu,terraform-profile"
	profileParamSpecJSON = "spec_json"

	StatusProvisioning = "provisioning"
	StatusProvisioned  = "provisioned"
	StatusCreating     = "creating"
	StatusCreated      = "created"
	StatusBooting      = "booting"
	StatusBooted       = "booted"
	StatusReady        = "ready"
)

var statusOrder = []string{
	StatusProvisioning,
	StatusProvisioned,
	StatusCreating,
	StatusCreated,
	StatusBooting,
	StatusBooted,
	StatusReady,
}

// fast lookup: status -> rank
var statusRank = func() map[string]int {
	m := make(map[string]int, len(statusOrder)+8)
	for i, s := range statusOrder {
		m[strings.ToLower(s)] = i + 1
	}
	// common site/output variants & synonyms
	m["starting"] = m[StatusBooting]
	m["started"] = m[StatusBooted]
	m["running"] = m[StatusBooted]
	m["up"] = m[StatusBooted]
	m["booted/provisioned"] = m[StatusBooted]
	// tolerate old terms if they appear
	m["swapped"] = m[StatusProvisioned]
	m["swapped-in"] = m[StatusProvisioned]
	return m
}()

func canon(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

func rankOf(s string) int {
	if r, ok := statusRank[canon(s)]; ok {
		return r
	}
	return 0
}
