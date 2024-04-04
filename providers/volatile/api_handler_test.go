package volatile_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/volatile"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	ticker1 = mmtypes.NewTicker("foo", "bar", 1, 1)
	ticker2 = mmtypes.NewTicker("bar", "foo", 1, 1)
)

func setupTest(t *testing.T) types.PriceAPIDataHandler {
	ttpc := make(map[mmtypes.Ticker]mmtypes.ProviderConfig)
	ttpc[ticker1] = mmtypes.ProviderConfig{
		Name:           volatile.Name,
		OffChainTicker: ticker1.String(),
	}
	pmm, err := types.NewProviderMarketMap(volatile.Name, ttpc)
	require.NoError(t, err)
	h, err := volatile.NewAPIHandler(pmm)
	require.NoError(t, err)
	return h
}

func TestCreateURL(t *testing.T) {
	volatileHandler := setupTest(t)
	url, err := volatileHandler.CreateURL(nil)
	require.NoError(t, err)
	require.Equal(t, "volatile-exchange-url", url)
}

func TestParseResponse(t *testing.T) {
	volatileHandler := setupTest(t)
	resp := volatileHandler.ParseResponse([]mmtypes.Ticker{ticker1, ticker2}, nil)
	require.Equal(t, 1, len(resp.Resolved))
	require.Equal(t, 1, len(resp.UnResolved))
	require.NotNilf(t, resp.Resolved[ticker1], "did not receive a response for ticker1")
}
