package portalclient

import "time"

type Options struct {
	Server  string
	Port    int
	Path    string
	CertPEM string
	KeyPEM  string
	Verify  bool
	Timeout time.Duration
}

type Config struct {
	Client  *Client
	Project string
	PemPath string
}
