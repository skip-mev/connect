package http

import (
	"net/http"

	"github.com/skip-mev/slinky/cmd/build"
)

const (
	// UserAgentHeaderKey is the key for the User-Agent header.
	UserAgentHeaderKey = "User-Agent"
)

// HeaderOption is a function that sets a (key, value) pair in the headers map.
type HeaderOption func(map[string]string)

func WithSlinkyVersionUserAgent() HeaderOption {
	return func(header map[string]string) {
		header[UserAgentHeaderKey] = build.Build
	}
}

func WithAuthentication(key, token string) HeaderOption {
	return func(header map[string]string) {
		header[key] = token
	}
}

// RoundTripperWithHeaders is a round tripper that adds headers to the request.
type RoundTripperWithHeaders struct {
	// Headers is the map of headers to add to the request.
	headers map[string]string

	// Next is the next round tripper in the chain.
	next http.RoundTripper
}

// NewRoundTripperWithHeaders creates a new RoundTripperWithHeaders.
func NewRoundTripperWithHeaders(next http.RoundTripper, headerOptions ...HeaderOption) *RoundTripperWithHeaders {
	headers := make(map[string]string)

	for _, option := range headerOptions {
		option(headers)
	}

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
