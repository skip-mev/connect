package oracle

import (
	"fmt"
	"math/big"
	"sync"

	"go.uber.org/zap"

	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
)

var _ types.PriceAggregator = &MedianAggregator{}

// MedianAggregator is an aggregator that calculates the median price for each ticker,
// resolved from a predefined set of conversion markets. A conversion market is a set of
// markets that can be used to convert the prices of a set of tickers to a common ticker.
// These are defined in the market map configuration.
type MedianAggregator struct {
	mtx     sync.Mutex
	logger  *zap.Logger
	cfg     mmtypes.MarketMap
	metrics oraclemetrics.Metrics

	// indexPrices cache the median prices for each ticker. These are unscaled prices.
	indexPrices types.AggregatorPrices
	// scaledPrices cache the scaled prices for each ticker. These are the prices that can be
	// consumed by external providers.
	scaledPrices types.AggregatorPrices
	// providerPrices cache the unscaled prices for each provider. These are indexed by
	// provider -> offChainTicker -> price.
	providerPrices types.AggregatedProviderPrices
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

	if metrics == nil {
		logger.Warn("metrics is nil; using a no-op metrics implementation")
		metrics = oraclemetrics.NewNopMetrics()
	}

	return &MedianAggregator{
		logger:         logger,
		cfg:            cfg,
		metrics:        metrics,
		indexPrices:    make(types.AggregatorPrices),
		scaledPrices:   make(types.AggregatorPrices),
		providerPrices: make(types.AggregatedProviderPrices),
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
	m.mtx.Lock()
	defer m.mtx.Unlock()

	indexPrices := make(types.AggregatorPrices)
	scaledPrices := make(types.AggregatorPrices)

	for ticker, market := range m.cfg.Markets {
		// Get the converted prices for set of convertable markets.
		// ex. BTC/USDT * Index USDT/USD = BTC/USD
		//     BTC/USDC * Index USDC/USD = BTC/USD
		target := market.Ticker
		convertedPrices := m.CalculateConvertedPrices(target, market.ProviderConfigs)

		// We need to have at least the minimum number of providers to calculate the median.
		if len(convertedPrices) < int(target.MinProviderCount) {
			m.logger.Error(
				"insufficient amount of converted prices",
				zap.String("target_ticker", ticker),
				zap.Int("num_converted_prices", len(convertedPrices)),
				zap.Any("converted_prices", convertedPrices),
				zap.Int("min_provider_count", int(target.MinProviderCount)),
			)

			continue
		}

		// Take the median of the converted prices. This takes the average of the middle two
		// prices if the number of prices is even.
		price := math.CalculateMedian(convertedPrices)
		indexPrices[target.String()] = new(big.Float).Copy(price)

		// Scale the price to the target ticker's decimals.
		scaledPrices[target.String()] = math.ScaleBigFloat(new(big.Float).Copy(price), target.Decimals)

		m.logger.Info(
			"calculated median price",
			zap.String("target_ticker", ticker),

			zap.String("unscaled_price", indexPrices[target.String()].String()),
			zap.String("scaled_price", scaledPrices[target.String()].String()),
			zap.Any("converted_prices", convertedPrices),
		)
		floatPrice, _ := price.Float64()
		m.metrics.AddTickerTick(target.String())
		m.metrics.UpdateAggregatePrice(target.String(), target.GetDecimals(), floatPrice)
	}

	// Update the aggregated data. These prices are going to be used as the index prices the
	// next time we calculate prices.
	m.logger.Info("calculated median prices for price feeds", zap.Int("num_prices", len(indexPrices)))
	m.indexPrices = indexPrices
	m.scaledPrices = scaledPrices
}

// CalculateConvertedPrices calculates the converted prices for a given set of paths and target ticker.
// The prices utilized are the prices most recently seen by the providers. Each price is within a
// MaxPriceAge window so is safe to use.
func (m *MedianAggregator) CalculateConvertedPrices(
	target mmtypes.Ticker,
	providers []mmtypes.ProviderConfig,
) []*big.Float {
	m.logger.Debug("calculating converted prices", zap.String("ticker", target.String()))
	if len(providers) == 0 {
		m.logger.Error(
			"no conversion paths",
			zap.String("target_ticker", target.String()),
		)

		return nil
	}

	convertedPrices := make([]*big.Float, 0, len(providers))
	for _, cfg := range providers {
		// Calculate the converted price.
		adjustedPrice, err := m.CalculateAdjustedPrice(cfg)
		if err != nil {
			m.logger.Debug(
				"failed to calculate converted price",
				zap.Error(err),
				zap.String("target_ticker", target.String()),
				zap.Any("provider", cfg.Name),
			)

			m.metrics.AddProviderTick(cfg.Name, target.String(), false)
			continue
		}

		convertedPrices = append(convertedPrices, adjustedPrice)
		m.logger.Debug(
			"calculated converted price",
			zap.String("target_ticker", target.String()),
			zap.String("price", adjustedPrice.String()),
			zap.Any("provider", cfg.Name),
		)

		m.metrics.AddProviderTick(cfg.Name, target.String(), true)
		floatPrice, _ := adjustedPrice.Float64()
		m.metrics.UpdatePrice(cfg.Name, target.String(), target.GetDecimals(), floatPrice)
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
	cfg mmtypes.ProviderConfig,
) (*big.Float, error) {
	price, err := m.GetProviderPrice(cfg)
	if err != nil {
		return nil, err
	}

	if cfg.NormalizeByPair == nil {
		return price, nil
	}

	adjustableByIndexPrice, err := m.GetIndexPrice(*cfg.NormalizeByPair)
	if err != nil {
		return nil, err
	}

	// Make sure that the price is adjusted by the market price.
	return new(big.Float).Mul(price, adjustableByIndexPrice), nil
}
