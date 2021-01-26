package testutilities

import "net/http"

// NewTestClient returns a new HTTP Client utilizing a custom Transport specification
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{Transport: RoundTripFunc(fn)}
}

// RoundTripFunc acts as a function type that implements the RoundTripper interface
// allowing us to inject our own response
type RoundTripFunc func(req *http.Request) (*http.Response, error)

// RoundTrip is a method that overrides RoundTrip's default behavior while maintining the
// same signature needed for an HTTP Client's Transport specification
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
