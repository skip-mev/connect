package constants_test

import (
	"encoding/json"
	"testing"

	"github.com/skip-mev/slinky/cmd/constants"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/stretchr/testify/require"
)

func TestMarkets(t *testing.T) {
	// Unmarshal the RaydiumMarketMapJSON into RaydiumMarketMap.
	var mm mmtypes.MarketMap
	require.NoError(t, json.Unmarshal([]byte(constants.RaydiumMarketMapJSON), &mm))
	require.NoError(t, mm.ValidateBasic())

	// Unmarshal the CoreMarketMapJSON into CoreMarketMap.
	var mm2 mmtypes.MarketMap
	require.NoError(t, json.Unmarshal([]byte(constants.CoreMarketMapJSON), &mm2))
	require.NoError(t, mm2.ValidateBasic())
}
