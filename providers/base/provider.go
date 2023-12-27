package base

import (
	"context"
	"fmt"
	"maps"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
)

// BaseProvider implements a base provider that can be used to build other providers.
type BaseProvider[K comparable, V any] struct { //nolint
	mu     sync.Mutex
	logger *zap.Logger

	// APIDataHandler is the handler for the API data. Developer's implement this interface
	// to extend the provider's functionality. For example, this could be used to fetch
	// prices from an API, where K is the currency pair and V is the price.
	handler APIDataHandler[K, V]

	// cfg is the provider's config. This contains the name, path, fetch timeout, and
	// fetch interval for the provider. To read more about the config, see the
	// oracle/config/provider.go file. To read more about how to configure a custom provider,
	// please see the providers/README.md file.
	cfg config.ProviderConfig

	// data is the latest set of key -> value pairs for the provider i.e. the latest prices
	// for a given set of currency pairs.
	data map[K]V

	// lastUpdated is the time at which the data was last updated/fetched.
	lastUpdate time.Time
}

// NewProvider returns a new Base provider.
func NewProvider[K comparable, V any](
	logger *zap.Logger,
	cfg config.ProviderConfig,
	handler APIDataHandler[K, V],
) (*BaseProvider[K, V], error) {
	if handler == nil {
		return nil, fmt.Errorf("no api data handler specified for provider %s", cfg.Name)
	}

	return &BaseProvider[K, V]{
		logger:  logger.With(zap.String("provider", cfg.Name)),
		cfg:     cfg,
		handler: handler,
		data:    make(map[K]V),
	}, nil
}

// Start starts the provider's main loop. The provider will fetch the data from the API
// and continuously update the data. This blocks until the provider is stopped.
func (p *BaseProvider[K, V]) Start(ctx context.Context) error {
	p.logger.Info("starting provider")

	// Start the main loop.
	return p.loop(ctx)
}

// Name returns the name of the provider.
func (p *BaseProvider[K, V]) Name() string {
	return p.cfg.Name
}

// GetData returns the latest data recorded by the provider. The data is constantly
// updated by the provider's main loop and provides access to the latest data - prices
// in constant time.
func (p *BaseProvider[K, V]) GetData() map[K]V {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Deep copy the prices into a new map.
	cpy := make(map[K]V)
	maps.Copy(cpy, p.data)

	return cpy
}

// LastUpdate returns the time at which the prices were last updated.
func (p *BaseProvider[K, V]) LastUpdate() time.Time {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.lastUpdate
}
