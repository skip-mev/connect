package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
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
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
							{
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
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
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
					},
					"USDT/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
							{
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
									Invert:       false,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
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
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
					},
					"USD/USDT": {
						CurrencyPair: slinkytypes.NewCurrencyPair("USD", "USDT"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
							{
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
									Invert:       false,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("USD", "USDT"),
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
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "BITCOIN"),
					},
					"USDT/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
							{
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "BITCOIN"),
									Invert:       true,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
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
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "MOG"),
					},
					"MOG/USDT": {
						CurrencyPair: slinkytypes.NewCurrencyPair("MOG", "USDT"),
					},
					"USDT/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
							{
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "MOG"),
									Invert:       false,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("MOG", "USDT"),
									Invert:       false,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
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
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "MOG"),
					},
					"MOG/PEPE": {
						CurrencyPair: slinkytypes.NewCurrencyPair("MOG", "PEPE"),
					},
					"SKIP/PEPE": {
						CurrencyPair: slinkytypes.NewCurrencyPair("SKIP", "PEPE"),
					},
					"SKIP/USDT": {
						CurrencyPair: slinkytypes.NewCurrencyPair("SKIP", "USDT"),
					},
					"USDT/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
							{
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "MOG"),
									Invert:       false,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("MOG", "PEPE"),
									Invert:       false,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("SKIP", "PEPE"),
									Invert:       true,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("SKIP", "USDT"),
									Invert:       false,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
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
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "BITCOIN"),
					},
					"USDT/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
					},
					"BITCOIN/MOG": {
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "MOG"),
					},
					"MOG/USDT": {
						CurrencyPair: slinkytypes.NewCurrencyPair("MOG", "USDT"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
							{
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "BITCOIN"),
									Invert:       true,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
							{
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "MOG"),
									Invert:       false,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("MOG", "USDT"),
									Invert:       false,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
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
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
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
						CurrencyPair: slinkytypes.NewCurrencyPair("MOG", "USD"),
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
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
							{
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "BITCOIN"),
									Invert:       true,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
							{
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "MOG"),
									Invert:       false,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("MOG", "USDT"),
									Invert:       false,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
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
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
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
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
					},
					"USDT/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
							{
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", ""),
									Invert:       false,
								},
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
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
						CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
					},
				},
				AggregatedFeeds: map[string]config.AggregateFeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
						Conversions: []config.Conversions{
							{
								{
									CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
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
