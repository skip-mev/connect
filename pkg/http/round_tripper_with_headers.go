package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/skip-mev/connect/v2/cmd/build"
)

const (
	// UserAgentHeaderKey is the key for the User-Agent header.
	UserAgentHeaderKey = "User-Agent"
)

// HeaderOption is a function that sets a (key, value) pair in the headers map.
type HeaderOption func(map[string]string)

func WithConnectVersionUserAgent() HeaderOption {
	return func(header map[string]string) {
		header[UserAgentHeaderKey] = fmt.Sprintf("connect/%s", build.Build)
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

// Client is a wrapper around the Go stdlib http client.
type Client struct {
	internal *http.Client
}

// NewClient returns a new Client with its internal http client
// set to the default client.
func NewClient() *Client {
	return &Client{
		internal: http.DefaultClient,
	}
}

type GetOptions func(*http.Request)

func WithHeader(key, value string) GetOptions {
	return func(r *http.Request) {
		r.Header.Add(key, value)
	}
}

func WithJSONAccept() GetOptions {
	return func(r *http.Request) {
		r.Header.Add("Accept", "application/json")
	}
}

// GetWithContext performs a Get request with the context provided.
func (c *Client) GetWithContext(ctx context.Context, url string, opts ...GetOptions) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.internal.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	return resp, nil
}
