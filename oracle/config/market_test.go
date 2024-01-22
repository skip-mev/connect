package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

func TestMarketConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      config.MarketConfig
		expectedErr bool
	}{
		{
			name: "good config",
			config: config.MarketConfig{
				Name: "test",
				CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
					"BITCOIN/USD": {
						Ticker:       "BTC/USD",
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "bad config with no name",
			config: config.MarketConfig{
				CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
					"BITCOIN/USD": {
						Ticker:       "BTC/USD",
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "bad config with bad currency pair format",
			config: config.MarketConfig{
				Name: "test",
				CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
					"BITCOINUSD": {
						Ticker:       "BTC/USD",
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "bad config with bad currency pair",
			config: config.MarketConfig{
				Name: "test",
				CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
					"BITCOIN/USD": {
						Ticker:       "BTC/USD",
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", ""),
					},
				},
			},
			expectedErr: true,
		},
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

func TestInvertMarketConfig(t *testing.T) {
	cfg := config.MarketConfig{
		Name: "test",
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"BITCOIN/USD": {
				Ticker:       "BTC/USD",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"BITCOIN/EURO": {
				Ticker:       "BTC/EURO",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "EURO"),
			},
		},
	}

	invertedCfg := cfg.Invert()
	require.Equal(t, "test", invertedCfg.Name)
	require.Equal(t, len(cfg.CurrencyPairToMarketConfigs), len(invertedCfg.MarketToCurrencyPairConfigs))

	for ticker, marketConfig := range cfg.CurrencyPairToMarketConfigs {
		invertedMarketConfig, ok := invertedCfg.MarketToCurrencyPairConfigs[marketConfig.Ticker]
		require.True(t, ok)
		require.Equal(t, ticker, invertedMarketConfig.CurrencyPair.ToString())
	}
}
