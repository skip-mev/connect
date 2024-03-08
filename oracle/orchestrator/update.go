package orchestrator

import (
	"math/big"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"go.uber.org/zap"
)

// UpdateWithMarketMap updates the orchestrator's market map and updates the providers'
// market maps. Specifically, it determines if the provider's market map has a diff,
// and if so, updates the provider's state.
func (o *ProviderOrchestrator) UpdateWithMarketMap(marketMap mmtypes.MarketMap) error {
	o.mut.Lock()
	defer o.mut.Unlock()

	if err := marketMap.ValidateBasic(); err != nil {
		o.logger.Error("failed to validate market map", zap.Error(err))
		return err
	}

	// Iterate over all of the existing providers and update their market maps.
	for name, state := range o.providers {
		providerMarketMap, err := types.ProviderMarketMapFromMarketMap(name, marketMap)
		if err != nil {
			o.logger.Error("failed to create provider market map", zap.String("provider", name), zap.Error(err))
			return err
		}

		// Update the provider's state.
		updatedState, err := o.UpdateProviderState(providerMarketMap, state)
		if err != nil {
			o.logger.Error("failed to update provider state", zap.String("provider", name), zap.Error(err))
			return err
		}

		o.providers[name] = updatedState
	}

	o.marketMap = marketMap
	return nil
}

// UpdateProviderState updates the provider's state based on the market map. Specifically,
// this will update the provider's query handler and the provider's market map.
func (o *ProviderOrchestrator) UpdateProviderState(marketMap types.ProviderMarketMap, state ProviderState) (ProviderState, error) {
	provider := state.Provider

	o.logger.Info("updating provider state", zap.String("provider_state", provider.Name()))
	switch provider.Type() {
	case providertypes.API:
		// Create and update the API query handler.
		handler, err := o.createAPIQueryHandler(state.Cfg, marketMap)
		if err != nil {
			return state, err
		}

		provider.Update(
			base.WithNewAPIHandler(handler),
			base.WithNewIDs[mmtypes.Ticker, *big.Int](marketMap.GetTickers()),
		)
	case providertypes.WebSockets:
		// Create and update the WebSocket query handler.
		handler, err := o.createWebSocketQueryHandler(state.Cfg, marketMap)
		if err != nil {
			return state, err
		}

		provider.Update(
			base.WithNewWebSocketHandler(handler),
			base.WithNewIDs[mmtypes.Ticker, *big.Int](marketMap.GetTickers()),
		)
	}

	// Update the provider's state.
	state.Market = marketMap
	state.Enabled = len(marketMap.GetTickers()) > 0
	o.logger.Info("updated provider state", zap.String("provider_state", provider.Name()))
	return state, nil
}
