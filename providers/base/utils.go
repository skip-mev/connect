package base

import (
	"fmt"

	providertypes "github.com/skip-mev/slinky/providers/types"
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
