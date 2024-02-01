package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

func TestProviderConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      config.ProviderConfig
		expectedErr bool
	}{
		{
			name: "good API config",
			config: config.ProviderConfig{
				API: config.APIConfig{
					Enabled:    true,
					Timeout:    time.Second,
					Interval:   time.Second,
					MaxQueries: 1,
					Name:       "test",
					Atomic:     true,
					URL:        "http://test.com",
				},
				Name: "test",
				Market: config.MarketConfig{
					Name: "test",
					CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
						"BITCOIN/USD/8": {
							Ticker:       "BTC/USD",
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
						},
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "good websocket config",
			config: config.ProviderConfig{
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
					MaxReadErrorCount:   config.DefaultMaxReadErrorCount,
				},
				Name: "test",
				Market: config.MarketConfig{
					Name: "test",
					CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
						"BITCOIN/USD/8": {
							Ticker:       "BTC/USD",
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
						},
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "no name",
			config: config.ProviderConfig{
				API: config.APIConfig{
					Enabled:    true,
					Timeout:    time.Second,
					Interval:   time.Second,
					MaxQueries: 1,
					Name:       "test",
					Atomic:     true,
					URL:        "http://test.com",
				},
				Market: config.MarketConfig{
					Name: "test",
					CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
						"BITCOIN/USD/8": {
							Ticker:       "BTC/USD",
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "no API or websocket config",
			config: config.ProviderConfig{
				Name: "test",
				Market: config.MarketConfig{
					Name: "test",
					CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
						"BITCOIN/USD/8": {
							Ticker:       "BTC/USD",
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "both API and websocket config",
			config: config.ProviderConfig{
				API: config.APIConfig{
					Enabled:    true,
					Timeout:    time.Second,
					Interval:   time.Second,
					MaxQueries: 1,
					Name:       "test",
					Atomic:     true,
					URL:        "http://test.com",
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
					MaxReadErrorCount:   config.DefaultMaxReadErrorCount,
				},
				Name: "test",
				Market: config.MarketConfig{
					Name: "test",
					CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
						"BITCOIN/USD/8": {
							Ticker:       "BTC/USD",
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "bad API config",
			config: config.ProviderConfig{
				API: config.APIConfig{
					Enabled:    true,
					Timeout:    2 * time.Second,
					Interval:   time.Second,
					MaxQueries: 1,
				},
				Name: "test",
				Market: config.MarketConfig{
					Name: "test",
					CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
						"BITCOIN/USD/8": {
							Ticker:       "BTC/USD",
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "bad websocket config",
			config: config.ProviderConfig{
				WebSocket: config.WebSocketConfig{
					Enabled:             true,
					ReconnectionTimeout: 2 * time.Second,
				},
				Name: "test",
				Market: config.MarketConfig{
					Name: "test",
					CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
						"BITCOIN/USD/8": {
							Ticker:       "BTC/USD",
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "mismatch names between provider and ws config",
			config: config.ProviderConfig{
				WebSocket: config.WebSocketConfig{
					Enabled:             true,
					MaxBufferSize:       1,
					ReconnectionTimeout: time.Second,
					WSS:                 "wss://test.com",
					Name:                "test",
				},
				Name: "test2",
				Market: config.MarketConfig{
					Name: "test2",
					CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
						"BITCOIN/USD/8": {
							Ticker:       "BTC/USD",
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "mismatch names between provider and api config",
			config: config.ProviderConfig{
				API: config.APIConfig{
					Enabled:    true,
					Timeout:    time.Second,
					Interval:   time.Second,
					MaxQueries: 1,
					Name:       "test",
					Atomic:     true,
					URL:        "http://test.com",
				},
				Name: "test2",
				Market: config.MarketConfig{
					Name: "test2",
					CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
						"BITCOIN/USD/8": {
							Ticker:       "BTC/USD",
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "bad market config",
			config: config.ProviderConfig{
				API: config.APIConfig{
					Enabled:    true,
					Timeout:    time.Second,
					Interval:   time.Second,
					MaxQueries: 1,
					Name:       "test",
					Atomic:     true,
					URL:        "http://test.com",
				},
				Name:   "test",
				Market: config.MarketConfig{},
			},
			expectedErr: true,
		},
		{
			name: "mismatch names between provider and market config",
			config: config.ProviderConfig{
				API: config.APIConfig{
					Enabled:    true,
					Timeout:    time.Second,
					Interval:   time.Second,
					MaxQueries: 1,
					Name:       "test",
					Atomic:     true,
					URL:        "http://test.com",
				},
				Name: "test",
				Market: config.MarketConfig{
					Name: "test2",
					CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
						"BITCOIN/USD/8": {
							Ticker:       "BTC/USD",
							CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
						},
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
