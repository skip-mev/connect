package oracle

import (
	"fmt"
	"math/big"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"go.uber.org/zap"
)

// MedianAggregator is an aggregator that calculates the median price for each currency pair,
// resolved from the median prices of all price feeds.
type MedianAggregator struct {
	logger *zap.Logger
	cfg    config.AggregateMarketConfig
}

// NewMedianAggregator returns a new Median aggregator.
func NewMedianAggregator(logger *zap.Logger, cfg config.AggregateMarketConfig) (*MedianAggregator, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if err := cfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &MedianAggregator{
		logger: logger,
		cfg:    cfg,
	}, nil
}

// AggregateFn returns the aggregate function for the median price calculation. Specifically, this
// aggregation function first resolves the median prices for all price feeds, and then calculates
// the median price for each currency pair using the conversion markets.
//
// For example, if the oracle receives price updates for
//   - BTC/USDT
//   - USDT/USD
//   - BTC/USDC
//   - USDC/USD
//   - BTC/USD
//
// This function will first calculate the median prices for BTC/USDT, USDT/USD, BTC/USDC, and USDC/USD,
// and then calculate the median price for BTC/USD using the conversion markets. Specifically, it will
// calculate the median price for BTC/USD using the following convertable markets:
//  1. BTC/USDT * USDT/USD = BTC/USD
//  2. BTC/USDC * USDC/USD = BTC/USD
//  3. BTC/USD = BTC/USD
//
// The final median price for BTC/USD will be the median of the prices calculated from the above
// calculations.
func (m *MedianAggregator) AggregateFn() aggregator.AggregateFn[string, map[oracletypes.CurrencyPair]*big.Int] {
	return func(
		feedsPerProvider aggregator.AggregatedProviderData[string, map[oracletypes.CurrencyPair]*big.Int],
	) map[oracletypes.CurrencyPair]*big.Int {
		// Calculate the median price for each price feed.
		feedMedians := aggregator.ComputeMedian()(feedsPerProvider)
		m.logger.Debug("calculated median prices for raw price feeds", zap.Any("num_prices", len(feedMedians)))

		// Scale all of the medians to a common number of decimals. This does not lose precision.
		for cp, price := range feedMedians {
			scaledPrice, err := ScaleUpCurrencyPairPrice(int64(cp.Decimals()), price)
			if err != nil {
				m.logger.Error("failed to scale price", zap.Error(err), zap.String("currency_pair", cp.String()), zap.String("price", price.String()))
				continue
			}

			feedMedians[cp] = scaledPrice
		}

		// Determine the final aggregated price for each currency pair.
		aggregatedMedians := make(map[oracletypes.CurrencyPair]*big.Int)
		for _, cfg := range m.cfg.AggregatedFeeds {
			// Get the converted price for each convertable feed.
			convertedPrices := m.CalculateConvertedPrices(cfg, feedMedians)

			// If there were no converted prices, log an error and continue.
			cp := cfg.CurrencyPair
			if len(convertedPrices) == 0 {
				m.logger.Debug("no converted prices", zap.String("currency_pair", cp.String()))
				continue
			}

			// Take the median of the converted prices.
			aggregatedMedians[cp] = aggregator.CalculateMedian(convertedPrices)
			m.logger.Debug("calculated median price", zap.String("currency_pair", cp.String()), zap.String("price", aggregatedMedians[cp].String()), zap.Any("converted_prices", convertedPrices))
		}

		// Scale all of the aggregated medians back to the original number of decimals.
		for cp, price := range aggregatedMedians {
			unscaledPrice, err := ScaleDownCurrencyPair(int64(cp.Decimals()), price)
			if err != nil {
				m.logger.Error("failed to scale price", zap.Error(err), zap.String("currency_pair", cp.String()), zap.String("price", price.String()))
				continue
			}

			aggregatedMedians[cp] = unscaledPrice
		}

		return aggregatedMedians
	}
}

// CalculateConvertedPrices calculates the converted prices for each currency pair using the
// provided median prices and the conversion markets.
//
// For example, if the oracle receives a price for BTC/USDT and USDT/USD, it can use the conversion
// market to convert the BTC/USDT price to BTC/USD. In this case, the medians map would contain
// the median prices for BTC/USDT and USDT/USD, and the conversions would contain a sorted list of
// operations to convert the price of BTC/USDT to BTC/USD i.e. BTC/USDT * USDT/USD = BTC/USD.
func (m *MedianAggregator) CalculateConvertedPrices(
	cfg config.AggregateFeedConfig,
	medians map[oracletypes.CurrencyPair]*big.Int,
) []*big.Int {
	convertedPrices := make([]*big.Int, 0)
	cp := cfg.CurrencyPair

	for _, conversion := range cfg.Conversions {
		// Ensure that the conversion is valid.
		if err := config.CheckSort(cp, conversion); err != nil {
			m.logger.Error("invalid conversion", zap.Error(err), zap.Any("conversions", conversion))
			continue
		}

		// Calculate the converted price.
		convertedPrice, err := m.CalculateConvertedPrice(cp, conversion, medians)
		if err != nil {
			m.logger.Debug("failed to calculate converted price", zap.Error(err), zap.Any("conversions", conversion))
			continue
		}

		convertedPrices = append(convertedPrices, convertedPrice)
	}

	return convertedPrices
}

// CalculateConvertedPrice converts the price of a currency pair to another currency pair using the
// provided median prices and the conversion markets. This function always scales the price to the
// maximum number of decimals in the conversion markets before normalizing the price at the end.
func (m *MedianAggregator) CalculateConvertedPrice(
	outcome oracletypes.CurrencyPair,
	operations []config.Conversion,
	medians map[oracletypes.CurrencyPair]*big.Int,
) (*big.Int, error) {
	if len(operations) == 0 {
		return nil, fmt.Errorf("no conversion operations")
	}

	// Scalers for the number of decimals.
	one := ScaledOne(ScaledDecimals)
	zero := big.NewInt(0)

	first := operations[0]
	cp := first.CurrencyPair

	// Get the median price for the first feed.
	price, ok := medians[cp]
	if !ok {
		return nil, fmt.Errorf("missing median price for feed %s", first.CurrencyPair.String())
	}

	if price.Cmp(zero) == 0 {
		return zero, nil
	}

	// If the first feed is inverted, invert the price scaled to the number of decimals.
	if first.Invert {
		price = InvertCurrencyPairPrice(price, ScaledDecimals)
	}

	m.logger.Debug(
		"got median price",
		zap.String("target_currency_pair", outcome.String()),
		zap.String("current_currency_pair", cp.String()),
		zap.String("tracking_price", price.String()),
		zap.String("median_price", price.String()),
	)

	for _, feed := range operations[1:] {
		// Get the median price for the feed.
		cp := feed.CurrencyPair
		median, ok := medians[cp]
		if !ok {
			return nil, fmt.Errorf("missing median price for feed %s", feed.CurrencyPair.String())
		}

		if median.Cmp(zero) == 0 {
			return zero, nil
		}

		// Invert the price if necessary.
		if feed.Invert {
			median = InvertCurrencyPairPrice(median, ScaledDecimals)
		}

		// Scale the median price to the number of decimals.
		price = price.Mul(price, median)
		price = price.Div(price, one)

		m.logger.Debug(
			"got median price",
			zap.String("target_currency_pair", outcome.String()),
			zap.String("conversion_currency_pair", cp.String()),
			zap.String("tracking_price", price.String()),
			zap.String("median_price", median.String()),
		)
	}

	m.logger.Debug("calculated converted price", zap.String("target_currency_pair", outcome.String()), zap.String("price", price.String()))
	return price, nil
}
