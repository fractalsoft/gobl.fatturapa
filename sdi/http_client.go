package sdi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"
)

// HTTPClientOptFunc defines function for customizing the client
type HTTPClientOptFunc func(*HTTPClientOpts)

// HTTPClientOpts defines the client parameters
type HTTPClientOpts struct {
	CaCertPool *x509.CertPool
	Context    context.Context
	Transport  Transport
	Timeout    uint
}

func defaultHTTPClientOpts() HTTPClientOpts {
	return HTTPClientOpts{
		Transport:  &DefaultTransport{},
		Timeout:    2500,
		CaCertPool: nil,
		Context:    context.Background(),
	}
}

// HTTPClient defines HTTPClient used to make http requests
type HTTPClient struct {
	HTTPClientOpts
}

// NewHTTPClient returns HTTP client with Cert Pool
func NewHTTPClient(opts ...HTTPClientOptFunc) *HTTPClient {
	o := defaultHTTPClientOpts()
	for _, fn := range opts {
		fn(&o)
	}
	return &HTTPClient{
		HTTPClientOpts: o,
	}
}

// WithDebugClient uses a more verbose client
func WithDebugClient() HTTPClientOptFunc {
	return func(o *HTTPClientOpts) {
		o.Transport = &LoggingTransport{}
	}
}

// WithCaCertPool sets the certificate pool used in http requests
func WithCaCertPool(pool *x509.CertPool) HTTPClientOptFunc {
	return func(o *HTTPClientOpts) {
		o.CaCertPool = pool
	}
}

// WithTimeout sets the time in miliseconds after which the requests should timeout
func WithTimeout(ms uint) HTTPClientOptFunc {
	return func(o *HTTPClientOpts) {
		o.Timeout = ms
	}
}

// WithContext adds context to the HTTPClient
func WithContext(ctx context.Context) HTTPClientOptFunc {
	return func(o *HTTPClientOpts) {
		o.Context = ctx
	}
}

// Build a HTTP Client with set properties
func (c HTTPClient) Build() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: c.CaCertPool,
		},
	}

	c.Transport.SetTransport(transport)
	c.Transport.WithContext(c.Context)

	client := &http.Client{
		Timeout:   time.Duration(c.Timeout) * time.Millisecond,
		Transport: c.Transport,
	}

	return client
}
