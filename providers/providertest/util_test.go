package providertest_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/providertest"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

var (
	usdtusdTicker = mmtypes.Ticker{
		CurrencyPair: connecttypes.CurrencyPair{
			Base:  "USDT",
			Quote: "USD",
		},
		Decimals:         8,
		MinProviderCount: 1,
		Enabled:          true,
	}

	usdtusdTickerDisabled = mmtypes.Ticker{
		CurrencyPair: connecttypes.CurrencyPair{
			Base:  "USDT",
			Quote: "USD",
		},
		Decimals:         8,
		MinProviderCount: 1,
		Enabled:          false,
	}

	usdtusdTickerMinProvider = mmtypes.Ticker{
		CurrencyPair: connecttypes.CurrencyPair{
			Base:  "USDT",
			Quote: "USD",
		},
		Decimals:         8,
		MinProviderCount: 10,
		Enabled:          true,
	}

	usdtusdProviderCfgOkx = mmtypes.ProviderConfig{
		Name:           "okx_ws",
		OffChainTicker: "USDC-USDT",
		Invert:         true,
	}

	usdtusdProviderCfgBitstamp = mmtypes.ProviderConfig{
		Name:           "bitstamp_api",
		OffChainTicker: "USDT-USD",
		Invert:         false,
	}

	usdtusdSingleProvider = mmtypes.Market{
		Ticker: usdtusdTicker,
		ProviderConfigs: []mmtypes.ProviderConfig{
			usdtusdProviderCfgOkx,
		},
	}

	usdtusdMultiProvider = mmtypes.Market{
		Ticker: usdtusdTicker,
		ProviderConfigs: []mmtypes.ProviderConfig{
			usdtusdProviderCfgOkx,
			usdtusdProviderCfgBitstamp,
		},
	}

	btctusdTicker = mmtypes.Ticker{
		CurrencyPair: connecttypes.CurrencyPair{
			Base:  "BTC",
			Quote: "USD",
		},
		Decimals:         8,
		MinProviderCount: 1,
		Enabled:          true,
	}

	btcusdProviderCfgOkx = mmtypes.ProviderConfig{
		Name:           "okx_ws",
		OffChainTicker: "USDC-BTC",
		Invert:         true,
	}

	btcusdProviderCfgBitstamp = mmtypes.ProviderConfig{
		Name:           "bitstamp_api",
		OffChainTicker: "BTC-USD",
		Invert:         false,
	}

	btcusdSingleProvider = mmtypes.Market{
		Ticker: btctusdTicker,
		ProviderConfigs: []mmtypes.ProviderConfig{
			btcusdProviderCfgOkx,
		},
	}

	btcusdMultiProvider = mmtypes.Market{
		Ticker: btctusdTicker,
		ProviderConfigs: []mmtypes.ProviderConfig{
			btcusdProviderCfgOkx,
			btcusdProviderCfgBitstamp,
		},
	}
)

func TestFilterMarketMapToProviders(t *testing.T) {
	tests := []struct {
		name  string
		input mmtypes.MarketMap
		want  map[string]mmtypes.MarketMap
	}{
		{
			name:  "empty",
			input: mmtypes.MarketMap{},
			want:  make(map[string]mmtypes.MarketMap),
		},
		{
			name: "single market with one provider",
			input: mmtypes.MarketMap{
				Markets: map[string]mmtypes.Market{
					usdtusdTicker.String(): usdtusdSingleProvider,
				},
			},
			want: map[string]mmtypes.MarketMap{
				"okx_ws": {
					Markets: map[string]mmtypes.Market{
						usdtusdTicker.String(): usdtusdSingleProvider,
					},
				},
			},
		},
		{
			name: "enable disabled markets for testing",
			input: mmtypes.MarketMap{
				Markets: map[string]mmtypes.Market{
					usdtusdTickerDisabled.String(): usdtusdSingleProvider,
				},
			},
			want: map[string]mmtypes.MarketMap{
				"okx_ws": {
					Markets: map[string]mmtypes.Market{
						usdtusdTicker.String(): usdtusdSingleProvider,
					},
				},
			},
		},
		{
			name: "set min provider count to 1 for testing",
			input: mmtypes.MarketMap{
				Markets: map[string]mmtypes.Market{
					usdtusdTickerMinProvider.String(): usdtusdSingleProvider,
				},
			},
			want: map[string]mmtypes.MarketMap{
				"okx_ws": {
					Markets: map[string]mmtypes.Market{
						usdtusdTicker.String(): usdtusdSingleProvider,
					},
				},
			},
		},
		{
			name: "single market with multi provider",
			input: mmtypes.MarketMap{
				Markets: map[string]mmtypes.Market{
					usdtusdTicker.String(): usdtusdMultiProvider,
				},
			},
			want: map[string]mmtypes.MarketMap{
				"okx_ws": {
					Markets: map[string]mmtypes.Market{
						usdtusdTicker.String(): {
							Ticker: usdtusdTicker,
							ProviderConfigs: []mmtypes.ProviderConfig{
								usdtusdProviderCfgOkx,
							},
						},
					},
				},
				"bitstamp_api": {
					Markets: map[string]mmtypes.Market{
						usdtusdTicker.String(): {
							Ticker: usdtusdTicker,
							ProviderConfigs: []mmtypes.ProviderConfig{
								usdtusdProviderCfgBitstamp,
							},
						},
					},
				},
			},
		},
		{
			name: "multi market with single provider",
			input: mmtypes.MarketMap{
				Markets: map[string]mmtypes.Market{
					usdtusdTicker.String(): usdtusdSingleProvider,
					btctusdTicker.String(): btcusdSingleProvider,
				},
			},
			want: map[string]mmtypes.MarketMap{
				"okx_ws": {
					Markets: map[string]mmtypes.Market{
						usdtusdTicker.String(): {
							Ticker: usdtusdTicker,
							ProviderConfigs: []mmtypes.ProviderConfig{
								usdtusdProviderCfgOkx,
							},
						},
						btctusdTicker.String(): {
							Ticker: btctusdTicker,
							ProviderConfigs: []mmtypes.ProviderConfig{
								btcusdProviderCfgOkx,
							},
						},
					},
				},
			},
		},
		{
			name: "multi market with multi provider",
			input: mmtypes.MarketMap{
				Markets: map[string]mmtypes.Market{
					usdtusdTicker.String(): usdtusdMultiProvider,
					btctusdTicker.String(): btcusdMultiProvider,
				},
			},
			want: map[string]mmtypes.MarketMap{
				"okx_ws": {
					Markets: map[string]mmtypes.Market{
						usdtusdTicker.String(): {
							Ticker: usdtusdTicker,
							ProviderConfigs: []mmtypes.ProviderConfig{
								usdtusdProviderCfgOkx,
							},
						},
						btctusdTicker.String(): {
							Ticker: btctusdTicker,
							ProviderConfigs: []mmtypes.ProviderConfig{
								btcusdProviderCfgOkx,
							},
						},
					},
				},
				"bitstamp_api": {
					Markets: map[string]mmtypes.Market{
						usdtusdTicker.String(): {
							Ticker: usdtusdTicker,
							ProviderConfigs: []mmtypes.ProviderConfig{
								usdtusdProviderCfgBitstamp,
							},
						},
						btctusdTicker.String(): {
							Ticker: btctusdTicker,
							ProviderConfigs: []mmtypes.ProviderConfig{
								btcusdProviderCfgBitstamp,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := providertest.FilterMarketMapToProviders(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}
