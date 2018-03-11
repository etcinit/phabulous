package core

import (
	"crypto/tls"
	"net/http"
)

// ClientOptions are options that can be set on the HTTP client.
type ClientOptions struct {
	APIToken string

	Cert       string
	CertUser   string
	SessionKey string

	InsecureSkipVerify bool
}

// makeHttpClient creates a new HTTP client for making API requests.
func makeHTTPClient(options *ClientOptions) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: options.InsecureSkipVerify,
			},
		},
	}
}
