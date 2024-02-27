package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func TestProviderMarketMap(t *testing.T) {
	testCases := []struct {
		name         string
		providerName string
		configs      types.TickerToProviderConfig
		expErr       bool
	}{
		{
			name:         "empty configs",
			providerName: "test",
			configs:      types.TickerToProviderConfig{},
			expErr:       false,
		},
		{
			name:         "empty provider name",
			providerName: "",
			configs: types.TickerToProviderConfig{
				constants.BITCOIN_USD: {
					Name:           "test",
					OffChainTicker: "BTC-USD",
				},
			},
			expErr: true,
		},
		{
			name:         "empty off-chain ticker",
			providerName: "test",
			configs: types.TickerToProviderConfig{
				constants.BITCOIN_USD: {
					Name:           "test",
					OffChainTicker: "",
				},
			},
			expErr: true,
		},
		{
			name:         "invalid ticker",
			providerName: "test",
			configs: types.TickerToProviderConfig{
				mmtypes.NewTicker("BTC", "USD", 8, 0): {
					Name:           "test",
					OffChainTicker: "BTC-USD",
				},
			},
			expErr: true,
		},
		{
			name:         "valid configs",
			providerName: "test",
			configs: types.TickerToProviderConfig{
				constants.BITCOIN_USD: {
					Name:           "test",
					OffChainTicker: "BTC-USD",
				},
				constants.BITCOIN_USDC: {
					Name:           "test",
					OffChainTicker: "BTC-USDC",
				},
			},
			expErr: false,
		},
		{
			name:         "mismatch provider name and config",
			providerName: "test",
			configs: types.TickerToProviderConfig{
				constants.BITCOIN_USD: {
					Name:           "invalid",
					OffChainTicker: "BTC-USD",
				},
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pmap, err := types.NewProviderMarketMap(tc.providerName, tc.configs)
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			err = pmap.ValidateBasic()
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProviderMarketMapFromMarketMap(t *testing.T) {
	testCases := []struct {
		name         string
		marketMap    mmtypes.MarketMap
		providerName string
		expectedMap  types.ProviderMarketMap
		expErr       bool
	}{
		{
			name:         "empty market map",
			marketMap:    mmtypes.MarketMap{},
			providerName: "coinbase",
			expectedMap: types.ProviderMarketMap{
				Name:          "coinbase",
				TickerConfigs: make(map[mmtypes.Ticker]mmtypes.ProviderConfig),
				OffChainMap:   map[string]mmtypes.Ticker{},
			},
			expErr: false,
		},
		{
			name: "valid market map with no entries for the given provider",
			marketMap: mmtypes.MarketMap{
				Tickers: map[string]mmtypes.Ticker{
					constants.BITCOIN_USD.String(): constants.BITCOIN_USD,
				},
				Providers: map[string]mmtypes.Providers{
					constants.BITCOIN_USD.String(): {
						Providers: []mmtypes.ProviderConfig{
							{
								Name:           "test",
								OffChainTicker: "BTC-USD",
							},
						},
					},
				},
			},
			providerName: "coinbase",
			expectedMap: types.ProviderMarketMap{
				Name:          "coinbase",
				TickerConfigs: make(map[mmtypes.Ticker]mmtypes.ProviderConfig),
				OffChainMap:   map[string]mmtypes.Ticker{},
			},
			expErr: false,
		},
		{
			name: "valid market map with entries for the given provider",
			marketMap: mmtypes.MarketMap{
				Tickers: map[string]mmtypes.Ticker{
					constants.BITCOIN_USD.String(): constants.BITCOIN_USD,
				},
				Providers: map[string]mmtypes.Providers{
					constants.BITCOIN_USD.String(): {
						Providers: []mmtypes.ProviderConfig{
							{
								Name:           "coinbase",
								OffChainTicker: "BTC-USD",
							},
						},
					},
				},
			},
			providerName: "coinbase",
			expectedMap: types.ProviderMarketMap{
				Name: "coinbase",
				TickerConfigs: types.TickerToProviderConfig{
					constants.BITCOIN_USD: {
						Name:           "coinbase",
						OffChainTicker: "BTC-USD",
					},
				},
				OffChainMap: map[string]mmtypes.Ticker{
					"BTC-USD": constants.BITCOIN_USD,
				},
			},
			expErr: false,
		},
		{
			name: "invalid market map",
			marketMap: mmtypes.MarketMap{
				Tickers: map[string]mmtypes.Ticker{
					constants.BITCOIN_USD.String(): constants.BITCOIN_USD,
				},
			},
			providerName: "coinbase",
			expectedMap:  types.ProviderMarketMap{},
			expErr:       true,
		},
		{
			name: "multiple providers for the same ticker",
			marketMap: mmtypes.MarketMap{
				Tickers: map[string]mmtypes.Ticker{
					constants.BITCOIN_USD.String(): constants.BITCOIN_USD,
				},
				Providers: map[string]mmtypes.Providers{
					constants.BITCOIN_USD.String(): {
						Providers: []mmtypes.ProviderConfig{
							{
								Name:           "coinbase",
								OffChainTicker: "BTC-USD",
							},
							{
								Name:           "test",
								OffChainTicker: "BTCs-USD",
							},
						},
					},
				},
			},
			providerName: "coinbase",
			expectedMap: types.ProviderMarketMap{
				Name: "coinbase",
				TickerConfigs: types.TickerToProviderConfig{
					constants.BITCOIN_USD: {
						Name:           "coinbase",
						OffChainTicker: "BTC-USD",
					},
				},
				OffChainMap: map[string]mmtypes.Ticker{
					"BTC-USD": constants.BITCOIN_USD,
				},
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pmap, err := types.ProviderMarketMapFromMarketMap(tc.providerName, tc.marketMap)
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedMap.Name, pmap.Name)
				require.Equal(t, tc.expectedMap.TickerConfigs, pmap.TickerConfigs)
				require.Equal(t, tc.expectedMap.OffChainMap, pmap.OffChainMap)
			}
		})
	}
}
