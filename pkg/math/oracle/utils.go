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
	cfg := m.GetMarketMap()
	ticker, ok := cfg.Tickers[operation.CurrencyPair.String()]
	if !ok {
		return mmtypes.Ticker{}, fmt.Errorf("missing ticker: %s", operation.CurrencyPair.String())
	}

	return ticker, nil
}

// GetProviderPrice returns the relevant provider price. Note that if the operation
// is for the index provider, then the price is retrieved from the previously calculated
// median prices. Otherwise, the price is retrieved from the provider cache. Additionally,
// this function normalizes (scales, inverts) the price to maintain the maximum precision.
func (m *MedianAggregator) GetProviderPrice(
	operation mmtypes.Operation,
) (*big.Int, error) {
	ticker, err := m.GetTickerFromOperation(operation)
	if err != nil {
		return nil, err
	}

	var cache types.TickerPrices
	if operation.Provider != mmtypes.IndexPrice {
		cache = m.PriceAggregator.GetDataByProvider(operation.Provider)
	} else {
		cache = m.GetAggregatedData()
	}

	price, ok := cache[ticker]
	if !ok {
		return nil, fmt.Errorf("missing %s price for ticker: %s", operation.Provider, ticker.String())
	}

	scaledPrice, err := ScaleUpCurrencyPairPrice(ticker.Decimals, price)
	if err != nil {
		return nil, err
	}

	if operation.Invert {
		scaledPrice = InvertCurrencyPairPrice(scaledPrice, ScaledDecimals)
	}

	return scaledPrice, nil
}
