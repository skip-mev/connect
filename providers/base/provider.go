package base

import (
	"context"
	"fmt"
	"maps"
	"sync"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// BaseProvider implements a base provider that can be used to build other providers.
type BaseProvider[K comparable, V any] struct { //nolint
	mu     sync.Mutex
	logger *zap.Logger

	// handler is the handler for the querying data. Developer's implement this interface
	// to extend the provider's functionality. For example, this could be used to fetch
	// prices from an API, where K is the currency pair and V is the price. For more information
	// on how to implement a custom handler, please see the providers/base/README.md file.
	handler handlers.QueryHandler[K, V]

	// cfg is the provider's config. This contains the name, path, fetch timeout, and
	// fetch interval for the provider. To read more about the config, see the
	// oracle/config/provider.go file. To read more about how to configure a custom provider,
	// please see the providers/README.md file.
	cfg config.ProviderConfig

	// data is the latest set of key -> value pairs for the provider i.e. the latest prices
	// for a given set of currency pairs.
	data map[K]providertypes.Result[V]

	// ids is the set of IDs that the provider will fetch data for.
	ids []K
}

// NewProvider returns a new Base provider.
func NewProvider[K comparable, V any](
	logger *zap.Logger,
	cfg config.ProviderConfig,
	handler handlers.QueryHandler[K, V],
	ids []K,
) (*BaseProvider[K, V], error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid provider config %s", err)
	}

	if logger == nil {
		return nil, fmt.Errorf("no logger specified for provider %s", cfg.Name)
	}

	if handler == nil {
		return nil, fmt.Errorf("no query handler specified for provider %s", cfg.Name)
	}

	return &BaseProvider[K, V]{
		logger:  logger.With(zap.String("provider", cfg.Name)),
		cfg:     cfg,
		handler: handler,
		data:    make(map[K]providertypes.Result[V]),
		ids:     ids,
	}, nil
}

// Start starts the provider's main loop. The provider will fetch the data from the handler
// and continuously update the data. This blocks until the provider is stopped.
func (p *BaseProvider[K, V]) Start(ctx context.Context) error {
	p.logger.Info("starting provider", zap.Duration("interval", p.cfg.Interval))

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
func (p *BaseProvider[K, V]) GetData() map[K]providertypes.Result[V] {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Deep copy the prices into a new map.
	cpy := make(map[K]providertypes.Result[V])
	maps.Copy(cpy, p.data)

	return cpy
}
