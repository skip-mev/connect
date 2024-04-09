package http

import (
	"net/http"
)

// RoundTripperWithHeaders is a round tripper that adds headers to the request.
type RoundTripperWithHeaders struct {
	// Headers is the map of headers to add to the request.
	headers map[string]string

	// Next is the next round tripper in the chain.
	next http.RoundTripper
}

// NewRoundTripperWithHeaders creates a new RoundTripperWithHeaders.
func NewRoundTripperWithHeaders(headers map[string]string, next http.RoundTripper) *RoundTripperWithHeaders {
	return &RoundTripperWithHeaders{
		headers: headers,
		next:    next,
	}
}

// RoundTrip updates the Requests' headers with the headers specified in the constructor, and calls the underlying RoundTripper.
func (r *RoundTripperWithHeaders) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	return r.next.RoundTrip(req)
}
