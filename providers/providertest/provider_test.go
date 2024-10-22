package providertest_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/providers/providertest"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

func TestUniswapMarkets(t *testing.T) {
	mm, err := mmtypes.ReadMarketMapFromFile("./output.json")
	require.NoError(t, err)

	providerMM := providertest.FilterMarketMapToProviders(mm)
	for provider, mm := range providerMM {
		if provider == "uniswapv3_api-ethereum" {
			ctx := context.Background()
			o, err := providertest.NewTestingOracle(ctx, []string{provider})
			require.NoError(t, err)

			res, err := o.RunMarketMap(ctx, mm, providertest.DefaultProviderTestConfig())
			require.NoError(t, err)

			t.Log(res)

		}
	}
}
