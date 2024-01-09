package handlers

import (
	"context"

	providertypes "github.com/skip-mev/slinky/providers/types"
)

// QueryHandler is an interface that encapsulates querying a data provider for info.
// The handler must respect the context timeout and cancel the request if the context
// is cancelled. All responses must be sent to the response channel. These are processed
// asynchronously by the provider.
//
//go:generate mockery --name QueryHandler --output ./mocks/ --case underscore
type QueryHandler[K comparable, V any] interface {
	Query(
		ctx context.Context,
		ids []K,
		responseCh chan<- providertypes.GetResponse[K, V],
	)
}
