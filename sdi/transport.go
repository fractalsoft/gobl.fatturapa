package sdi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
)

// Transport defines interface for http client transport methods
type Transport interface {
	SetTransport(t *http.Transport)
	WithContext(ctx context.Context)
	RoundTrip(r *http.Request) (*http.Response, error)
}

// DefaultTransport default configuration for HTTP communications
type DefaultTransport struct {
	HTTPTransport *http.Transport
	Context       context.Context
}

// SetTransport sets the underlying HTTPTransport value
func (s *DefaultTransport) SetTransport(t *http.Transport) {
	s.HTTPTransport = t
}

// WithContext adds context to the Transport
func (s *DefaultTransport) WithContext(ctx context.Context) {
	s.Context = ctx
}

// RoundTrip executes a single HTTP transaction, returning,
// a Response for the provided Request.
func (s *DefaultTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return s.HTTPTransport.RoundTrip(r.WithContext(s.Context))
}

// LoggingTransport allows us to debug HTTP communications
type LoggingTransport struct {
	HTTPTransport *http.Transport
	Context       context.Context
}

// SetTransport sets the underlying HTTPTransport value
func (s *LoggingTransport) SetTransport(t *http.Transport) {
	s.HTTPTransport = t
}

// WithContext adds context to the Transport
func (s *LoggingTransport) WithContext(ctx context.Context) {
	s.Context = ctx
}

// RoundTrip executes a single HTTP transaction, returning,
// a Response for the provided Request, but also logging that things for developer.
func (s *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.WithContext(s.Context)

	bytes, _ := httputil.DumpRequestOut(r, true)
	fmt.Printf("%s\n", bytes)

	resp, err := s.HTTPTransport.RoundTrip(r)
	// err is returned after dumping the response
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}

	respBytes, _ := httputil.DumpResponse(resp, true)
	bytes = append(bytes, respBytes...)

	fmt.Printf("%s\n", bytes)

	return resp, err
}
