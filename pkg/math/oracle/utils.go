package oracle

import (
	"fmt"
	"math/big"

	"github.com/skip-mev/slinky/oracle/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// GetTickerFromOperation returns the ticker for the given operation.
func (m *MedianAggregator) GetTickerFromOperation(
	operation mmtypes.Operation,
) (mmtypes.Ticker, error) {
	ticker, ok := m.cfg.Tickers[operation.CurrencyPair.String()]
	if !ok {
		return mmtypes.Ticker{}, fmt.Errorf("missing ticker: %s", operation.CurrencyPair.String())
	}

	return ticker, nil
}

// GetProviderPrice returns the relevant provider price.
func (m *MedianAggregator) GetProviderPrice(
	operation mmtypes.Operation,
) (*big.Int, error) {
	ticker, err := m.GetTickerFromOperation(operation)
	if err != nil {
		return nil, err
	}

	// If the provider is not the index provider, then we can get the price
	// from the provider cache. Otherwise we want to retrieve the previously
	// calculated median price (index price).
	var cache types.TickerPrices
	if operation.Provider != IndexProviderPrice {
		cache = m.PriceAggregator.GetDataByProvider(operation.Provider)
	} else {
		cache = m.PriceAggregator.GetAggregatedData()
	}

	price, ok := cache[ticker]
	if !ok {
		return nil, fmt.Errorf("missing %s price for ticker: %s", operation.Provider, ticker.String())
	}

	// We scale the price up to the maximum precision to ensure that we can
	// perform the necessary calculations. If the price is inverted, then we
	// we can higher conversion precision.
	scaledPrice, err := ScaleUpCurrencyPairPrice(ticker.Decimals, price)
	if err != nil {
		return nil, err
	}

	if operation.Invert {
		scaledPrice = InvertCurrencyPairPrice(scaledPrice, ScaledDecimals)
	}

	return scaledPrice, nil
}
