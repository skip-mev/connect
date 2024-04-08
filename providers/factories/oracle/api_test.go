package oracle_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	oracleconfig "github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/oracle/types"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/factories/oracle"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func TestAPIQueryHandlerFactory(t *testing.T) {
	t.Run("if batch-size * max-queries < len(tickers) - fail", func(t *testing.T) {
		cfg := oracleconfig.ProviderConfig{
			Name: "test",
			API: oracleconfig.APIConfig{
				BatchSize:        10,
				MaxQueries:       1,
				Name:             "test",
				Enabled:          true,
				Interval:         1,
				Timeout:          1,
				ReconnectTimeout: 1,
				URL:              "http://test.com",
			},
			Type: "test",
		}
		mm := oracletypes.ProviderMarketMap{
			TickerConfigs: make(oracletypes.TickerToProviderConfig),
		}

		// add 11 tickers
		for i := 0; i < 11; i++ {
			mm.TickerConfigs[mmtypes.Ticker{
				CurrencyPair: slinkytypes.NewCurrencyPair("BTC", fmt.Sprintf("USD%d", i)),
			}] = mmtypes.ProviderConfig{}
		}

		_, err := oracle.APIQueryHandlerFactory(nil, cfg, nil, mm)
		require.Error(t, err)
		require.Equal(t, "number of tickers to fetch for: 11 is greater than the batch-size (10) * max-queries (1)", err.Error())
	})

	t.Run("if batch-size * max-queries < len(tickers), and batch-Size is inferred - fail", func(t *testing.T) {
		cfg := oracleconfig.ProviderConfig{
			Name: "test",
			API: oracleconfig.APIConfig{
				MaxQueries:       10,
				Name:             "test",
				Enabled:          true,
				Interval:         1,
				Timeout:          1,
				ReconnectTimeout: 1,
				URL:              "http://test.com",
			},
			Type: "test",
		}
		mm := oracletypes.ProviderMarketMap{
			TickerConfigs: make(oracletypes.TickerToProviderConfig),
		}

		// add 11 tickers
		for i := 0; i < 11; i++ {
			mm.TickerConfigs[mmtypes.Ticker{
				CurrencyPair: slinkytypes.NewCurrencyPair("BTC", fmt.Sprintf("USD%d", i)),
			}] = mmtypes.ProviderConfig{}
		}

		_, err := oracle.APIQueryHandlerFactory(nil, cfg, nil, mm)
		require.Error(t, err)
		require.Equal(t, "number of tickers to fetch for: 11 is greater than the batch-size (1) * max-queries (10)", err.Error())
	})
}
