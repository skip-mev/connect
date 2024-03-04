package base

import (
	"context"
	"fmt"
	"maps"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	wshandlers "github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// Provider implements a base provider that can be used to build other providers.
type Provider[K providertypes.ResponseKey, V providertypes.ResponseValue] struct {
	mu      sync.Mutex
	logger  *zap.Logger
	running atomic.Bool

	// name is the name of the provider.
	name string

	// api is the handler for the querying api data. Developer's implement this interface
	// to extend the provider's functionality. For example, this could be used to fetch
	// prices from an API, where K is the currency pair and V is the price. For more information
	// on how to implement a custom handler, please see the providers/base/README.md file.
	api apihandlers.APIQueryHandler[K, V]

	// apiCfg is the API configuration for the provider.
	apiCfg config.APIConfig

	// ws is the handler for the websocket data. Developers implement this interface to extend
	// the provider's functionality. For example, this could be used to fetch prices from a
	// websocket, where K is the currency pair and V is the price. For more information on how
	// to implement a custom handler, please see the providers/base/README.md file.
	ws wshandlers.WebSocketQueryHandler[K, V]

	// wsCfg is the websocket configuration for the provider.
	wsCfg config.WebSocketConfig

	// data is the latest set of key -> value pairs for the provider i.e. the latest prices
	// for a given set of currency pairs.
	data map[K]providertypes.Result[V]

	// ids is the set of IDs that the provider will fetch data for.
	ids []K

	// metrics is the metrics implementation for the provider.
	metrics providermetrics.ProviderMetrics

	// restartCh is the channel that is used to signal the provider to restart.
	restartCh chan struct{}

	// stopCh is the channel that is used to signal the provider to stop.
	stopCh chan struct{}
}

// NewProvider returns a new Base provider.
func NewProvider[K providertypes.ResponseKey, V providertypes.ResponseValue](opts ...ProviderOption[K, V]) (*Provider[K, V], error) {
	p := &Provider[K, V]{
		logger:    zap.NewNop(),
		ids:       make([]K, 0),
		data:      make(map[K]providertypes.Result[V]),
		restartCh: make(chan struct{}, 1),
		stopCh:    make(chan struct{}, 1),
	}

	for _, opt := range opts {
		opt(p)
	}

	switch {
	case p.api != nil && p.ws != nil:
		return nil, fmt.Errorf("cannot configure both api and websocket")
	case p.api == nil && p.ws == nil:
		return nil, fmt.Errorf("must configure either api or websocket")
	case p.apiCfg.ValidateBasic() != nil:
		return nil, fmt.Errorf("invalid api config")
	case p.wsCfg.ValidateBasic() != nil:
		return nil, fmt.Errorf("invalid websocket config")
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

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// If the config updater is set, the provider may update it's internal configurations
	// on the fly. As such, we need to listen for updates to the config updater and restart
	// the provider's main loop when the configuration changes.
	wg := sync.WaitGroup{}

	// Start the main loop. At a high level, the main loop will continuously fetch data from
	// the handler and update the provider's data. It allows for the provider to be restarted
	// when the configuration changes. The main loop will exit either when the context is
	// cancelled or the provider gets an unexpected error.
	var retErr error
MainLoop:
	for {
		// Create a new context for the fetch loop. This allows us to cancel the fetch loop
		// when the provider needs to be restarted.
		fetchCtx, cancelFetch := context.WithCancel(ctx)
		defer cancelFetch()

		// Start the fetch loop.
		errCh := make(chan error)
		wg.Add(1)
		go func() {
			errCh <- p.fetch(fetchCtx)
			wg.Done()
		}()

		select {
		case <-p.restartCh:
			// If any of the provider's configurations have changed, the provider will
			// be signalled to restart.
			p.logger.Info("restarting provider")
			cancelFetch()

			// Wait for the fetch loop to stop.
			err := <-errCh
			p.logger.Debug("provider fetch loop stopped", zap.Error(err))
			continue MainLoop
		case err := <-errCh:
			// If the fetch loop stops unexpectedly, we should return.
			retErr = err
			break MainLoop
		case <-ctx.Done():
			// If the context is cancelled, we should return. We expect the fetch go routine
			// to exit when the context is cancelled.
			retErr = <-errCh
			break MainLoop
		case <-p.stopCh:
			// If the provider is manually stopped, we stop the fetch loop and return.
			p.logger.Debug("stopping provider")
			cancel()
			retErr = <-errCh
			break MainLoop
		}
	}

	wg.Wait()
	p.logger.Info("wait group done")

	return retErr
}

// Stop stops the provider's main loop.
func (p *Provider[K, V]) Stop() {
	if !p.running.Load() {
		return
	}

	p.logger.Info("received manual stop signal")
	p.stopCh <- struct{}{}
}

// IsRunning returns true if the provider is running.
func (p *Provider[K, V]) IsRunning() bool {
	return p.running.Load()
}

// Name returns the name of the provider.
func (p *Provider[K, V]) Name() string {
	return p.name
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

// Type returns the type of data handler the provider uses.
func (p *Provider[K, V]) Type() providertypes.ProviderType {
	switch {
	case p.apiCfg.Enabled:
		return providertypes.API
	case p.wsCfg.Enabled:
		return providertypes.WebSockets
	default:
		return "unknown"
	}
}
