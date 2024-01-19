package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

func TestProviderMarketConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      config.ProviderMarketConfig
		expectedErr bool
	}{}

	oracleCfg := config.OracleConfig{
		UpdateInterval: time.Second,
		Providers: []config.ProviderConfig{
			config.ProviderConfig{
				API: config.APIConfig{
					Enabled:    true,
					Timeout:    time.Second,
					Interval:   time.Second,
					MaxQueries: 1,
				},
				Name: "test",
				MarketConfig: config.ProviderMarketConfig{
					Name: "test",
					CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
						"BITCOIN/USDT": {
							Ticker:       "BTCUSDT",
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
						},
						"ETHEREUM/USDT": {
							Ticker:       "ETHUSDT",
							CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
						},
					},
				},
			},
			{
				WebSocket: config.WebSocketConfig{
					Enabled:             true,
					MaxBufferSize:       1,
					ReconnectionTimeout: time.Second,
				},
				Name: "test2",
				MarketConfig: config.ProviderMarketConfig{
					Name: "test2",
					CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
						"BITCOIN/USDT": {
							Ticker:       "BTC-USDT",
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
						},
						"ETHEREUM/USDT": {
							Ticker:       "ETH-USDT",
							CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
						},
					},
				},
			},
		},
	}

	f, err := os.Create("config.toml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(oracleCfg); err != nil {
		panic(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.ValidateBasic()
			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
