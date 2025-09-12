package portalclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	portal "github.com/csc478-wcu/portalctl/portal"
)

// Re-export types for provider packages.
type Client = portal.Client
type StatusPayload = portal.StatusPayload
type EmulabResponse = portal.EmulabResponse

// New returns a configured XML-RPC client.
func New(o Options) (*portal.Client, error) {
	return portal.New(portal.Options{
		Server:  o.Server,
		Port:    o.Port,
		Path:    o.Path,
		CertPEM: o.CertPEM,
		KeyPEM:  o.KeyPEM,
		Verify:  o.Verify,
		Timeout: o.Timeout,
	})
}

// ----- High-level helpers (provider-friendly) -----
// StartExperiment is pass-through; your params already contain project/name.
func StartExperiment(c *portal.Client, params map[string]any) (*portal.EmulabResponse, error) {
	return c.StartExperiment(params)
}

// Status calls portal.experimentStatus with "project,exp" combined into
// the XML-RPC "experiment" parameter, since the backend accepts comma form.
func Status(c *portal.Client, project, exp string, asJSON, withCert, refresh bool) (*portal.EmulabResponse, error) {
	combined := fmt.Sprintf("%s,%s", strings.TrimSpace(project), strings.TrimSpace(exp))
	return c.ExperimentStatus(combined, asJSON, withCert, refresh)
}

// Terminate invokes portal.terminateExperiment with "project,exp".
func Terminate(c *portal.Client, project, exp string) (*portal.EmulabResponse, error) {
	combined := fmt.Sprintf("%s,%s", strings.TrimSpace(project), strings.TrimSpace(exp))
	return c.TerminateExperiment(combined)
}

// Manifests invokes portal.experimentManifests with "project,exp".
func Manifests(c *portal.Client, project, exp string) (*portal.EmulabResponse, error) {
	combined := fmt.Sprintf("%s,%s", strings.TrimSpace(project), strings.TrimSpace(exp))
	return c.ExperimentManifests(combined)
}

// Wait delegates to your portal client’s waiter (if you keep it);
// otherwise prefer helper/retry.StateChangeConf in the resource.
func Wait(ctx context.Context, c *portal.Client, exp string, interval, timeout time.Duration, done func(*portal.StatusPayload) bool) (*portal.StatusPayload, error) {
	return c.WaitForStatus(ctx, exp, interval, timeout, done)
}

// Parse/Flatten re-exports.
func ParseStatusJSON(s string) (*portal.StatusPayload, error) { return portal.ParseStatusJSON(s) }
func FlattenNodes(p *portal.StatusPayload) map[string]portal.StatusNode {
	return portal.FlattenNodes(p)
}
