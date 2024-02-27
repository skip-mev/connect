package orchestrator

import (
	"github.com/skip-mev/slinky/oracle/types"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// UpdateMarketMap updates the orchestrator's market map and updates the providers' market maps.
// Specifically, it determines if the provider's market map has a diff, and
func (o *ProviderOrchestrator) UpdateMarketMap(marketMap mmtypes.MarketMap) error {
	o.mut.Lock()
	defer o.mut.Unlock()

	if err := marketMap.ValidateBasic(); err != nil {
		return err
	}

	// Iterate over all of the existing providers and update their market maps.
	for name, state := range o.providers {
		providerMarketMap, err := types.ProviderMarketMapFromMarketMap(name, marketMap)
		if err != nil {
			return err
		}

		provider := state.Provider
		updater := provider.GetConfigUpdater()
		if updater == nil {
			continue
		}

		switch provider.Type() {
		case providertypes.API:
			// Create and update the API query handler.
			handler, err := o.createAPIQueryHandler(state.Cfg, providerMarketMap)
			if err != nil {
				return err
			}

			updater.UpdateAPIHandler(handler)
		case providertypes.WebSockets:
			// Create and update the WebSocket query handler.
			handler, err := o.createWebSocketQueryHandler(state.Cfg, providerMarketMap)
			if err != nil {
				return err
			}

			updater.UpdateWebSocketHandler(handler)
		}

		// Update the set of IDs that the provider is responsible for.
		updater.UpdateIDs(providerMarketMap.GetTickers())
		// Update the provider's state.
		state.Market = providerMarketMap
		state.Enabled = len(providerMarketMap.GetTickers()) > 0
		o.providers[name] = state
	}

	o.marketMap = marketMap
	return nil
}
