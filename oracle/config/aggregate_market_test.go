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
			cfg:       config.AggregateMarketConfig{},
			expectErr: false,
		},
		{
			name: "valid config where we have the exact feed that we want",
			cfg: config.AggregateMarketConfig{
				Feeds: map[string]config.FeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
							{
								{
									CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
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
			name: "valid config where we have to aggregate to get the feed that we want",
			cfg: config.AggregateMarketConfig{
				Feeds: map[string]config.FeedConfig{
					"BITCOIN/USDT": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
					},
					"USDT/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
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
				Feeds: map[string]config.FeedConfig{
					"BITCOIN/USDT": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
					},
					"USD/USDT": {
						CurrencyPair: oracletypes.NewCurrencyPair("USD", "USDT"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
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
				Feeds: map[string]config.FeedConfig{
					"USDT/BITCOIN": {
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "BITCOIN"),
					},
					"USDT/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
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
				Feeds: map[string]config.FeedConfig{
					"BITCOIN/MOG": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "MOG"),
					},
					"MOG/USDT": {
						CurrencyPair: oracletypes.NewCurrencyPair("MOG", "USDT"),
					},
					"USDT/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
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
			name: "valid config with 1 currency pair with 1 convertable set of feeds with 5 conversions",
			cfg: config.AggregateMarketConfig{
				Feeds: map[string]config.FeedConfig{
					"BITCOIN/MOG": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "MOG"),
					},
					"MOG/PEPE": {
						CurrencyPair: oracletypes.NewCurrencyPair("MOG", "PEPE"),
					},
					"SKIP/PEPE": {
						CurrencyPair: oracletypes.NewCurrencyPair("SKIP", "PEPE"),
					},
					"SKIP/USDT": {
						CurrencyPair: oracletypes.NewCurrencyPair("SKIP", "USDT"),
					},
					"USDT/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
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
									CurrencyPair: oracletypes.NewCurrencyPair("SKIP", "PEPE"),
									Invert:       true,
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
				Feeds: map[string]config.FeedConfig{
					"USDT/BITCOIN": {
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "BITCOIN"),
					},
					"USDT/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
					},
					"BITCOIN/MOG": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "MOG"),
					},
					"MOG/USDT": {
						CurrencyPair: oracletypes.NewCurrencyPair("MOG", "USDT"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
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
			name: "invalid config with bad currency pair format",
			cfg: config.AggregateMarketConfig{
				Feeds: map[string]config.FeedConfig{
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
				Feeds: map[string]config.FeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("MOG", "USD"),
					},
				},
			},
			expectErr: true,
		},
		{
			name: "invalid config where feed in conversions is not supported",
			cfg: config.AggregateMarketConfig{
				Feeds: map[string]config.FeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
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
			name: "invalid config with no conversions",
			cfg: config.AggregateMarketConfig{
				Feeds: map[string]config.FeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions:  []config.Conversions{},
					},
				},
			},
			expectErr: true,
		},
		{
			name: "invalid config with bad conversion format",
			cfg: config.AggregateMarketConfig{
				Feeds: map[string]config.FeedConfig{
					"BITCOIN/USDT": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
					},
					"USDT/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
							{
								{
									CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", ""),
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
			name: "invalid config with mismatched outcome in conversion",
			cfg: config.AggregateMarketConfig{
				Feeds: map[string]config.FeedConfig{
					"ETHEREUM/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
							{
								{
									CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
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
