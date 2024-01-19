package base

import (
	"context"
	"fmt"
	"maps"
	"sync"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	wshandlers "github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// BaseProvider implements a base provider that can be used to build other providers.
type BaseProvider[K comparable, V any] struct { //nolint
	mu     sync.Mutex
	logger *zap.Logger

	// name is the name of the provider.
	name string

	// api is the handler for the querying api data. Developer's implement this interface
	// to extend the provider's functionality. For example, this could be used to fetch
	// prices from an API, where K is the currency pair and V is the price. For more information
	// on how to implement a custom handler, please see the providers/base/README.md file.
	api apihandlers.APIQueryHandler[K, V]

	// apiCfg is the API configuration for the provider.
	apiCfg config.APIConfig

	// ws is the handler for the web socket data. Developers implement this interface to extend
	// the provider's functionality. For example, this could be used to fetch prices from a
	// websocket, where K is the currency pair and V is the price. For more information on how
	// to implement a custom handler, please see the providers/base/README.md file.
	ws wshandlers.WebSocketQueryHandler[K, V]

	// wsCfg is the web socket configuration for the provider.
	wsCfg config.WebSocketConfig

	// data is the latest set of key -> value pairs for the provider i.e. the latest prices
	// for a given set of currency pairs.
	data map[K]providertypes.Result[V]

	// ids is the set of IDs that the provider will fetch data for.
	ids []K

	// metrics is the metrics implementation for the provider.
	metrics providermetrics.ProviderMetrics
}

// NewProvider returns a new Base provider.
func NewProvider[K comparable, V any](opts ...ProviderOption[K, V]) (providertypes.Provider[K, V], error) {
	p := &BaseProvider[K, V]{
		logger: zap.NewNop(),
		ids:    make([]K, 0),
		data:   make(map[K]providertypes.Result[V]),
	}

	for _, opt := range opts {
		opt(p)
	}

	switch {
	case p.api != nil && p.ws != nil:
		return nil, fmt.Errorf("cannot configure both api and web socket")
	case p.api == nil && p.ws == nil:
		return nil, fmt.Errorf("must configure either api or web socket")
	case p.apiCfg.ValidateBasic() != nil:
		return nil, fmt.Errorf("invalid api config")
	case p.wsCfg.ValidateBasic() != nil:
		return nil, fmt.Errorf("invalid web socket config")
	}

	if p.metrics == nil {
		p.metrics = providermetrics.NewNopProviderMetrics()
	}

	return p, nil
}

// Start starts the provider's main loop. The provider will fetch the data from the handler
// and continuously update the data. This blocks until the provider is stopped.
func (p *BaseProvider[K, V]) Start(ctx context.Context) error {
	p.logger.Info("starting provider")
	if len(p.ids) == 0 {
		p.logger.Warn("no ids to fetch")
	}

	// Start the main loop.
	return p.fetch(ctx)
}

// Name returns the name of the provider.
func (p *BaseProvider[K, V]) Name() string {
	return p.name
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
