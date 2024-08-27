package oracle

import (
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/base"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

// UpdateMarketMap updates the oracle's market map and updates the providers'
// market maps. Specifically, it determines if the provider's market map has a diff,
// and if so, updates the provider's state.
func (o *OracleImpl) UpdateMarketMap(marketMap mmtypes.MarketMap) error {
	o.mut.Lock()
	defer o.mut.Unlock()

	if err := marketMap.ValidateBasic(); err != nil {
		o.logger.Error("failed to validate market map", zap.Error(err))
		return err
	}

	// Iterate over all existing price providers and update their market maps.
	for name, state := range o.priceProviders {
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

		o.priceProviders[name] = updatedState
	}

	o.marketMap = marketMap
	if o.aggregator != nil {
		o.aggregator.UpdateMarketMap(o.marketMap)
	}

	return nil
}

// UpdateProviderState updates the provider's state based on the market map. Specifically,
// this will update the provider's query handler and the provider's market map.
func (o *OracleImpl) UpdateProviderState(providerTickers []types.ProviderTicker, state ProviderState) (ProviderState, error) {
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

	// Ignore sampling limits for provider update logs via injecting provider name in message
	o.logger.Info(
		fmt.Sprintf("updated %s provider state", provider.Name()),
		zap.String("provider", provider.Name()),
		zap.Int("num_tickers", len(provider.GetIDs())),
	)
	return state, nil
}

func (o *OracleImpl) fetchAllPrices() {
	o.logger.Debug("starting price fetch loop")
	defer func() {
		if r := recover(); r != nil {
			o.logger.Error("fetchAllPrices tick panicked", zap.Error(fmt.Errorf("%v", r)))
		}
	}()

	o.aggregator.Reset()

	// Retrieve the latest prices from each provider.
	o.mut.Lock()
	for _, provider := range o.priceProviders {
		o.fetchPrices(provider.Provider)
	}
	o.mut.Unlock()

	o.logger.Debug("oracle fetched prices from providers")

	// Compute aggregated prices and update the oracle.
	o.aggregator.AggregatePrices()
	o.setLastSyncTime(time.Now().UTC())

	// update the last sync time
	o.metrics.AddTick()
}

func (o *OracleImpl) fetchPrices(provider *types.PriceProvider) {
	defer func() {
		if r := recover(); r != nil {
			o.logger.Error(
				"provider panicked",
				zap.String("provider_name", provider.Name()),
				zap.Error(fmt.Errorf("%v", r)),
			)
		}
	}()

	if !provider.IsRunning() {
		o.logger.Debug(
			"provider is not running",
			zap.String("provider", provider.Name()),
		)

		return
	}

	o.logger.Debug(
		"retrieving prices",
		zap.String("provider", provider.Name()),
		zap.String("data handler type",
			string(provider.Type())),
	)

	// Fetch and set prices from the provider.
	prices := provider.GetData()
	if prices == nil {
		o.logger.Debug(
			"provider returned nil prices",
			zap.String("provider", provider.Name()),
			zap.String("data handler type", string(provider.Type())),
		)

		return
	}

	timeFilteredPrices := make(types.Prices)
	for pair, result := range prices {
		// If the price is older than the maxCacheAge, skip it.
		diff := time.Now().UTC().Sub(result.Timestamp)
		if diff > o.cfg.MaxPriceAge {
			o.logger.Debug(
				"skipping price",
				zap.String("provider", provider.Name()),
				zap.String("data handler type", string(provider.Type())),
				zap.String("pair", pair.String()),
				zap.Duration("diff", diff),
			)

			continue
		}

		o.logger.Debug(
			"adding price",
			zap.String("provider", provider.Name()),
			zap.String("data handler type", string(provider.Type())),
			zap.String("pair", pair.String()),
			zap.String("price", result.Value.String()),
			zap.Duration("diff", diff),
		)
		timeFilteredPrices[pair.GetOffChainTicker()] = result.Value
	}

	o.logger.Debug("provider returned prices",
		zap.String("provider", provider.Name()),
		zap.String("data handler type", string(provider.Type())),
		zap.Int("prices", len(prices)),
	)
	o.aggregator.SetProviderPrices(provider.Name(), timeFilteredPrices)
}

func (o *OracleImpl) setLastSyncTime(t time.Time) {
	o.mut.Lock()
	defer o.mut.Unlock()

	o.lastPriceSync = t
}
