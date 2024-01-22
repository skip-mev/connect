package types

import (
	"go.uber.org/zap"
	"golang.org/x/net/context"

	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"

	"github.com/skip-mev/slinky/oracle/config"
)

// Provider defines an interface a data provider must implement.
//
//go:generate mockery --name Provider --filename mock_provider.go
type Provider[K comparable, V GetResult] interface {
	// Name returns the name of the provider.
	Name() string

	// GetData returns the aggregated data for the given (key, value) pairs.
	// For example, if the provider is fetching prices for a set of currency
	// pairs, the data returned by this function would be the latest prices
	// for those currency pairs.
	GetData() map[K]Result[V]

	// Start starts the provider.
	Start(context.Context) error

	// Type returns the type of the provider data handler.
	Type() providermetrics.ProviderType
}

// ProviderFactory inputs the oracle configuration and returns a set of providers. Developers
// can implement their own provider factory to create their own providers.
type ProviderFactory[K comparable, V GetResult] func(
	*zap.Logger,
	config.OracleConfig,
	config.OracleMetricsConfig,
) ([]Provider[K, V], error)
