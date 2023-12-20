package base

import (
	"context"
)

// APIDataHandler interface defines the methods that need to be implemented by the extender.
//
//go:generate mockery --name APIDataHandler --output ./mocks/ --case underscore
type APIDataHandler[K comparable, V any] interface {
	// Get is used to fetch data from the API.
	Get(ctx context.Context) (map[K]V, error)
}
