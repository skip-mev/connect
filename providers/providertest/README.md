# Provider testing

## Example

The following example can be used as a base for testing providers.

```go
package providertest_test

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/providertest"
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
	// take in a market map and filter it to output N market maps with only a single provider
	marketsPerProvider := providertest.FilterMarketMapToProviders(mm)

	// run this check for each provider (here only okx_ws)
	for provider, marketMap := range marketsPerProvider {
		ctx := context.Background()
		p, err := providertest.NewTestingOracle(ctx, provider)
		require.NoError(t, err)

		results, err := p.RunMarketMap(ctx, marketMap, providertest.DefaultProviderTestConfig())
		require.NoError(t, err)

		p.Logger.Info("results", zap.Any("results", results))
	}
}

```