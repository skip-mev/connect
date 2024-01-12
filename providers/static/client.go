package static

import (
	"context"
	"net/http"

	"github.com/skip-mev/slinky/providers/base/api/handlers"
)

var _ handlers.RequestHandler = (*StaticMockClient)(nil)

// StaticMockClient is meant to be paired with the StaticMockAPIHandler. It
// should only be used for testing.
type StaticMockClient struct{} //nolint

func NewStaticMockClient() *StaticMockClient {
	return &StaticMockClient{}
}

// Do is a no-op.
func (s *StaticMockClient) Do(_ context.Context, _ string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
	}, nil
}
