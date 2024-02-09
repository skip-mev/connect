package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
)

func TestOracleConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      config.OracleConfig
		expectedErr bool
	}{
		{
			name: "good config",
			config: config.OracleConfig{
				UpdateInterval: time.Second,
				Providers: []config.ProviderConfig{
					{
						Name: "test",
						Market: config.MarketConfig{
							Name: "test",
							CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
								"BITCOIN/USD": {
									Ticker:       "BTC/USD",
									CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
								},
							},
						},
						WebSocket: config.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							WSS:                 "wss://test.com",
							Name:                "test",
							ReadBufferSize:      config.DefaultReadBufferSize,
							WriteBufferSize:     config.DefaultWriteBufferSize,
							HandshakeTimeout:    config.DefaultHandshakeTimeout,
							EnableCompression:   config.DefaultEnableCompression,
							ReadTimeout:         config.DefaultReadTimeout,
							WriteTimeout:        config.DefaultWriteTimeout,
						},
					},
				},
				Market: config.AggregateMarketConfig{
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
			},
			expectedErr: false,
		},
		{
			name:        "bad config with no update interval",
			config:      config.OracleConfig{},
			expectedErr: true,
		},
		{
			name: "bad config with bad currency pair format",
			config: config.OracleConfig{
				UpdateInterval: time.Second,
				Providers: []config.ProviderConfig{
					{
						Name: "test",
						Market: config.MarketConfig{
							Name: "test",
							CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
								"BITCOIN/USD": {
									Ticker:       "BTC/USD",
									CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
								},
							},
						},
						WebSocket: config.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							WSS:                 "wss://test.com",
							Name:                "test",
							ReadBufferSize:      config.DefaultReadBufferSize,
							WriteBufferSize:     config.DefaultWriteBufferSize,
							HandshakeTimeout:    config.DefaultHandshakeTimeout,
							EnableCompression:   config.DefaultEnableCompression,
							ReadTimeout:         config.DefaultReadTimeout,
							WriteTimeout:        config.DefaultWriteTimeout,
						},
					},
				},
				Market: config.AggregateMarketConfig{
					Feeds: map[string]config.FeedConfig{
						"BITCOINUSD": {
							CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "bad config with bad metrics",
			config: config.OracleConfig{
				UpdateInterval: time.Second,
				Providers: []config.ProviderConfig{
					{
						Name: "test",
						Market: config.MarketConfig{
							Name: "test",
							CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
								"BITCOIN/USD": {
									Ticker:       "BTC/USD",
									CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
								},
							},
						},
						WebSocket: config.WebSocketConfig{
							Enabled:             true,
							MaxBufferSize:       1,
							ReconnectionTimeout: time.Second,
							WSS:                 "wss://test.com",
							Name:                "test",
							ReadBufferSize:      config.DefaultReadBufferSize,
							WriteBufferSize:     config.DefaultWriteBufferSize,
							HandshakeTimeout:    config.DefaultHandshakeTimeout,
							EnableCompression:   config.DefaultEnableCompression,
							ReadTimeout:         config.DefaultReadTimeout,
							WriteTimeout:        config.DefaultWriteTimeout,
						},
					},
				},
				Market: config.AggregateMarketConfig{
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
										CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
										Invert:       false,
									},
								},
							},
						},
					},
				},
				Metrics: config.MetricsConfig{
					Enabled: true,
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
