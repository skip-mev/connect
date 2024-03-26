package oracle

import (
	"fmt"
	"math/big"

	"github.com/skip-mev/slinky/oracle/types"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
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
	var (
		err          error
		cache        types.TickerPrices
		targetTicker = ticker
	)
	isNormalize := providerConfig.NormalizeByPair != nil

	if !isNormalize {
		cache = m.GetDataByProvider(providerConfig.Name)
	} else {
		cache = m.GetAggregatedData()
		targetTicker, err = m.GetTickerFromCurrencyPair(*providerConfig.NormalizeByPair)
	}
	if err != nil {
		return nil, err
	}

	price, ok := cache[targetTicker]
	if !ok {
		return nil, fmt.Errorf("missing %s price for ticker: %s", providerConfig.Name, targetTicker.String())
	}

	scaledPrice, err := ScaleUpCurrencyPairPrice(targetTicker.Decimals, price)
	if err != nil {
		return nil, err
	}

	if providerConfig.Invert && !isNormalize {
		scaledPrice = InvertCurrencyPairPrice(scaledPrice, ScaledDecimals)
	}

	return scaledPrice, nil
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
