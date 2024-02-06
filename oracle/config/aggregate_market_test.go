package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
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
			expectErr: true,
		},
		{
			name: "valid config where we have the exact feed that we want",
			cfg: config.AggregateMarketConfig{
				Feeds: map[string]config.FeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
					},
				},
				AggregatedFeeds: map[string][][]config.Conversion{
					"BITCOIN/USD": {
						{
							{
								CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
								Invert:       false,
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
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT", oracletypes.DefaultDecimals),
					},
					"USDT/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD", oracletypes.DefaultDecimals),
					},
				},
				AggregatedFeeds: map[string][][]config.Conversion{
					"BITCOIN/USD": {
						{
							{
								CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT", oracletypes.DefaultDecimals),
								Invert:       false,
							},
							{
								CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD", oracletypes.DefaultDecimals),
								Invert:       false,
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
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT", oracletypes.DefaultDecimals),
					},
					"USD/USDT": {
						CurrencyPair: oracletypes.NewCurrencyPair("USD", "USDT", oracletypes.DefaultDecimals),
					},
				},
				AggregatedFeeds: map[string][][]config.Conversion{
					"BITCOIN/USD": {
						{
							{
								CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT", oracletypes.DefaultDecimals),
								Invert:       false,
							},
							{
								CurrencyPair: oracletypes.NewCurrencyPair("USD", "USDT", oracletypes.DefaultDecimals),
								Invert:       true,
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
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "BITCOIN", oracletypes.DefaultDecimals),
					},
					"USDT/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD", oracletypes.DefaultDecimals),
					},
				},
				AggregatedFeeds: map[string][][]config.Conversion{
					"BITCOIN/USD": {
						{
							{
								CurrencyPair: oracletypes.NewCurrencyPair("USDT", "BITCOIN", oracletypes.DefaultDecimals),
								Invert:       true,
							},
							{
								CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD", oracletypes.DefaultDecimals),
								Invert:       false,
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
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "MOG", oracletypes.DefaultDecimals),
					},
					"MOG/USDT": {
						CurrencyPair: oracletypes.NewCurrencyPair("MOG", "USDT", oracletypes.DefaultDecimals),
					},
					"USDT/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD", oracletypes.DefaultDecimals),
					},
				},
				AggregatedFeeds: map[string][][]config.Conversion{
					"BITCOIN/USD": {
						{
							{
								CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "MOG", oracletypes.DefaultDecimals),
								Invert:       false,
							},
							{
								CurrencyPair: oracletypes.NewCurrencyPair("MOG", "USDT", oracletypes.DefaultDecimals),
								Invert:       false,
							},
							{
								CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD", oracletypes.DefaultDecimals),
								Invert:       false,
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
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "MOG", oracletypes.DefaultDecimals),
					},
					"MOG/PEPE": {
						CurrencyPair: oracletypes.NewCurrencyPair("MOG", "PEPE", oracletypes.DefaultDecimals),
					},
					"SKIP/PEPE": {
						CurrencyPair: oracletypes.NewCurrencyPair("SKIP", "PEPE", oracletypes.DefaultDecimals),
					},
					"SKIP/USDT": {
						CurrencyPair: oracletypes.NewCurrencyPair("SKIP", "USDT", oracletypes.DefaultDecimals),
					},
					"USDT/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD", oracletypes.DefaultDecimals),
					},
				},
				AggregatedFeeds: map[string][][]config.Conversion{
					"BITCOIN/USD": {
						{
							{
								CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "MOG", oracletypes.DefaultDecimals),
								Invert:       false,
							},
							{
								CurrencyPair: oracletypes.NewCurrencyPair("MOG", "PEPE", oracletypes.DefaultDecimals),
								Invert:       false,
							},
							{
								CurrencyPair: oracletypes.NewCurrencyPair("SKIP", "PEPE", oracletypes.DefaultDecimals),
								Invert:       true,
							},
							{
								CurrencyPair: oracletypes.NewCurrencyPair("SKIP", "USDT", oracletypes.DefaultDecimals),
								Invert:       false,
							},
							{
								CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD", oracletypes.DefaultDecimals),
								Invert:       false,
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
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "BITCOIN", oracletypes.DefaultDecimals),
					},
					"USDT/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD", oracletypes.DefaultDecimals),
					},
					"BITCOIN/MOG": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "MOG", oracletypes.DefaultDecimals),
					},
					"MOG/USDT": {
						CurrencyPair: oracletypes.NewCurrencyPair("MOG", "USDT", oracletypes.DefaultDecimals),
					},
				},
				AggregatedFeeds: map[string][][]config.Conversion{
					"BITCOIN/USD": {
						{
							{
								CurrencyPair: oracletypes.NewCurrencyPair("USDT", "BITCOIN", oracletypes.DefaultDecimals),
								Invert:       true,
							},
							{
								CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD", oracletypes.DefaultDecimals),
								Invert:       false,
							},
						},
						{
							{
								CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "MOG", oracletypes.DefaultDecimals),
								Invert:       false,
							},
							{
								CurrencyPair: oracletypes.NewCurrencyPair("MOG", "USDT", oracletypes.DefaultDecimals),
								Invert:       false,
							},
							{
								CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD", oracletypes.DefaultDecimals),
								Invert:       false,
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
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
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
						CurrencyPair: oracletypes.NewCurrencyPair("MOG", "USD", oracletypes.DefaultDecimals),
					},
				},
			},
			expectErr: true,
		},
		{
			name: "invalid config with insufficient aggregations",
			cfg: config.AggregateMarketConfig{
				Feeds: map[string]config.FeedConfig{
					"BITCOIN/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
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
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
					},
				},
				AggregatedFeeds: map[string][][]config.Conversion{
					"BITCOIN/USD": {
						{
							{
								CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT", oracletypes.DefaultDecimals),
								Invert:       false,
							},
							{
								CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD", oracletypes.DefaultDecimals),
								Invert:       false,
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
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
					},
				},
				AggregatedFeeds: map[string][][]config.Conversion{
					"BITCOIN/USD": {},
				},
			},
			expectErr: true,
		},
		{
			name: "invalid config with bad conversion format",
			cfg: config.AggregateMarketConfig{
				Feeds: map[string]config.FeedConfig{
					"BITCOIN/USDT": {
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT", oracletypes.DefaultDecimals),
					},
					"USDT/USD": {
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD", oracletypes.DefaultDecimals),
					},
				},
				AggregatedFeeds: map[string][][]config.Conversion{
					"BITCOIN/USD": {
						{
							{
								CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "", oracletypes.DefaultDecimals),
								Invert:       false,
							},
							{
								CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD", oracletypes.DefaultDecimals),
								Invert:       false,
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
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD", oracletypes.DefaultDecimals),
					},
				},
				AggregatedFeeds: map[string][][]config.Conversion{
					"BITCOIN/USD": {
						{
							{
								CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD", oracletypes.DefaultDecimals),
								Invert:       false,
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
				require.Error(t, err, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
