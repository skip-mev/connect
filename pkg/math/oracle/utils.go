package oracle

import (
	"fmt"
	"maps"
	"math/big"

	"github.com/skip-mev/slinky/oracle/types"
	pkgtypes "github.com/skip-mev/slinky/pkg/types"
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
)

// GetProviderPrice returns the relevant provider price. Note that the aggregator's
// provider data cache stores prices in the form of providerName -> offChainTicker -> price.
func (m *MedianAggregator) GetProviderPrice(
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

	if cfg.Invert {
		return new(big.Float).Quo(big.NewFloat(1), price), nil
	}

	return price, nil
}

// GetIndexPrice returns the relevant index price. Note that the aggregator's
// index price cache stores prices in the form of ticker -> price.
func (m *MedianAggregator) GetIndexPrice(
	cp pkgtypes.CurrencyPair,
) (*big.Float, error) {
	price, ok := m.indexPrices[cp.String()]
	if !ok {
		return nil, fmt.Errorf("missing index price for ticker: %s", cp)
	}

	return price, nil
}

// SetIndexPrice sets the index price for the given currency pair.
func (m *MedianAggregator) SetIndexPrices(
	prices types.AggregatorPrices,
) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.indexPrices = prices
}

// GetIndexPrices returns the index prices the aggregator has.
func (m *MedianAggregator) GetIndexPrices() types.AggregatorPrices {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	cpy := make(types.AggregatorPrices)
	maps.Copy(cpy, m.indexPrices)

	return cpy
}

// UpdateMarketMap updates the market map for the oracle.
func (m *MedianAggregator) UpdateMarketMap(marketMap mmtypes.MarketMap) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.cfg = marketMap
}

// GetMarketMap returns the market map for the oracle.
func (m *MedianAggregator) GetMarketMap() *mmtypes.MarketMap {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return &m.cfg
}
