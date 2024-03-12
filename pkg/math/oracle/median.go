package oracle

import (
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math/median"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

const (
	// MaxConversionOperations is the maximum number of conversion operations that can be used
	// to convert a price from one ticker to another.
	MaxConversionOperations = 2

	// IndexProviderPrice is the provider name for the index price.
	IndexProviderPrice = "index"
)

// MedianAggregator is an aggregator that calculates the median price for each ticker,
// resolved from a predefined set of conversion markets. In particular, these conversion
// markets are defined in the marketmap configuration that is provided to the aggregator.
type MedianAggregator struct {
	*types.PriceAggregator
	logger *zap.Logger
	cfg    mmtypes.MarketMap
}

// NewMedianAggregator returns a new Median aggregator.
func NewMedianAggregator(logger *zap.Logger, cfg mmtypes.MarketMap) (*MedianAggregator, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if err := cfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	m := &MedianAggregator{
		logger:          logger,
		cfg:             cfg,
		PriceAggregator: types.NewPriceAggregator(),
	}

	return m, nil
}

// AggregateFn returns the aggregate function for the median price calculation. Specifically, this
// aggregation function utilizes the previously cached prices to determine the current median prices
// by converting the prices using the predefined conversion markets. These conversion markets define
// the set of markets that can be used to convert the prices of a set of tickers to a common ticker.
func (m *MedianAggregator) AggregatedData() {
	updatedPrices := make(types.TickerPrices)
	for ticker, paths := range m.cfg.Paths {
		target, ok := m.cfg.Tickers[ticker]
		if !ok {
			m.logger.Error(
				"failed to get ticker; skipping aggregation",
				zap.String("ticker", ticker),
			)

			continue
		}

		// Get the converted prices for set of convertable markets.
		// ex. BTC/USDT * Index USDT/USD = BTC/USD
		//     BTC/USDC * Index USDC/USD = BTC/USD
		convertedPrices := m.CalculateConvertedPrices(
			target,
			paths,
		)

		// We need to have at least the minimum number of providers to calculate the median.
		if len(convertedPrices) < int(target.MinProviderCount) {
			m.logger.Error(
				"insufficient amount of converted prices",
				zap.String("ticker", ticker),
				zap.Int("num_converted_prices", len(convertedPrices)),
				zap.Any("converted_prices", convertedPrices),
				zap.Int("min_provider_count", int(target.MinProviderCount)),
			)

			continue
		}

		// Take the median of the converted prices.
		price := median.CalculateMedian(convertedPrices)
		updatedPrices[target] = price
		m.logger.Info(
			"calculated median price",
			zap.String("ticker", ticker),
			zap.String("price", price.String()),
			zap.Any("converted_prices", convertedPrices),
		)

	}

	m.logger.Info("calculated median prices for price feeds", zap.Int("num_prices", len(updatedPrices)))
	m.PriceAggregator.SetAggregatedData(updatedPrices)
}

// CalculateConvertedPrices calculates the converted prices for each ticker using the
// provided median prices and the conversion markets.
func (m *MedianAggregator) CalculateConvertedPrices(
	target mmtypes.Ticker,
	paths mmtypes.Paths,
) []*big.Int {
	m.logger.Debug("calculating converted prices", zap.String("ticker", target.String()))
	if len(paths.Paths) == 0 {
		m.logger.Error(
			"no conversion paths",
			zap.String("ticker", target.String()),
		)

		return nil
	}

	convertedPrices := make([]*big.Int, 0)
	for _, path := range paths.Paths {
		// Ensure that the number of operations is valid.
		if len(path.Operations) > MaxConversionOperations || len(path.Operations) == 0 {
			m.logger.Error(
				"invalid number of operations",
				zap.String("ticker", target.String()),
				zap.Any("path", path),
			)

			continue
		}

		// Calculate the converted price.
		adjustedPrice, err := m.CalculateAdjustedPrice(target, path.Operations)
		if err != nil {
			m.logger.Debug(
				"failed to calculate converted price",
				zap.Error(err),
				zap.String("ticker", target.String()),
				zap.Any("conversions", path),
			)

			continue
		}

		convertedPrices = append(convertedPrices, adjustedPrice)
		m.logger.Debug(
			"calculated converted price",
			zap.String("ticker", target.String()),
			zap.String("price", adjustedPrice.String()),
			zap.Any("conversions", path.Operations),
		)
	}

	return convertedPrices
}

// CalculateAdjustedPrice calculates an adjusted price for a given set of operations (if applicable).
// In particular, this assumes that every operation is either:
//
//  1. A direct conversion from the base ticker to the target ticker i.e. we want BTC/USD and
//     we have BTC/USD from a provider (e.g. Coinbase).
//  2. We need to convert the price of a given asset against the index price of the asset.
//
// In the first case, we can simply return the price of the provider. In the second case, we need
// to adjust the price by the index price of the asset. If the index price is not available, we
// return an error.
func (m *MedianAggregator) CalculateAdjustedPrice(
	target mmtypes.Ticker,
	operations []mmtypes.Operation,
) (*big.Int, error) {
	price, err := m.GetProviderPrice(operations[0])
	if err != nil {
		return nil, err
	}

	// If we have a single operation, then we can simply return the price. This implies that
	// we have a direct conversion from the base ticker to the target ticker.
	if len(operations) == 1 {
		return ScaleDownCurrencyPairPrice(target.Decimals, price)
	}

	adjustableByMarketPrice, err := m.GetProviderPrice(operations[1])
	if err != nil {
		return nil, err
	}

	// Make sure that the price is adjusted by the market price.
	adjustedPrice := big.NewInt(0).Mul(price, adjustableByMarketPrice)
	adjustedPrice = adjustedPrice.Div(adjustedPrice, ScaledOne(ScaledDecimals))

	return ScaleDownCurrencyPairPrice(target.Decimals, adjustedPrice)
}
