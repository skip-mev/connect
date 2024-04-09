package orchestrator

import (
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
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
		providerTickers, err := types.ProviderTickersFromMarketMap(name, marketMap)
		if err != nil {
			o.logger.Error("failed to create provider market map", zap.String("provider", name), zap.Error(err))
			return err
		}

		// Update the provider's state.
		updatedState, err := o.UpdateProviderState(providerTickers, state)
		if err != nil {
			o.logger.Error("failed to update provider state", zap.String("provider", name), zap.Error(err))
			return err
		}

		o.providers[name] = updatedState
	}

	o.marketMap = marketMap
	if o.aggregator != nil {
		o.aggregator.UpdateMarketMap(o.marketMap)
	}

	return nil
}

// UpdateProviderState updates the provider's state based on the market map. Specifically,
// this will update the provider's query handler and the provider's market map.
func (o *ProviderOrchestrator) UpdateProviderState(providerTickers []types.ProviderTicker, state ProviderState) (ProviderState, error) {
	provider := state.Provider

	o.logger.Info("updating provider state", zap.String("provider_state", provider.Name()))
	provider.Update(base.WithNewIDs[types.ProviderTicker, *big.Float](providerTickers))

	switch {
	case len(providerTickers) == 0:
		provider.Stop()
	case len(providerTickers) > 0 && !provider.IsRunning():
		o.wg.Add(1)
		go func() {
			defer o.wg.Done()
			o.execProviderFn(o.mainCtx, provider)
		}()
	}

	// Update the provider's state.
	o.logger.Info("updated provider state", zap.String("provider_state", provider.Name()))
	return state, nil
}
