package marketmaps_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/cmd/constants/marketmaps"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func TestMarkets(t *testing.T) {
	// Unmarshal the RaydiumMarketMapJSON into RaydiumMarketMap.
	var mm mmtypes.MarketMap
	require.NoError(t, json.Unmarshal([]byte(marketmaps.RaydiumMarketMapJSON), &mm))
	require.NoError(t, mm.ValidateBasic())

	// Unmarshal the CoreMarketMapJSON into CoreMarketMap.
	var mm2 mmtypes.MarketMap
	require.NoError(t, json.Unmarshal([]byte(marketmaps.CoreMarketMapJSON), &mm2))
	require.NoError(t, mm2.ValidateBasic())

	// Unmarshal the UniswapV3BaseMarketMapJSON into UniswapV3BaseMarketMap.
	var mm3 mmtypes.MarketMap
	require.NoError(t, json.Unmarshal([]byte(marketmaps.UniswapV3BaseMarketMapJSON), &mm3))
	require.NoError(t, marketmaps.UniswapV3BaseMarketMap.ValidateBasic())

	// Unmarshal the CoinGeckoMarketMapJSON into CoinGeckoMarketMap.
	var mm4 mmtypes.MarketMap
	require.NoError(t, json.Unmarshal([]byte(marketmaps.CoinGeckoMarketMapJSON), &mm4))
	require.NoError(t, marketmaps.CoinGeckoMarketMap.ValidateBasic())

	// Unmarshal the CoinMarketCapMarketMapJSON into CoinMarketCapMarketMap.
	var mm5 mmtypes.MarketMap
	require.NoError(t, json.Unmarshal([]byte(marketmaps.CoinMarketCapMarketMapJSON), &mm5))
	require.NoError(t, marketmaps.CoinMarketCapMarketMap.ValidateBasic())

	// Unmarshal the OsmosisMarketMapJSON into OsmosisMarketMap.
	var mm6 mmtypes.MarketMap
	require.NoError(t, json.Unmarshal([]byte(marketmaps.OsmosisMarketMapJSON), &mm6))
	require.NoError(t, marketmaps.OsmosisMarketMap.ValidateBasic())
}
