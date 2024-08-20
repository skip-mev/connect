package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/types"
	pkgtypes "github.com/skip-mev/connect/v2/pkg/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

func TestProviderTickersFromMarketMap(t *testing.T) {
	cases := []struct {
		name     string
		provider string
		market   mmtypes.MarketMap
		expected []types.ProviderTicker
		err      bool
	}{
		{
			name:     "empty market map",
			provider: "test",
			market:   mmtypes.MarketMap{},
			expected: nil,
			err:      false,
		},
		{
			name:     "single disabled market",
			provider: "test",
			market: mmtypes.MarketMap{
				Markets: map[string]mmtypes.Market{
					"BTC/USD": {
						Ticker: mmtypes.NewTicker("BTC", "USD", 8, 1, false),
						ProviderConfigs: []mmtypes.ProviderConfig{
							{
								Name:           "test",
								OffChainTicker: "BTC/USDT",
								Metadata_JSON:  "{}",
							},
						},
					},
				},
			},
			expected: []types.ProviderTicker{},
			err:      false,
		},
		{
			name:     "single market for the provider",
			provider: "test",
			market: mmtypes.MarketMap{
				Markets: map[string]mmtypes.Market{
					"BTC/USD": {
						Ticker: mmtypes.NewTicker("BTC", "USD", 8, 1, true),
						ProviderConfigs: []mmtypes.ProviderConfig{
							{
								Name:           "test",
								OffChainTicker: "BTC/USDT",
								Metadata_JSON:  "{}",
							},
						},
					},
				},
			},
			expected: []types.ProviderTicker{
				types.NewProviderTicker(
					"BTC/USDT",
					"{}",
				),
			},
			err: false,
		},
		{
			name:     "provider has no configs in market map",
			provider: "test",
			market: mmtypes.MarketMap{
				Markets: map[string]mmtypes.Market{
					"BTC/USD": {
						Ticker: mmtypes.NewTicker("BTC", "USD", 8, 1, true),
						ProviderConfigs: []mmtypes.ProviderConfig{
							{
								Name:           "other",
								OffChainTicker: "BTC/USDT",
								Metadata_JSON:  "{}",
							},
						},
					},
				},
			},
			expected: nil,
			err:      false,
		},
		{
			name:     "duplicate markets for the provider",
			provider: "test",
			market: mmtypes.MarketMap{
				Markets: map[string]mmtypes.Market{
					"ETH/USD": {
						Ticker: mmtypes.NewTicker("ETH", "USD", 8, 1, true),
						ProviderConfigs: []mmtypes.ProviderConfig{
							{
								Name:           "test",
								OffChainTicker: "ETH/USDT",
								NormalizeByPair: &pkgtypes.CurrencyPair{
									Base:  "USDT",
									Quote: "USD",
								},
								Metadata_JSON: "{}",
							},
						},
					},
					"USDT/USD": {
						Ticker: mmtypes.NewTicker("USDT", "USD", 8, 1, true),
						ProviderConfigs: []mmtypes.ProviderConfig{
							{
								Name:           "test",
								OffChainTicker: "ETH/USDT",
								Invert:         true,
								NormalizeByPair: &pkgtypes.CurrencyPair{
									Base:  "ETH",
									Quote: "USD",
								},
								Metadata_JSON: "{}",
							},
						},
					},
				},
			},
			expected: []types.ProviderTicker{
				types.NewProviderTicker(
					"ETH/USDT",
					"{}",
				),
			},
			err: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := types.ProviderTickersFromMarketMap(tc.provider, tc.market)
			if tc.err {
				require.Error(t, err)
				return
			}

			expectedCache := make(map[types.ProviderTicker]struct{})
			for _, ticker := range tc.expected {
				expectedCache[ticker] = struct{}{}
			}
			actualCache := make(map[types.ProviderTicker]struct{})
			for _, ticker := range actual {
				actualCache[ticker] = struct{}{}
			}
			require.Equal(t, expectedCache, actualCache)
		})
	}
}
