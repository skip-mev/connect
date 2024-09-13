package marketmaps_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/cmd/constants/marketmaps"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

func TestMarkets(t *testing.T) {
	markets := []map[string]marketmaptypes.Market{
		marketmaps.RaydiumMarketMap.Markets,
		marketmaps.CoreMarketMap.Markets,
		marketmaps.UniswapV3BaseMarketMap.Markets,
		marketmaps.CoinGeckoMarketMap.Markets,
		marketmaps.OsmosisMarketMap.Markets,
		marketmaps.PolymarketMarketMap.Markets,
		marketmaps.ForexMarketMap.Markets,
	}
	for _, m := range markets {
		require.NotEmpty(t, m)
	}
}
