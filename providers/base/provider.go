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

// Provider implements a base provider that can be used to build other providers.
type Provider[K providertypes.ResponseKey, V providertypes.ResponseValue] struct {
	mu     sync.Mutex
	logger *zap.Logger

	// api is the handler for the querying api data. Developer's implement this interface
	// to extend the provider's functionality. For example, this could be used to fetch
	// prices from an API, where K is the currency pair and V is the price. For more information
	// on how to implement a custom handler, please see the providers/base/README.md file.
	api apihandlers.APIQueryHandler[K, V]

	// ws is the handler for the web socket data. Developers implement this interface to extend
	// the provider's functionality. For example, this could be used to fetch prices from a
	// websocket, where K is the currency pair and V is the price. For more information on how
	// to implement a custom handler, please see the providers/base/README.md file.
	ws wshandlers.WebSocketQueryHandler[K, V]

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

	// metrics is the metrics implementation for the provider.
	metrics providermetrics.ProviderMetrics
}

// NewProvider returns a new Base provider.
func NewProvider[K providertypes.ResponseKey, V providertypes.ResponseValue](
	cfg config.ProviderConfig,
	opts ...ProviderOption[K, V],
) (providertypes.Provider[K, V], error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid provider config %s", err)
	}

	p := &Provider[K, V]{
		cfg:    cfg,
		logger: zap.NewNop(),
		ids:    make([]K, 0),
		data:   make(map[K]providertypes.Result[V]),
	}

	for _, opt := range opts {
		opt(p)
	}

	if p.api == nil && p.ws == nil {
		return nil, fmt.Errorf("no query handler specified for provider %s", cfg.Name)
	}

	if p.api != nil && p.ws != nil {
		return nil, fmt.Errorf("cannot specify both an api and web socket query handler for provider %s", cfg.Name)
	}

	if p.metrics == nil {
		p.metrics = providermetrics.NewNopProviderMetrics()
	}

	return p, nil
}

// Start starts the provider's main loop. The provider will fetch the data from the handler
// and continuously update the data. This blocks until the provider is stopped.
func (p *Provider[K, V]) Start(ctx context.Context) error {
	p.logger.Info("starting provider")
	if len(p.ids) == 0 {
		p.logger.Warn("no ids to fetch")
	}

	// Start the main loop.
	return p.fetch(ctx)
}

// Name returns the name of the provider.
func (p *Provider[K, V]) Name() string {
	return p.cfg.Name
}

// GetData returns the latest data recorded by the provider. The data is constantly
// updated by the provider's main loop and provides access to the latest data - prices
// in constant time.
func (p *Provider[K, V]) GetData() map[K]providertypes.Result[V] {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Deep copy the prices into a new map.
	cpy := make(map[K]providertypes.Result[V])
	maps.Copy(cpy, p.data)

	return cpy
}

// Type returns the type of data handler the provider uses
func (p *Provider[K, V]) Type() providermetrics.ProviderType {
	switch {
	case p.cfg.API.Enabled:
		return providermetrics.API
	case p.cfg.WebSocket.Enabled:
		return providermetrics.WebSockets
	default:
		return "unknown"
	}
}
