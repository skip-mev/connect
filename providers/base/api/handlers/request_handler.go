package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// RequestHandler is an interface that encapsulates sending an HTTP request to a data provider.
//
//go:generate mockery --name RequestHandler --output ./mocks/ --case underscore
type RequestHandler interface {
	// Do is used to send a request with the given URL to the data provider.
	Do(ctx context.Context, url string, body io.Reader) (*http.Response, error)
}

var _ RequestHandler = (*RequestHandlerImpl)(nil)

// HTTPClient is the interface expected by the golang upstream HTTP library.
//
//go:generate mockery --name HTTPClient --output ./mocks/ --case underscore
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RequestHandlerImpl is the default implementation of the RequestHandler interface.
type RequestHandlerImpl struct {
	// client is the HTTP client to use when sending requests.
	client HTTPClient

	// method is the HTTP method to use when sending requests.
	method string

	// request headers is a map of HTTP request header (key, value)
	// pairs to transmit along with the request.
	requestHeaderPairs map[string]string
}

// NewRequestHandlerImpl creates a new RequestHandlerImpl. It manages making HTTP requests.
func NewRequestHandlerImpl(client HTTPClient, opts ...Option) (RequestHandler, error) {
	h := &RequestHandlerImpl{
		client: client,
		method: http.MethodGet,
	}

	for _, opt := range opts {
		opt(h)
	}

	if h.method == "" {
		return nil, fmt.Errorf("http request method cannot be empty")
	}

	return h, nil
}

// Do is used to send a request with the given URL to the data provider. It first
// wraps the request with the given context before sending it to the data provider.
func (r *RequestHandlerImpl) Do(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, r.method, url, body)
	if err != nil {
		return nil, err
	}

	// set request headers
	for key, value := range r.requestHeaderPairs {
		req.Header.Set(key, value)
	}

	return r.client.Do(req)
}
