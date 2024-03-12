package oracle

import (
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math/median"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// MedianAggregator is an aggregator that calculates the median price for each ticker,
// resolved from the median prices of all price feeds. Specifically, this aggregator
// first resolves the median prices for all price feeds, and then calculates the median
// price for each ticker using the conversion markets. If no conversion markets are provided
// for a certain ticker, the final aggregated price will be the median of the prices calculated
// from the raw price feeds.
type MedianAggregator struct {
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

	return &MedianAggregator{
		logger: logger,
		cfg:    cfg,
	}, nil
}

// AggregateFn returns the aggregate function for the median price calculation. Specifically, this
// aggregation function first resolves the median prices for all price feeds, and then calculates
// the median price for each ticker using the conversion markets. If no conversion markets are provided
// for a certain ticker, the final aggregated price will be the median of the prices calculated from
// the raw price feeds.
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
func (m *MedianAggregator) AggregateFn() types.PriceAggregationFn {
	return func(feedsPerProvider types.AggregatedProviderPrices) types.TickerPrices {
		// Calculate the median price for each price feed.
		feedMedians := median.ComputeMedian()(feedsPerProvider)
		m.logger.Info("calculated median prices for raw price feeds", zap.Int("num_prices", len(feedMedians)))

		for ticker, price := range feedMedians {
			m.logger.Info(
				"got median price",
				zap.String("ticker", ticker.String()),
				zap.String("price", price.String()),
			)
		}

		// Calculate the converted prices for each ticker using the conversion markets.
		aggregatedMedians := m.GetConvertedPrices(feedMedians)
		m.logger.Info("calculated median prices for converted price feeds", zap.Int("num_prices", len(aggregatedMedians)))

		// Replace the median prices for each ticker with the aggregated median prices.
		// This will overwrite the median prices for each ticker with the final aggregated median prices.
		// This is necessary because the conversion markets may not be provided for all tickers.
		for ticker, price := range aggregatedMedians {
			m.logger.Info(
				"replacing median price",
				zap.String("ticker", ticker.String()),
				zap.String("price", price.String()),
			)

			feedMedians[ticker] = price
		}

		return feedMedians
	}
}

// GetConvertedPrices returns the converted prices for each ticker using the provided median
// prices and the conversion markets.
func (m *MedianAggregator) GetConvertedPrices(feedMedians types.TickerPrices) types.TickerPrices {
	// Scale all medians to a common number of decimals. This does not lose precision.
	scaledMedians := make(types.TickerPrices, len(feedMedians))
	for ticker, price := range feedMedians {
		scaledPrice, err := ScaleUpCurrencyPairPrice(ticker.Decimals, price)
		if err != nil {
			m.logger.Error(
				"failed to scale price",
				zap.Error(err),
				zap.String("ticker", ticker.String()),
				zap.Uint64("decimals", ticker.Decimals),
				zap.String("price", price.String()),
			)

			continue
		}

		scaledMedians[ticker] = scaledPrice
	}

	// Determine the final aggregated price for each ticker, specifically only for the tickers
	// that have a set of conversion markets.
	aggregatedMedians := make(types.TickerPrices)
	for tickerStr, paths := range m.cfg.Paths {
		ticker, ok := m.cfg.Tickers[tickerStr]
		if !ok {
			m.logger.Error(
				"failed to get ticker",
				zap.String("ticker", tickerStr),
			)

			continue
		}

		// Get the converted prices for set of convertable markets.
		// ex. BTC/USDT * USDT/USD = BTC/USD
		//     BTC/USDC * USDC/USD = BTC/USD
		convertedPrices := m.CalculateConvertedPrices(ticker, paths, scaledMedians)

		// If there were no converted prices, log an error and continue. In this case,
		// the final aggregated price will be the median of the prices calculated from the raw
		// price feeds - if any.
		if len(convertedPrices) == 0 {
			m.logger.Error("no converted prices", zap.String("ticker", ticker.String()))
			continue
		}

		// Take the median of the converted prices.
		aggregatedMedians[ticker] = median.CalculateMedian(convertedPrices)
	}

	// Scale all of the aggregated medians back to the original number of decimals.
	for ticker, price := range aggregatedMedians {
		unscaledPrice, err := ScaleDownCurrencyPairPrice(ticker.Decimals, price)
		if err != nil {
			m.logger.Error(
				"failed to scale price",
				zap.Error(err),
				zap.String("ticker", ticker.String()),
				zap.Uint64("decimals", ticker.Decimals),
				zap.String("price", price.String()),
			)

			continue
		}

		aggregatedMedians[ticker] = unscaledPrice
	}

	return aggregatedMedians
}

// CalculateConvertedPrices calculates the converted prices for each ticker using the
// provided median prices and the conversion markets.
//
// For example, if the oracle receives a price for BTC/USDT and USDT/USD, it can use the conversion
// market to convert the BTC/USDT price to BTC/USD. In this case, the medians map would contain
// the median prices for BTC/USDT and USDT/USD, and the conversions would contain a sorted list of
// operations to convert the price of BTC/USDT to BTC/USD i.e. BTC/USDT * USDT/USD = BTC/USD.
func (m *MedianAggregator) CalculateConvertedPrices(
	ticker mmtypes.Ticker,
	paths mmtypes.Paths,
	medians types.TickerPrices,
) []*big.Int {
	convertedPrices := make([]*big.Int, 0)

	for _, path := range paths.Paths {
		// Calculate the converted price.
		convertedPrice, err := m.CalculateConvertedPrice(ticker, path, medians)
		if err != nil {
			m.logger.Debug(
				"failed to calculate converted price",
				zap.Error(err),
				zap.String("ticker", ticker.String()),
				zap.Any("conversions", path),
			)

			continue
		}

		convertedPrices = append(convertedPrices, convertedPrice)
	}

	return convertedPrices
}

// CalculateConvertedPrice converts a set of median prices to a target ticker using a set of
// conversion operations.
func (m *MedianAggregator) CalculateConvertedPrice(
	target mmtypes.Ticker,
	path mmtypes.Path,
	medians types.TickerPrices,
) (*big.Int, error) {
	if err := path.ValidateBasic(); err != nil {
		return nil, err
	}

	// Ensure that the conversion is valid.
	if !path.Match(target.String()) {
		return nil, fmt.Errorf("path does not match target %s: %s", target.String(), path.String())
	}

	// Scalers for the number of decimals.
	one := ScaledOne(ScaledDecimals)
	zero := big.NewInt(0)

	operations := path.Operations
	if len(operations) == 0 {
		return zero, fmt.Errorf("no operations in path")
	}

	first := operations[0]
	cp := first.CurrencyPair

	// Get the median price for the first feed.
	price, err := m.getMedianPrice(cp, medians)
	if err != nil {
		return nil, fmt.Errorf("failed to get median price for feed %s: %w", cp.String(), err)
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
		zap.String("target_ticker", target.String()),
		zap.String("current_ticker", cp.String()),
		zap.String("tracking_price", price.String()),
		zap.String("median_price", price.String()),
	)

	for _, feed := range operations[1:] {
		// Get the median price for the feed.
		cp := feed.CurrencyPair
		median, err := m.getMedianPrice(cp, medians)
		if err != nil {
			return nil, fmt.Errorf("failed to get median price for feed %s: %w", cp.String(), err)
		}

		if median.Cmp(zero) == 0 {
			return zero, nil
		}

		// Invert the price if necessary.
		if feed.Invert {
			median = InvertCurrencyPairPrice(median, ScaledDecimals)
		}

		// Make the conversion.
		price = price.Mul(price, median)
		price = price.Div(price, one)

		m.logger.Debug(
			"got median price",
			zap.String("target_ticker", target.String()),
			zap.String("conversion_ticker", cp.String()),
			zap.String("tracking_price", price.String()),
			zap.String("median_price", median.String()),
		)
	}

	m.logger.Debug(
		"calculated converted price",
		zap.String("target_ticker", target.String()),
		zap.String("price", price.String()),
	)

	return price, nil
}

// getMedianPrice returns the median price for a given ticker from the provided prices.
func (m *MedianAggregator) getMedianPrice(cp slinkytypes.CurrencyPair, prices types.TickerPrices) (*big.Int, error) {
	ticker, ok := m.cfg.Tickers[cp.String()]
	if !ok {
		return nil, fmt.Errorf("invalid ticker: %s", cp.String())
	}

	// Get the price for the ticker.
	price, ok := prices[ticker]
	if !ok {
		return nil, fmt.Errorf("missing price for ticker %s", ticker.String())
	}

	return price, nil
}
