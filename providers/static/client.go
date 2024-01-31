package static

import (
	"context"
	"net/http"

	"github.com/skip-mev/slinky/providers/base/api/handlers"
)

var _ handlers.RequestHandler = (*MockClient)(nil)

// MockClient is meant to be paired with the MockAPIHandler. It
// should only be used for testing.
type MockClient struct{}

func NewStaticMockClient() *MockClient {
	return &MockClient{}
}

// Do is a no-op.
func (s *MockClient) Do(_ context.Context, _ string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
	}, nil
}

// Type returns the HTTP method used to send requests.
func (s *MockClient) Type() string {
	return http.MethodGet
}
