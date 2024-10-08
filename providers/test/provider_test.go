package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

var (
	usdtusd = mmtypes.Market{
		Ticker: mmtypes.Ticker{
			CurrencyPair: connecttypes.CurrencyPair{
				Base:  "USDT",
				Quote: "USD",
			},
			Decimals:         8,
			MinProviderCount: 1,
			Enabled:          true,
		},
		ProviderConfigs: []mmtypes.ProviderConfig{
			{
				Name:           "okx_ws",
				OffChainTicker: "USDC-USDT",
				Invert:         true,
			},
		},
	}

	mm = mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			usdtusd.Ticker.String(): usdtusd,
		},
	}
)

func TestProvider(t *testing.T) {
	marketsPerProvider := FilterMarketMapToProviders(mm)

	for provider, marketMap := range marketsPerProvider {
		ctx := context.Background()
		p, err := NewTestingOracle(ctx, provider)
		require.NoError(t, err)

		err = p.RunMarketMap(ctx, marketMap, DefaultProviderTestConfig())
		require.NoError(t, err)
	}

}
