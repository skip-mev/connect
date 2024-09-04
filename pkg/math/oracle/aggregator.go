package oracle

import (
	"fmt"
	"math/big"
	"sync"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle"
	oraclemetrics "github.com/skip-mev/connect/v2/oracle/metrics"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

var _ oracle.PriceAggregator = &IndexPriceAggregator{}

// IndexPriceAggregator is an aggregator that calculates the median price for each ticker,
// resolved from a predefined set of conversion markets. A conversion market is a set of
// markets that can be used to convert the prices of a set of tickers to a common ticker.
// These are defined in the market map configuration.
type IndexPriceAggregator struct {
	mtx     sync.Mutex
	logger  *zap.Logger
	cfg     mmtypes.MarketMap
	metrics oraclemetrics.Metrics

	// indexPrices cache the median prices for each ticker. These are unscaled prices.
	indexPrices types.Prices
	// scaledPrices cache the scaled prices for each ticker. These are the prices that can be
	// consumed by consumers.
	scaledPrices types.Prices
	// providerPrices cache the unscaled prices for each provider. These are indexed by
	// provider -> offChainTicker -> price.
	providerPrices map[string]types.Prices
}

// NewIndexPriceAggregator returns a new Index Price Aggregator.
func NewIndexPriceAggregator(
	logger *zap.Logger,
	cfg mmtypes.MarketMap,
	metrics oraclemetrics.Metrics,
) (*IndexPriceAggregator, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if metrics == nil {
		logger.Warn("metrics is nil; using a no-op metrics implementation")
		metrics = oraclemetrics.NewNopMetrics()
	}

	return &IndexPriceAggregator{
		logger:         logger.With(zap.String("process", "index_price_aggregator")),
		cfg:            cfg,
		metrics:        metrics,
		indexPrices:    make(types.Prices),
		scaledPrices:   make(types.Prices),
		providerPrices: make(map[string]types.Prices),
	}, nil
}

// AggregatePrices implements the aggregate function for the median price calculation. Specifically, this
// aggregation function aggregates the prices seen by each provider by first converting each price to a
// common ticker and then calculating the median of the converted prices. Prices are converted either
//
//  1. Directly from the base ticker to the target ticker. i.e. I have BTC/USD and I want BTC/USD.
//  2. Using the index price of an asset. i.e. I have BTC/USDT and I want BTC/USD. I can convert
//     BTC/USDT to BTC/USD using the index price of USDT/USD.
//
// The index price cache contains the previously calculated median prices.
func (m *IndexPriceAggregator) AggregatePrices() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	indexPrices := make(types.Prices)
	scaledPrices := make(types.Prices)

	var missingPrices []string

	for ticker, market := range m.cfg.Markets {
		if !market.Ticker.Enabled {
			m.logger.Debug("skipping disabled market", zap.Any("market", market))
			continue
		}

		// Get the converted prices for set of convertible markets.
		// ex. BTC/USDT * Index USDT/USD = BTC/USD
		//     BTC/USDC * Index USDC/USD = BTC/USD
		target := market.Ticker
		convertedPrices := m.CalculateConvertedPrices(market)
		m.metrics.AddProviderCountForMarket(target.String(), len(convertedPrices))

		// We need to have at least the minimum number of providers to calculate the median.
		if len(convertedPrices) < int(target.MinProviderCount) { //nolint:gosec
			missingPrices = append(missingPrices, ticker)
			m.logger.Debug(
				"insufficient amount of converted prices",
				zap.String("target_ticker", ticker),
				zap.Int("num_converted_prices", len(convertedPrices)),
				zap.Any("converted_prices", convertedPrices),
				zap.Int("min_provider_count", int(target.MinProviderCount)), //nolint:gosec
			)

			continue
		}

		// Take the median of the converted prices. This takes the average of the middle two
		// prices if the number of prices is even.
		price := math.CalculateMedian(convertedPrices)
		indexPrices[target.String()] = new(big.Float).Copy(price)

		// Scale the price to the target ticker's decimals.
		scaledPrices[target.String()] = math.ScaleBigFloat(new(big.Float).Copy(price), target.Decimals)

		m.logger.Debug(
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
	m.logger.Debug("calculated median prices for price feeds", zap.Int("num_prices", len(indexPrices)))
	m.metrics.MissingPrices(missingPrices)
	if len(missingPrices) > 0 {
		m.logger.Info("failed to calculate prices for price feeds", zap.Strings("missing_prices", missingPrices))
	}
	m.indexPrices = indexPrices
	m.scaledPrices = scaledPrices
}

// CalculateConvertedPrices calculates the converted prices for a given set of paths and target ticker.
// The prices utilized are the prices most recently seen by the providers. Each price is within a
// MaxPriceAge window so is safe to use.
func (m *IndexPriceAggregator) CalculateConvertedPrices(
	market mmtypes.Market,
) []*big.Float {
	m.logger.Debug("calculating converted prices", zap.String("ticker", market.Ticker.String()))
	if len(market.ProviderConfigs) == 0 {
		m.logger.Error(
			"no conversion paths",
			zap.String("target_ticker", market.Ticker.String()),
		)

		return nil
	}

	convertedPrices := make([]*big.Float, 0, len(market.ProviderConfigs))
	for _, cfg := range market.ProviderConfigs {
		// Calculate the converted price.
		adjustedPrice, err := m.CalculateAdjustedPrice(cfg)
		if err != nil {
			m.logger.Debug(
				"failed to calculate converted price",
				zap.Error(err),
				zap.String("target_ticker", market.Ticker.String()),
				zap.Any("provider", cfg.Name),
			)

			m.metrics.AddProviderTick(cfg.Name, market.Ticker.String(), false)
			continue
		}

		convertedPrices = append(convertedPrices, adjustedPrice)
		m.logger.Debug(
			"calculated converted price",
			zap.String("target_ticker", market.Ticker.String()),
			zap.String("price", adjustedPrice.String()),
			zap.Any("provider", cfg.Name),
		)

		m.metrics.AddProviderTick(cfg.Name, market.Ticker.String(), true)
		floatPrice, _ := adjustedPrice.Float64()
		m.metrics.UpdatePrice(cfg.Name, market.Ticker.String(), market.Ticker.GetDecimals(), floatPrice)
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
func (m *IndexPriceAggregator) CalculateAdjustedPrice(
	cfg mmtypes.ProviderConfig,
) (*big.Float, error) {
	price, err := m.GetProviderPrice(cfg)
	if err != nil {
		return nil, err
	}

	if cfg.NormalizeByPair == nil {
		return price, nil
	}

	normalizeByIndexPrice, err := m.GetIndexPrice(*cfg.NormalizeByPair)
	if err != nil {
		return nil, err
	}

	// Make sure that the price is adjusted by the market price.
	return new(big.Float).Mul(price, normalizeByIndexPrice), nil
}
