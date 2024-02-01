package config_test

import (
	"testing"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestAggregateMarketConfig(t *testing.T) {
	testCases := []struct {
		name      string
		cfg       config.AggregateMarketConfig
		expectErr bool
	}{
		{
			name:      "empty config",
			cfg:       config.NewAggregateMarketConfig(),
			expectErr: false,
		},
		{
			name: "valid config with 1 currency pair with no convertable markets",
			cfg: config.AggregateMarketConfig{
				CurrencyPairs: map[string]config.AggregateCurrencyPairConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					},
				},
			},
			expectErr: false,
		},
		{
			name: "valid config with 1 currency pair with 1 convertable market",
			cfg: config.AggregateMarketConfig{
				CurrencyPairs: map[string]config.AggregateCurrencyPairConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						ConvertableMarkets: [][]config.ConvertableMarket{
							{
								{
									CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
									Invert:       false,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
					},
				},
			},
			expectErr: false,
		},
		{
			name: "valid config with 1 currency pair and a invertable market conversion (end)",
			cfg: config.AggregateMarketConfig{
				CurrencyPairs: map[string]config.AggregateCurrencyPairConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						ConvertableMarkets: [][]config.ConvertableMarket{
							{
								{
									CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
									Invert:       false,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("USD", "USDT"),
									Invert:       true,
								},
							},
						},
					},
				},
			},
			expectErr: false,
		},
		{
			name: "valid config with 1 currency pair and a invertable market conversion (start)",
			cfg: config.AggregateMarketConfig{
				CurrencyPairs: map[string]config.AggregateCurrencyPairConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						ConvertableMarkets: [][]config.ConvertableMarket{
							{
								{
									CurrencyPair: oracletypes.NewCurrencyPair("USDT", "BITCOIN"),
									Invert:       true,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
					},
				},
			},
			expectErr: false,
		},
		{
			name: "valid config with 1 currency pair with 1 convertable market with 3 conversions",
			cfg: config.AggregateMarketConfig{
				CurrencyPairs: map[string]config.AggregateCurrencyPairConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						ConvertableMarkets: [][]config.ConvertableMarket{
							{
								{
									CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "MOG"),
									Invert:       false,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("MOG", "USDT"),
									Invert:       false,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
					},
				},
			},
			expectErr: false,
		},
		{
			name: "valid config with 1 currency pair with 1 convertable market with 5 conversions",
			cfg: config.AggregateMarketConfig{
				CurrencyPairs: map[string]config.AggregateCurrencyPairConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						ConvertableMarkets: [][]config.ConvertableMarket{
							{
								{
									CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "MOG"),
									Invert:       false,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("MOG", "PEPE"),
									Invert:       false,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("PEPE", "SKIP"),
									Invert:       false,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("SKIP", "USDT"),
									Invert:       false,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
					},
				},
			},
			expectErr: false,
		},
		{
			name: "valid config with 1 currency pair with 2 convertable markets",
			cfg: config.AggregateMarketConfig{
				CurrencyPairs: map[string]config.AggregateCurrencyPairConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						ConvertableMarkets: [][]config.ConvertableMarket{
							{
								{
									CurrencyPair: oracletypes.NewCurrencyPair("USDT", "BITCOIN"),
									Invert:       true,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
							{
								{
									CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "MOG"),
									Invert:       false,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("MOG", "USDT"),
									Invert:       false,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("USD", "USDT"),
									Invert:       true,
								},
							},
						},
					},
				},
			},
			expectErr: false,
		},
		{
			name: "invalid config with bad currency pair format",
			cfg: config.AggregateMarketConfig{
				CurrencyPairs: map[string]config.AggregateCurrencyPairConfig{
					"BITCOINUSD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					},
				},
			},
			expectErr: true,
		},
		{
			name: "invalid config with mismatched currency pairs",
			cfg: config.AggregateMarketConfig{
				CurrencyPairs: map[string]config.AggregateCurrencyPairConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("MOG", "USD"),
					},
				},
			},
			expectErr: true,
		},
		{
			name: "invalid config with insufficient convertable markets",
			cfg: config.AggregateMarketConfig{
				CurrencyPairs: map[string]config.AggregateCurrencyPairConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						ConvertableMarkets: [][]config.ConvertableMarket{
							{
								{
									CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
									Invert:       false,
								},
							},
						},
					},
				},
			},
			expectErr: true,
		},
		{
			name: "invalid config with mismatched outcome in conversion",
			cfg: config.AggregateMarketConfig{
				CurrencyPairs: map[string]config.AggregateCurrencyPairConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						ConvertableMarkets: [][]config.ConvertableMarket{
							{
								{
									CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
									Invert:       false,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
					},
				},
			},
			expectErr: true,
		},
		{
			name: "invalid config with mismatched outcome in conversion (inverted)",
			cfg: config.AggregateMarketConfig{
				CurrencyPairs: map[string]config.AggregateCurrencyPairConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						ConvertableMarkets: [][]config.ConvertableMarket{
							{
								{
									CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
									Invert:       false,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
									Invert:       true,
								},
							},
						},
					},
				},
			},
			expectErr: true,
		},
		{
			name: "invalid config with mismatched outcome in conversion",
			cfg: config.AggregateMarketConfig{
				CurrencyPairs: map[string]config.AggregateCurrencyPairConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						ConvertableMarkets: [][]config.ConvertableMarket{
							{
								{
									CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
									Invert:       false,
								},
								{
									CurrencyPair: oracletypes.NewCurrencyPair("USD", "USDT"),
									Invert:       false,
								},
							},
						},
					},
				},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
