package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/apis/coinbase"
	"github.com/skip-mev/connect/v2/x/marketmap/types"
)

var (
	btcusdtCP = slinkytypes.NewCurrencyPair("BTC", "USDT")

	btcusdt = types.Market{
		Ticker: types.Ticker{
			CurrencyPair:     btcusdtCP,
			Decimals:         8,
			MinProviderCount: 1,
			Enabled:          true,
		},
		ProviderConfigs: []types.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "btc-usdt",
			},
		},
	}

	btcusd = types.Market{
		Ticker: types.Ticker{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "BTC",
				Quote: "USD",
			},
			Decimals:         8,
			MinProviderCount: 1,
			Enabled:          true,
		},
		ProviderConfigs: []types.ProviderConfig{
			{
				Name:            "kucoin",
				OffChainTicker:  "btc-usdt",
				NormalizeByPair: &usdtusd.Ticker.CurrencyPair,
			},
		},
	}

	btcusdInvalid = types.Market{
		Ticker: types.Ticker{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "BTC",
				Quote: "USD",
			},
			Decimals:         8,
			MinProviderCount: 1,
			Enabled:          true,
		},
		ProviderConfigs: []types.ProviderConfig{
			{
				Name:            "kucoin",
				OffChainTicker:  "btc-usdt",
				NormalizeByPair: &usdtusdDisabled.Ticker.CurrencyPair,
			},
		},
	}

	usdtusd = types.Market{
		Ticker: types.Ticker{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "USDT",
				Quote: "USD",
			},
			Decimals:         8,
			MinProviderCount: 1,
			Enabled:          true,
		},
		ProviderConfigs: []types.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "usdt-usd",
			},
		},
	}

	usdtusdDisabled = types.Market{
		Ticker: types.Ticker{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "USDT",
				Quote: "USD",
			},
			Decimals:         8,
			MinProviderCount: 1,
			Enabled:          false,
		},
		ProviderConfigs: []types.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "usdt-usd",
			},
		},
	}

	usdcusd = types.Market{
		Ticker: types.Ticker{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "USDC",
				Quote: "USD",
			},
			Decimals:         8,
			MinProviderCount: 1,
			Enabled:          true,
		},
		ProviderConfigs: []types.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "usdc-usd",
			},
		},
	}

	ethusdt = types.Market{
		Ticker: types.Ticker{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "ETHEREUM",
				Quote: "USDT",
			},
			Decimals:         8,
			MinProviderCount: 1,
			Enabled:          true,
		},
		ProviderConfigs: []types.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "eth-usdt",
			},
		},
	}

	ethusd = types.Market{
		Ticker: types.Ticker{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "ETHEREUM",
				Quote: "USD",
			},
			Decimals:         8,
			MinProviderCount: 3,
			Enabled:          true,
		},
		ProviderConfigs: []types.ProviderConfig{
			{
				Name:            "kucoin",
				OffChainTicker:  "eth-usdt",
				NormalizeByPair: &usdtusd.Ticker.CurrencyPair,
			},
			{
				Name:            "binance",
				OffChainTicker:  "eth-usdt",
				NormalizeByPair: &usdtusd.Ticker.CurrencyPair,
			},
			{
				Name:            "mexc",
				OffChainTicker:  "eth-usdt",
				NormalizeByPair: &usdtusd.Ticker.CurrencyPair,
			},
		},
	}

	markets = map[string]types.Market{
		btcusdt.Ticker.String(): btcusdt,
		btcusd.Ticker.String():  btcusd,
		usdtusd.Ticker.String(): usdtusd,
		usdcusd.Ticker.String(): usdcusd,
		ethusdt.Ticker.String(): ethusdt,
		ethusd.Ticker.String():  ethusd,
	}

	invalidMarkets = map[string]types.Market{
		btcusdInvalid.Ticker.String():   btcusdInvalid,
		usdtusdDisabled.Ticker.String(): usdtusdDisabled,
	}
)

func TestMarketMapValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		marketMap types.MarketMap
		expectErr bool
	}{
		{
			name:      "valid empty",
			marketMap: types.MarketMap{},
			expectErr: false,
		},
		{
			name: "valid map",
			marketMap: types.MarketMap{
				Markets: markets,
			},
			expectErr: false,
		},
		{
			name: "invalid disabled normalizeByPair",
			marketMap: types.MarketMap{
				Markets: invalidMarkets,
			},
			expectErr: true,
		},
		{
			name: "market with no ticker",
			marketMap: types.MarketMap{
				Markets: map[string]types.Market{
					btcusdtCP.String(): {
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           coinbase.Name,
								OffChainTicker: "BTC-USD",
							},
						},
					},
				},
			},
			expectErr: true,
		},
		{
			name: "empty market",
			marketMap: types.MarketMap{
				Markets: map[string]types.Market{
					btcusdtCP.String(): {},
				},
			},
			expectErr: true,
		},
		{
			name: "provider config includes a ticker that is not supported",
			marketMap: types.MarketMap{
				Markets: map[string]types.Market{
					btcusdtCP.String(): {
						Ticker: types.Ticker{
							CurrencyPair:     btcusdtCP,
							Decimals:         8,
							MinProviderCount: 1,
						},
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:            coinbase.Name,
								OffChainTicker:  "btc-usd",
								NormalizeByPair: &slinkytypes.CurrencyPair{Base: "not", Quote: "real"},
								Invert:          false,
								Metadata_JSON:   "",
							},
						},
					},
				},
			},
			expectErr: true,
		},
		{
			name: "empty provider name",
			marketMap: types.MarketMap{
				Markets: map[string]types.Market{
					btcusdtCP.String(): {
						Ticker: types.Ticker{
							CurrencyPair:     btcusdtCP,
							Decimals:         8,
							MinProviderCount: 1,
						},
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           "",
								OffChainTicker: "btc-usd",
								Invert:         false,
								Metadata_JSON:  "",
							},
						},
					},
				},
			},
			expectErr: true,
		},
		{
			name: "no provider configs",
			marketMap: types.MarketMap{
				Markets: map[string]types.Market{
					btcusdtCP.String(): {
						Ticker: types.Ticker{
							CurrencyPair:     btcusdtCP,
							Decimals:         8,
							MinProviderCount: 1,
						},
						ProviderConfigs: []types.ProviderConfig{},
					},
				},
			},
			expectErr: true,
		},
		{
			name: "market-map with invalid key",
			marketMap: types.MarketMap{
				Markets: map[string]types.Market{
					ethusd.String(): {
						Ticker: types.Ticker{
							CurrencyPair:     btcusdtCP,
							Decimals:         8,
							MinProviderCount: 1,
						},
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           coinbase.Name,
								OffChainTicker: "BTC-USD",
							},
						},
					},
				},
			},
			expectErr: true,
		},
		{
			name: "valid single provider",
			marketMap: types.MarketMap{
				Markets: map[string]types.Market{
					btcusdtCP.String(): {
						Ticker: types.Ticker{
							CurrencyPair:     btcusdtCP,
							Decimals:         8,
							MinProviderCount: 1,
						},
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           coinbase.Name,
								OffChainTicker: "BTC-USD",
							},
						},
					},
				},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.marketMap.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMarketMapEqual(t *testing.T) {
	cases := []struct {
		name      string
		marketMap types.MarketMap
		other     types.MarketMap
		expect    bool
	}{
		{
			name:      "empty market map",
			marketMap: types.MarketMap{},
			other:     types.MarketMap{},
			expect:    true,
		},
		{
			name: "same market map",
			marketMap: types.MarketMap{
				Markets: map[string]types.Market{
					ethusdt.Ticker.String(): ethusdt,
				},
			},
			other: types.MarketMap{
				Markets: map[string]types.Market{
					ethusdt.Ticker.String(): ethusdt,
				},
			},
			expect: true,
		},
		{
			name: "different tickers",
			marketMap: types.MarketMap{
				Markets: map[string]types.Market{
					ethusdt.Ticker.String(): ethusdt,
				},
			},
			other: types.MarketMap{
				Markets: map[string]types.Market{
					btcusdt.Ticker.String(): btcusdt,
				},
			},
			expect: false,
		},
		{
			name: "different providers",
			marketMap: types.MarketMap{
				Markets: map[string]types.Market{
					ethusdt.Ticker.String(): ethusdt,
				},
			},
			other: types.MarketMap{
				Markets: map[string]types.Market{
					ethusdt.Ticker.String(): {
						Ticker:          ethusdt.Ticker,
						ProviderConfigs: btcusdt.ProviderConfigs,
					},
				},
			},
			expect: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expect, tc.marketMap.Equal(tc.other))
		})
	}
}
