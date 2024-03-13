package oracle

import (
	"fmt"
	"math/big"

	"github.com/skip-mev/slinky/oracle/types"
	pkgtypes "github.com/skip-mev/slinky/pkg/types"
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
		cache = m.GetDataByProvider(operation.Provider)
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

// ValidateMarketMap validates the market map configuration and its expected configuration for
// this aggregator. In particular, this will
//
//  1. Ensure that the market map is valid (ValidateBasic). This ensure's that each of the provider's
//     markets are supported by the market map.
//  2. Ensure that each path has a corresponding ticker.
//  3. Ensure that each path has a valid number of operations.
//  4. Ensure that each operation has a valid ticker and that the provider supports the ticker.
func ValidateMarketMap(
	marketMap mmtypes.MarketMap,
) error {
	if err := marketMap.ValidateBasic(); err != nil {
		return fmt.Errorf("valid basic failed for market map: %w", err)
	}

	for ticker, paths := range marketMap.Paths {
		// The ticker must be supported by the market map. Otherwise we do not how to resolve the
		// prices.
		if _, ok := marketMap.Tickers[ticker]; !ok {
			return fmt.Errorf("path includes a ticker that is not supported: %s", ticker)
		}

		for _, path := range paths.Paths {
			operations := path.Operations
			if len(operations) == 0 || len(operations) > MaxConversionOperations {
				return fmt.Errorf(
					"the expected number of operations is between 1 and %d; got %d operations for %s",
					MaxConversionOperations,
					len(operations),
					ticker,
				)
			}

			first := operations[0]
			if _, ok := marketMap.Tickers[first.CurrencyPair.String()]; !ok {
				return fmt.Errorf("operation included a ticker that is not supported: %s", first.CurrencyPair.String())
			}
			if err := checkIfProviderSupportsTicker(first.Provider, first.CurrencyPair, marketMap); err != nil {
				return err
			}

			if len(operations) != 2 {
				continue
			}

			second := operations[1]
			if second.Provider != IndexPrice {
				return fmt.Errorf("expected index price provider for second operation; got %s", second.Provider)
			}
			if _, ok := marketMap.Tickers[second.CurrencyPair.String()]; !ok {
				return fmt.Errorf("index operation included a ticker that is not supported: %s", second.CurrencyPair.String())
			}
		}
	}

	return nil
}

// checkIfProviderSupportsTicker checks if the provider supports the given ticker.
func checkIfProviderSupportsTicker(
	provider string,
	cp pkgtypes.CurrencyPair,
	marketMap mmtypes.MarketMap,
) error {
	providers, ok := marketMap.Providers[cp.String()]
	if !ok {
		return fmt.Errorf("provider %s included a ticker %s that has no providers supporting it", provider, cp.String())
	}

	for _, p := range providers.Providers {
		if p.Name == provider {
			return nil
		}
	}

	return fmt.Errorf("provider %s does not support ticker: %s", provider, cp.String())
}
