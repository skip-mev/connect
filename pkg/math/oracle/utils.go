package oracle

import (
	"fmt"
	"math/big"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// GetTickerFromCurrencyPair returns the ticker for the given currency pair.
func (m *MedianAggregator) GetTickerFromCurrencyPair(
	cp slinkytypes.CurrencyPair,
) (mmtypes.Ticker, error) {
	cfg := m.GetMarketMap()
	market, ok := cfg.Markets[cp.String()]
	if !ok {
		return mmtypes.Ticker{}, fmt.Errorf("missing ticker: %s", cp.String())
	}

	return market.Ticker, nil
}

// GetProviderPrice returns the relevant provider price. Note that if the operation
// is for the index provider, then the price is retrieved from the previously calculated
// median prices. Otherwise, the price is retrieved from the provider cache. Additionally,
// this function normalizes (scales, inverts) the price to maintain the maximum precision.
func (m *MedianAggregator) GetProviderPrice(
	ticker mmtypes.Ticker,
	providerConfig mmtypes.ProviderConfig,
) (*big.Int, error) {
	cache := m.GetDataByProvider(providerConfig.Name)

	price, ok := cache[ticker]
	if !ok {
		return nil, fmt.Errorf("missing %s price for ticker: %s", providerConfig.Name, ticker.String())
	}

	scaledPrice, err := ScaleUpCurrencyPairPrice(ticker.Decimals, price)
	if err != nil {
		return nil, err
	}

	if providerConfig.Invert {
		scaledPrice = InvertCurrencyPairPrice(scaledPrice, ScaledDecimals)
	}

	return scaledPrice, nil
}

// GetIndexPrice returns the aggregated index price.
func (m *MedianAggregator) GetIndexPrice(
	currencyPair slinkytypes.CurrencyPair,
) (*big.Int, error) {
	cache := m.GetAggregatedData()
	targetTicker, err := m.GetTickerFromCurrencyPair(currencyPair)
	if err != nil {
		return nil, err
	}

	price, ok := cache[targetTicker]
	if !ok {
		return nil, fmt.Errorf("missing index price for ticker: %s", targetTicker.String())
	}

	return ScaleUpCurrencyPairPrice(targetTicker.Decimals, price)
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
