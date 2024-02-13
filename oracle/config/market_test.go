package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
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
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
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
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "bad config with poorly formatted currency pair key",
			config: config.MarketConfig{
				Name: "test",
				CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
					"BITCOINUSD": {
						Ticker:       "BTC/USD",
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "bad config with bad currency pair in config",
			config: config.MarketConfig{
				Name: "test",
				CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
					"BITCOIN/USD": {
						Ticker:       "BTC/USD",
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", ""),
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

				// Validate that the config was inverted correctly.
				require.Equal(t, len(tc.config.CurrencyPairToMarketConfigs), len(tc.config.TickerToMarketConfigs))
				for _, marketConfig := range tc.config.CurrencyPairToMarketConfigs {
					require.Contains(t, tc.config.TickerToMarketConfigs, marketConfig.Ticker)
					require.Equal(t, marketConfig, tc.config.TickerToMarketConfigs[marketConfig.Ticker])
				}
			}
		})
	}
}
