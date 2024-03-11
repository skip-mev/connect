package standard

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
)

type MedianAggregator struct {
	logger *zap.Logger
	cfg    mmtypes.MarketMap
}

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

func (m *MedianAggregator) AggregateFn() types.PriceAggregationFn {
	return func(feedsPerProvider types.AggregatedProviderPrices) types.TickerPrices {
		usdPricesPerProvider := make(types.AggregatedProviderPrices)
		for provider, prices := range feedsPerProvider {
			usdPrices := m.convertToUSD(prices)
			usdPricesPerProvider[provider] = usdPrices
		}

		// Calculate the median of the USD prices
		pricesPerTicker := make(types.TickerPrices)
		for provider, prices := range usdPricesPerProvider {
			for ticker, price := range prices {
				if _, ok := pricesPerTicker[ticker]; !ok {
					pricesPerTicker[ticker] = make([]types.Price, 0)
				}

				pricesPerTicker[ticker] = append(pricesPerTicker[ticker], price)
			}
		}
	}
}

// GetUSDTPrice returns the price of USD for each provider. This is a function of the price
// of USDT, USDC, and other preconfigured paths.
func (m *MedianAggregator) GetUSDTPrice(feedsPerProvider types.AggregatedProviderPrices) {}
