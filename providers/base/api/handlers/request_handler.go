package handlers

import (
	"context"
	"net/http"
)

// RequestHandler is an interface that encapsulates sending an HTTP request to a data provider.
//
//go:generate mockery --name RequestHandler --output ./mocks/ --case underscore
type RequestHandler interface {
	// Do is used to send a request with the given URL to the data provider.
	Do(ctx context.Context, url string) (*http.Response, error)
}

var _ RequestHandler = (*RequestHandlerImpl)(nil)

// RequestHandlerImpl is the default implementation of the RequestHandler interface.
type RequestHandlerImpl struct {
	client *http.Client
}

// NewRequestHandlerImpl creates a new RequestHandlerImpl. It manages making HTTP requests.
func NewRequestHandlerImpl(client *http.Client) RequestHandler {
	return &RequestHandlerImpl{
		client: client,
	}
}

// Do is used to send a request with the given URL to the data provider. It first
// wraps the request with the given context before sending it to the data provider.
func (r *RequestHandlerImpl) Do(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return r.client.Do(req)
}
