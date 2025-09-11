package portalclient

import (
	"context"
	"time"

	portal "github.com/csc478-wcu/portalctl/portal"
)

type Client = portal.Client
type StatusPayload = portal.StatusPayload
type EmulabResponse = portal.EmulabResponse

func New(o Options) (*portal.Client, error) {
	return portal.New(portal.Options{
		Server:  o.Server,
		Port:    o.Port,
		Path:    o.Path,
		CertPEM: o.CertPEM,
		KeyPEM:  o.KeyPEM,
		Verify:  o.Verify,  // =false
		Timeout: o.Timeout,
	})
}

func StartExperiment(c *portal.Client, params map[string]any) (*portal.EmulabResponse, error) {
	return c.StartExperiment(params)
}

func Terminate(c *portal.Client, exp string) (*portal.EmulabResponse, error) {
	return c.TerminateExperiment(exp)
}

func Status(c *portal.Client, exp string, asJSON, withCert, refresh bool) (*portal.EmulabResponse, error) {
	return c.ExperimentStatus(exp, asJSON, withCert, refresh)
}

func Wait(ctx context.Context, c *portal.Client, exp string, interval, timeout time.Duration, done func(*portal.StatusPayload) bool) (*portal.StatusPayload, error) {
	return c.WaitForStatus(ctx, exp, interval, timeout, done)
}

func ParseStatusJSON(s string) (*portal.StatusPayload, error) { return portal.ParseStatusJSON(s) }
func FlattenNodes(p *portal.StatusPayload) map[string]portal.StatusNode {
	return portal.FlattenNodes(p)
}
