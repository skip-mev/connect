package handlers

import (
	"context"
	"fmt"
	"net/http"
)

// RequestHandler is an interface that encapsulates sending an HTTP request to a data provider.
//
//go:generate mockery --name RequestHandler --output ./mocks/ --case underscore
type RequestHandler interface {
	// Do is used to send a request with the given URL to the data provider.
	Do(ctx context.Context, url string) (*http.Response, error)

	// Type defines the type of the RequestHandler based on the type of
	// HTTP requests it makes  - GET, POST, etc.
	Type() string
}

var _ RequestHandler = (*RequestHandlerImpl)(nil)

// RequestHandlerImpl is the default implementation of the RequestHandler interface.
type RequestHandlerImpl struct {
	client *http.Client

	// method is the HTTP method to use when sending requests.
	method string
	// headers is the HTTP headers to use when sending requests.
	headers map[string]string
}

// NewRequestHandlerImpl creates a new RequestHandlerImpl. It manages making HTTP requests.
func NewRequestHandlerImpl(client *http.Client, opts ...Option) (RequestHandler, error) {
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
func (r *RequestHandlerImpl) Do(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, r.method, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	return r.client.Do(req)
}

// Type returns the HTTP method used to send requests.
func (r *RequestHandlerImpl) Type() string {
	return r.method
}
