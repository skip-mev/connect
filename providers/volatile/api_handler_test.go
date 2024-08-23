package volatile_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/volatile"
)

var (
	ticker1 = types.NewProviderTicker("foo/bar", "{}")
	ticker2 = types.NewProviderTicker("foo/baz", "{}")
)

func setupTest(t *testing.T) types.PriceAPIDataHandler {
	t.Helper()
	h := volatile.NewAPIHandler()
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
	resp := volatileHandler.ParseResponse([]types.ProviderTicker{ticker1, ticker2}, nil)
	require.Equal(t, 2, len(resp.Resolved))
	require.NotNilf(t, resp.Resolved[ticker1], "did not receive a response for ticker1")
}
