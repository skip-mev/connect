package oracle

import (
	"fmt"
	"maps"
	"math/big"

	"github.com/skip-mev/connect/v2/oracle/types"
	pkgtypes "github.com/skip-mev/connect/v2/pkg/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

// GetProviderPrice returns the relevant provider price. Note that the aggregator's
// provider data cache stores prices in the form of providerName -> offChainTicker -> price.
func (m *IndexPriceAggregator) GetProviderPrice(
	cfg mmtypes.ProviderConfig,
) (*big.Float, error) {
	cache, ok := m.providerPrices[cfg.Name]
	if !ok {
		return nil, fmt.Errorf("missing provider prices for provider: %s", cfg.Name)
	}

	price, ok := cache[cfg.OffChainTicker]
	if !ok {
		return nil, fmt.Errorf("missing %s price for ticker: %s", cfg.Name, cfg.OffChainTicker)
	}

	if price == nil {
		return nil, fmt.Errorf("price for %s ticker %s is nil", cfg.Name, cfg.OffChainTicker)
	}

	if cfg.Invert {
		return new(big.Float).Quo(big.NewFloat(1), price), nil
	}

	return price, nil
}

// GetIndexPrice returns the relevant index price. Note that the aggregator's
// index price cache stores prices in the form of ticker -> price.
func (m *IndexPriceAggregator) GetIndexPrice(
	cp pkgtypes.CurrencyPair,
) (*big.Float, error) {
	price, ok := m.indexPrices[cp.String()]
	if !ok {
		return nil, fmt.Errorf("missing index price for ticker: %s", cp)
	}

	if price == nil {
		return nil, fmt.Errorf("index price for ticker %s is nil", cp)
	}

	return price, nil
}

// SetIndexPrices sets the index price for the given currency pair.
func (m *IndexPriceAggregator) SetIndexPrices(
	prices types.Prices,
) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.indexPrices = prices
}

// GetIndexPrices returns the index prices the aggregator has.
func (m *IndexPriceAggregator) GetIndexPrices() types.Prices {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	cpy := make(types.Prices)
	maps.Copy(cpy, m.indexPrices)

	return cpy
}

// UpdateMarketMap updates the market map for the oracle.
func (m *IndexPriceAggregator) UpdateMarketMap(marketMap mmtypes.MarketMap) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.cfg = marketMap
}

// GetMarketMap returns the market map for the oracle.
func (m *IndexPriceAggregator) GetMarketMap() *mmtypes.MarketMap {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return &m.cfg
}

// SetProviderPrices updates the data aggregator with the given provider and data.
func (m *IndexPriceAggregator) SetProviderPrices(provider string, data types.Prices) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if data == nil {
		data = make(types.Prices)
	}

	m.providerPrices[provider] = data
}

// Reset resets the data aggregator for all providers.
func (m *IndexPriceAggregator) Reset() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.providerPrices = make(map[string]types.Prices)
}

// GetPrices returns the aggregated data the aggregator has. Specifically, the
// prices returned are the scaled prices - where each price is scaled by the
// respective ticker's decimals.
func (m *IndexPriceAggregator) GetPrices() types.Prices {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	cpy := make(types.Prices)
	maps.Copy(cpy, m.scaledPrices)

	return cpy
}
