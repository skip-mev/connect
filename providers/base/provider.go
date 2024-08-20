package base

import (
	"context"
	"fmt"
	"maps"
	"sync"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	apihandlers "github.com/skip-mev/connect/v2/providers/base/api/handlers"
	providermetrics "github.com/skip-mev/connect/v2/providers/base/metrics"
	wshandlers "github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

// Provider implements a base provider that can be used to build other providers.
type Provider[K providertypes.ResponseKey, V providertypes.ResponseValue] struct {
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

	// ws is the handler for the websocket data. Developers implement this interface to extend
	// the provider's functionality. For example, this could be used to fetch prices from a
	// websocket, where K is the currency pair and V is the price. For more information on how
	// to implement a custom handler, please see the providers/base/README.md file.
	ws wshandlers.WebSocketQueryHandler[K, V]

	// wsCfg is the websocket configuration for the provider.
	wsCfg config.WebSocketConfig

	// data is the latest set of key -> value pairs for the provider i.e. the latest prices
	// for a given set of currency pairs.
	data map[K]providertypes.ResolvedResult[V]

	// ids is the set of IDs that the provider will fetch data for.
	ids []K

	// metrics is the metrics implementation for the provider.
	metrics providermetrics.ProviderMetrics

	// fetchCtx is the context for the fetch function.
	fetchCtx context.Context

	// cancelFetchFn is the function that is used to cancel the fetch loop.
	cancelFetchFn context.CancelFunc

	// mainCtx is the context for the main loop.
	mainCtx context.Context

	// cancelMainFn is the function that is used to cancel the main loop.
	cancelMainFn context.CancelFunc

	// responseCh is the channel that is used to receive the response(s) from the query handler.
	responseCh chan providertypes.GetResponse[K, V]
}

// NewProvider returns a new Base provider.
func NewProvider[K providertypes.ResponseKey, V providertypes.ResponseValue](opts ...ProviderOption[K, V]) (*Provider[K, V], error) {
	p := &Provider[K, V]{
		logger: zap.NewNop(),
		ids:    make([]K, 0),
		data:   make(map[K]providertypes.ResolvedResult[V]),
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
	if ctx == nil {
		p.logger.Error("context is nil; exiting")
		return nil
	}

	p.logger.Info("starting provider")
	mainCtx, mainCancel := p.setMainCtx(ctx)
	defer mainCancel()

	wg := sync.WaitGroup{}

	// Start the main loop. At a high level, the main loop will continuously fetch data from
	// the handler and update the provider's data. It allows for the provider to be restarted
	// when the configuration changes. The main loop will exit either when the context is
	// cancelled or the provider gets an unexpected error.
	var (
		retErr error
	)
	for {
		// Ensure that the provider has IDs set. This could be reset if the provider is
		// restarted / reconfigured.
		if len(p.GetIDs()) == 0 {
			p.logger.Debug("no ids set on provider; exiting")
			return nil
		}

		// Create the response channel for the provider. This channel is used to receive the
		// response(s) from the query handler.
		if err := p.createResponseCh(); err != nil {
			return err
		}

		// Create a new context for the fetch loop. This allows us to cancel the fetch loop
		// when the provider needs to be restarted.
		fetchCtx, fetchCancel := p.setFetchCtx(mainCtx)

		// Start the receive loop.
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.recv(fetchCtx)
		}()

		// Start the fetch loop.
		errCh := make(chan error, 1)
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- p.fetch(fetchCtx)
			fetchCancel()
			close(p.responseCh)
		}()

		// Wait for the fetch loop to return or the context to be cancelled.
		p.logger.Debug("started provider fetch and recv routines")
		wg.Wait()
		retErr = <-errCh
		p.logger.Debug("provider routines stopped", zap.Error(retErr))

		// If the provider was stopped due to a context cancellation, then we should
		// not restart the provider.
		mainCtx, _ := p.getMainCtx()
		if mainCtx.Err() != nil {
			p.logger.Info(
				"main provider context has been cancelled; provider is exiting",
				zap.Error(mainCtx.Err()),
			)

			break
		}
	}

	return retErr
}

// Stop stops the provider's main loop.
func (p *Provider[K, V]) Stop() {
	mainCtx, cancelMain := p.getMainCtx()
	if mainCtx == nil {
		p.logger.Debug("provider is not running")
		return
	}

	select {
	case <-mainCtx.Done():
		// The provider is already stopped.
		p.logger.Debug("provider is not running")
		return
	default:
		// Cancel the main context to stop the provider.
		p.logger.Debug("manually stopping provider")
		cancelMain()
	}
}

// IsRunning returns true if the provider is running.
func (p *Provider[K, V]) IsRunning() bool {
	mainCtx, _ := p.getMainCtx()
	if mainCtx == nil {
		return false
	}

	select {
	case <-mainCtx.Done():
		return false
	default:
		return true
	}
}

// Name returns the name of the provider.
func (p *Provider[K, V]) Name() string {
	return p.name
}

// GetData returns the latest data recorded by the provider. The data is constantly
// updated by the provider's main loop and provides access to the latest data - prices
// in constant time.
func (p *Provider[K, V]) GetData() map[K]providertypes.ResolvedResult[V] {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Deep copy the prices into a new map.
	cpy := make(map[K]providertypes.ResolvedResult[V])
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
