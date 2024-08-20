package base

import (
	"context"
	"fmt"

	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

// createResponseCh creates the response channel for the provider.
func (p *Provider[K, V]) createResponseCh() error {
	// responseCh is used to receive the response(s) from the query handler.
	switch {
	case p.Type() == providertypes.API:
		// If the provider is an API provider, then the buffer size is set to the number of IDs.
		p.responseCh = make(chan providertypes.GetResponse[K, V], len(p.GetIDs()))
	case p.Type() == providertypes.WebSockets:
		// Otherwise, the buffer size is set to the max buffer size configured for the websocket.
		p.responseCh = make(chan providertypes.GetResponse[K, V], p.wsCfg.MaxBufferSize)
	default:
		return fmt.Errorf("no api or websocket configured")
	}

	return nil
}

// setMainCtx sets the main context for the provider.
func (p *Provider[K, V]) setMainCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.mainCtx, p.cancelMainFn = context.WithCancel(ctx)
	return p.mainCtx, p.cancelMainFn
}

// getMainCtx returns the main context for the provider.
func (p *Provider[K, V]) getMainCtx() (context.Context, context.CancelFunc) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.mainCtx, p.cancelMainFn
}

// setFetchCtx sets the fetch context for the provider.
func (p *Provider[K, V]) setFetchCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.fetchCtx, p.cancelFetchFn = context.WithCancel(ctx)
	return p.fetchCtx, p.cancelFetchFn
}

// getFetchCtx returns the fetch context for the provider.
func (p *Provider[K, V]) getFetchCtx() (context.Context, context.CancelFunc) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.fetchCtx, p.cancelFetchFn
}
