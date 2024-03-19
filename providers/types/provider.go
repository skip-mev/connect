package types

import (
	"golang.org/x/net/context"
)

type ProviderType string

const (
	WebSockets ProviderType = "websockets"
	API        ProviderType = "api"
)

// Provider defines an interface a data provider must implement.
//
//go:generate mockery --name Provider --filename mock_provider.go
type Provider[K ResponseKey, V ResponseValue] interface {
	// Name returns the name of the provider.
	Name() string

	// GetData returns the aggregated data for the given (key, value) pairs.
	// For example, if the provider is fetching prices for a set of currency
	// pairs, the data returned by this function would be the latest prices
	// for those currency pairs.
	GetData() map[K]ResolvedResult[V]

	// Start starts the provider.
	Start(context.Context) error

	// Type returns the type of the provider data handler.
	Type() ProviderType

	// IsRunning returns whether the provider is running.
	IsRunning() bool
}
