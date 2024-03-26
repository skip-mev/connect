package oracle

import (
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/aggregator"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math/median"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// MedianAggregator is an aggregator that calculates the median price for each ticker,
// resolved from a predefined set of conversion markets. A conversion market is a set of
// markets that can be used to convert the prices of a set of tickers to a common ticker.
// These are defined in the market map configuration.
type MedianAggregator struct {
	*aggregator.DataAggregator[string, types.TickerPrices]
	logger  *zap.Logger
	cfg     mmtypes.MarketMap
	metrics oraclemetrics.Metrics
}

// NewMedianAggregator returns a new Median aggregator.
func NewMedianAggregator(
	logger *zap.Logger,
	cfg mmtypes.MarketMap,
	metrics oraclemetrics.Metrics,
) (*MedianAggregator, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if err := cfg.ValidateBasic(); err != nil {
		return nil, err
	}

	if cfg.AggregationType != mmtypes.AggregationType_INDEX_PRICE_AGGREGATION {
		logger.Error("invalid aggregation type; please ensure the correct market config is used")
		return nil, fmt.Errorf("invalid aggregation type; expected %s got: %s", mmtypes.AggregationType_INDEX_PRICE_AGGREGATION, cfg.AggregationType)
	}

	if metrics == nil {
		logger.Warn("metrics is nil; using a no-op metrics implementation")
		metrics = oraclemetrics.NewNopMetrics()
	}

	return &MedianAggregator{
		logger:         logger,
		cfg:            cfg,
		metrics:        metrics,
		DataAggregator: aggregator.NewDataAggregator[string, types.TickerPrices](),
	}, nil
}

// AggregateData implements the aggregate function for the median price calculation. Specifically, this
// aggregation function aggregates the prices seen by each provider by first converting each price to a
// common ticker and then calculating the median of the converted prices. Prices are converted either
//
//  1. Directly from the base ticker to the target ticker. i.e. I have BTC/USD and I want BTC/USD.
//  2. Using the index price of an asset. i.e. I have BTC/USDT and I want BTC/USD. I can convert
//     BTC/USDT to BTC/USD using the index price of USDT/USD.
//
// The index price cache contains the previously calculated median prices.
func (m *MedianAggregator) AggregateData() {
	cfg := m.GetMarketMap()
	updatedPrices := make(types.TickerPrices)
	for ticker, paths := range cfg.Paths {
		target, ok := cfg.Tickers[ticker]
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

		// Take the median of the converted prices. This takes the average of the middle two
		// prices if the number of prices is even.
		price := median.CalculateMedian(convertedPrices)
		updatedPrices[target] = price
		m.logger.Info(
			"calculated median price",
			zap.String("ticker", ticker),
			zap.String("price", price.String()),
			zap.Any("converted_prices", convertedPrices),
		)
		m.metrics.AddTickerTick(target.String())

		floatPrice, _ := price.Float64()
		m.metrics.UpdateAggregatePrice(target.String(), target.GetDecimals(), floatPrice)
	}

	// Update the aggregated data. These prices are going to be used as the index prices the
	// next time we calculate prices.
	m.DataAggregator.SetAggregatedData(updatedPrices)
	m.logger.Info("calculated median prices for price feeds", zap.Int("num_prices", len(updatedPrices)))
}

// CalculateConvertedPrices calculates the converted prices for a given set of paths and target ticker.
// The prices utilized are the prices most recently seen by the providers. Each price is within a
// MaxPriceAge window so is safe to use.
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

	convertedPrices := make([]*big.Int, 0, len(paths.Paths))
	for _, path := range paths.Paths {
		if len(path.Operations) == 0 {
			m.logger.Error(
				"no operations",
				zap.String("ticker", target.String()),
			)

			continue
		}

		// Calculate the converted price.
		provider := path.Operations[0].Provider
		adjustedPrice, err := m.CalculateAdjustedPrice(target, path.Operations)
		if err != nil {
			m.logger.Debug(
				"failed to calculate converted price",
				zap.Error(err),
				zap.String("ticker", target.String()),
				zap.Any("conversions", path),
			)

			m.metrics.AddProviderTick(provider, target.String(), false)
			continue
		}

		convertedPrices = append(convertedPrices, adjustedPrice)
		m.logger.Debug(
			"calculated converted price",
			zap.String("ticker", target.String()),
			zap.String("price", adjustedPrice.String()),
			zap.Any("conversions", path.Operations),
		)

		m.metrics.AddProviderTick(provider, target.String(), true)
		floatPrice, _ := adjustedPrice.Float64()
		m.metrics.UpdatePrice(provider, target.String(), target.GetDecimals(), floatPrice)
	}

	return convertedPrices
}

// CalculateAdjustedPrice calculates an adjusted price for a given set of operations (if applicable).
// In particular, this assumes that every operation is either:
//
//  1. A direct conversion from the base ticker to the target ticker i.e. we want BTC/USD and
//     we have BTC/USD from a provider (e.g. Coinbase).
//  2. We need to convert the price of a given asset against the index price of an asset.
//
// In the first case, we can simply return the price of the provider. In the second case, we need
// to adjust the price by the index price of the asset. If the index price is not available, we
// return an error.
func (m *MedianAggregator) CalculateAdjustedPrice(
	target mmtypes.Ticker,
	operations []mmtypes.Operation,
) (*big.Int, error) {
	// Sanity check the number of operations. This should be [1, 2] operations.
	if len(operations) == 0 {
		return nil, fmt.Errorf("no operations")
	} else if len(operations) > mmtypes.MaxConversionOperations {
		return nil, fmt.Errorf("too many operations: %d", len(operations))
	}

	price, err := m.GetProviderPrice(operations[0])
	if err != nil {
		return nil, err
	}

	// If we have a single operation, then we can simply return the price. This implies that
	// we have a direct conversion from the base ticker to the target ticker.
	if len(operations) == 1 {
		return ScaleDownCurrencyPairPrice(target.Decimals, price)
	}

	// If we have more than one operation, then can only adjust the price using the index.
	if operations[1].Provider != mmtypes.IndexPrice {
		return nil, fmt.Errorf("expected index price but got %s", operations[1].Provider)
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
