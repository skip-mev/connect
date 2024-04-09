package oracle

import (
	"fmt"
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
	cache := m.GetDataByProvider(cfg.Name)
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
	price, ok := m.GetIndexPrices()[cp.String()]
	if !ok {
		return nil, fmt.Errorf("missing index price for ticker: %s", cp)
	}

	return price, nil
}

// UpdateMarketMap updates the market map for the oracle.
func (m *MedianAggregator) UpdateMarketMap(marketMap mmtypes.MarketMap) {
	m.Lock()
	defer m.Unlock()

	m.cfg = marketMap
}

// GetMarketMap returns the market map for the oracle.
func (m *MedianAggregator) GetMarketMap() *mmtypes.MarketMap {
	m.Lock()
	defer m.Unlock()

	return &m.cfg
}

// SetPrices sets the prices for the oracle. These are the scaled prices that can be consumed
// by external providers as well as the unscaled index prices.
func (m *MedianAggregator) SetPrices(
	indexPrices, scaledPrices types.AggregatorPrices,
) {
	m.SetIndexPrices(indexPrices)
	m.DataAggregator.SetAggregatedData(scaledPrices)
}

// GetIndexPrices returns the index prices for the oracle.
func (m *MedianAggregator) GetIndexPrices() types.AggregatorPrices {
	m.Lock()
	defer m.Unlock()

	return m.indexPrices
}

// SetIndexPrices sets the index prices for the oracle.
func (m *MedianAggregator) SetIndexPrices(indexPrices types.AggregatorPrices) {
	m.Lock()
	defer m.Unlock()

	m.indexPrices = indexPrices
}
