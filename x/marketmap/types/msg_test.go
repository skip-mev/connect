package types_test

import (
	"testing"

	"github.com/skip-mev/chaintestutil/sample"
	"github.com/stretchr/testify/require"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/marketmap/types"
)

func TestValidateBasicMsgUpsertMarket(t *testing.T) {
	validCurrencyPair := connecttypes.CurrencyPair{
		Base:  "BTC",
		Quote: "ETH",
	}

	validTicker := types.Ticker{
		CurrencyPair:     validCurrencyPair,
		Decimals:         8,
		MinProviderCount: 2,
	}

	tcs := []struct {
		name       string
		msg        types.MsgUpsertMarkets
		expectPass bool
	}{
		{
			"if the authority is not an acc-address - fail",
			types.MsgUpsertMarkets{
				Authority: "invalid",
			},
			false,
		},
		{
			"if there are no creates -  fail",
			types.MsgUpsertMarkets{
				Authority: sample.Address(sample.Rand()),
			},
			false,
		},
		{
			"invalid ticker (0 decimals) - fail",
			types.MsgUpsertMarkets{
				Authority: sample.Address(sample.Rand()),
				Markets: []types.Market{
					{
						Ticker: types.Ticker{
							CurrencyPair:     validCurrencyPair,
							Decimals:         0,
							MinProviderCount: 0,
						},
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           "kucoin",
								OffChainTicker: "btc-eth",
							},
							{
								Name:           "mexc",
								OffChainTicker: "btceth",
							},
						},
					},
				},
			},
			false,
		},
		{
			"duplicate tickers",
			types.MsgUpsertMarkets{
				Authority: sample.Address(sample.Rand()),
				Markets: []types.Market{
					{
						Ticker: types.Ticker{
							CurrencyPair:     validCurrencyPair,
							Decimals:         0,
							MinProviderCount: 0,
						},
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           "kucoin",
								OffChainTicker: "btc-eth",
							},
							{
								Name:           "mexc",
								OffChainTicker: "btceth",
							},
						},
					},
					{
						Ticker: types.Ticker{
							CurrencyPair:     validCurrencyPair,
							Decimals:         0,
							MinProviderCount: 0,
						},
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           "kucoin",
								OffChainTicker: "btc-eth",
							},
							{
								Name:           "mexc",
								OffChainTicker: "btceth",
							},
						},
					},
				},
			},
			false,
		},
		{
			"valid message",
			types.MsgUpsertMarkets{
				Authority: sample.Address(sample.Rand()),
				Markets: []types.Market{
					{
						Ticker: validTicker,
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           "kucoin",
								OffChainTicker: "btc-eth",
							},
							{
								Name:           "mexc",
								OffChainTicker: "btceth",
							},
						},
					},
				},
			},
			true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if !tc.expectPass {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestValidateBasicMsgCreateMarket(t *testing.T) {
	validCurrencyPair := connecttypes.CurrencyPair{
		Base:  "BTC",
		Quote: "ETH",
	}

	validTicker := types.Ticker{
		CurrencyPair:     validCurrencyPair,
		Decimals:         8,
		MinProviderCount: 2,
	}

	tcs := []struct {
		name       string
		msg        types.MsgCreateMarkets
		expectPass bool
	}{
		{
			"if the authority is not an acc-address - fail",
			types.MsgCreateMarkets{
				Authority: "invalid",
			},
			false,
		},
		{
			"if there are no creates -  fail",
			types.MsgCreateMarkets{
				Authority: sample.Address(sample.Rand()),
			},
			false,
		},
		{
			"invalid ticker (0 decimals) - fail",
			types.MsgCreateMarkets{
				Authority: sample.Address(sample.Rand()),
				CreateMarkets: []types.Market{
					{
						Ticker: types.Ticker{
							CurrencyPair:     validCurrencyPair,
							Decimals:         0,
							MinProviderCount: 0,
						},
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           "kucoin",
								OffChainTicker: "btc-eth",
							},
							{
								Name:           "mexc",
								OffChainTicker: "btceth",
							},
						},
					},
				},
			},
			false,
		},
		{
			"invalid num providers (need more than 1) - fail",
			types.MsgCreateMarkets{
				Authority: sample.Address(sample.Rand()),
				CreateMarkets: []types.Market{
					{
						Ticker: validTicker,
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           "kucoin",
								OffChainTicker: "btc-eth",
							},
						},
					},
				},
			},
			false,
		},
		{
			"invalid empty off-chain ticker - fail",
			types.MsgCreateMarkets{
				Authority: sample.Address(sample.Rand()),
				CreateMarkets: []types.Market{
					{
						Ticker: validTicker,
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           "kucoin",
								OffChainTicker: "btc-eth",
							},
							{
								Name:           "mexc",
								OffChainTicker: "",
							},
						},
					},
				},
			},
			false,
		},
		{
			"valid message",
			types.MsgCreateMarkets{
				Authority: sample.Address(sample.Rand()),
				CreateMarkets: []types.Market{
					{
						Ticker: validTicker,
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           "kucoin",
								OffChainTicker: "btc-eth",
							},
							{
								Name:           "mexc",
								OffChainTicker: "btceth",
							},
						},
					},
				},
			},
			true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if !tc.expectPass {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestValidateBasicMsgUpdateMarket(t *testing.T) {
	validCurrencyPair := connecttypes.CurrencyPair{
		Base:  "BTC",
		Quote: "ETH",
	}

	validTicker := types.Ticker{
		CurrencyPair:     validCurrencyPair,
		Decimals:         8,
		MinProviderCount: 2,
	}

	tcs := []struct {
		name       string
		msg        types.MsgUpdateMarkets
		expectPass bool
	}{
		{
			"if the Authority is not an acc-address - fail",
			types.MsgUpdateMarkets{
				Authority: "invalid",
			},
			false,
		},
		{
			"if there are no creates -  fail",
			types.MsgUpdateMarkets{
				Authority: sample.Address(sample.Rand()),
			},
			false,
		},
		{
			"invalid ticker (decimals) - fail",
			types.MsgUpdateMarkets{
				Authority: sample.Address(sample.Rand()),
				UpdateMarkets: []types.Market{
					{
						Ticker: types.Ticker{
							CurrencyPair:     validCurrencyPair,
							Decimals:         0,
							MinProviderCount: 0,
						},
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           "kucoin",
								OffChainTicker: "btc-eth",
							},
							{
								Name:           "mexc",
								OffChainTicker: "btceth",
							},
						},
					},
				},
			},
			false,
		},
		{
			"invalid num providers (need more than 1) - fail",
			types.MsgUpdateMarkets{
				Authority: sample.Address(sample.Rand()),
				UpdateMarkets: []types.Market{
					{
						Ticker: validTicker,
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           "kucoin",
								OffChainTicker: "btc-eth",
							},
						},
					},
				},
			},
			false,
		},
		{
			"invalid empty offchain ticker - fail",
			types.MsgUpdateMarkets{
				Authority: sample.Address(sample.Rand()),
				UpdateMarkets: []types.Market{
					{
						Ticker: validTicker,
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           "kucoin",
								OffChainTicker: "btc-eth",
							},
							{
								Name:           "mexc",
								OffChainTicker: "",
							},
						},
					},
				},
			},
			false,
		},
		{
			"valid message",
			types.MsgUpdateMarkets{
				Authority: sample.Address(sample.Rand()),
				UpdateMarkets: []types.Market{
					{
						Ticker: validTicker,
						ProviderConfigs: []types.ProviderConfig{
							{
								Name:           "kucoin",
								OffChainTicker: "btc-eth",
							},
							{
								Name:           "mexc",
								OffChainTicker: "btceth",
							},
						},
					},
				},
			},
			true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if !tc.expectPass {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestValidateBasicMsgParams(t *testing.T) {
	rng := sample.Rand()

	tcs := []struct {
		name       string
		msg        types.MsgParams
		expectPass bool
	}{
		{
			"if the Authority is not an acc-address - fail",
			types.MsgParams{
				Authority: "invalid",
			},
			false,
		},
		{
			name: "invalid params (no authorities) - fail",
			msg: types.MsgParams{
				Params: types.Params{
					MarketAuthorities: nil,
				},
				Authority: sample.Address(rng),
			},
			expectPass: false,
		},
		{
			name: "valid params",
			msg: types.MsgParams{
				Params: types.Params{
					MarketAuthorities: []string{sample.Address(rng)},
					Admin:             sample.Address(rng),
				},
				Authority: sample.Address(rng),
			},
			expectPass: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if !tc.expectPass {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestValidateBasicMsgRemoveMarketAuthorities(t *testing.T) {
	rng := sample.Rand()

	sampleAuth := sample.Address(rng)

	tcs := []struct {
		name       string
		msg        types.MsgRemoveMarketAuthorities
		expectPass bool
	}{
		{
			"if the Admin is not an acc-address - fail",
			types.MsgRemoveMarketAuthorities{
				Admin: "invalid",
			},
			false,
		},
		{
			name: "invalid message (no authorities) - fail",
			msg: types.MsgRemoveMarketAuthorities{
				RemoveAddresses: nil,
				Admin:           sample.Address(rng),
			},
			expectPass: false,
		},
		{
			name: "valid message",
			msg: types.MsgRemoveMarketAuthorities{
				RemoveAddresses: []string{sample.Address(rng)},
				Admin:           sample.Address(rng),
			},
			expectPass: true,
		},
		{
			name: "invalid message (duplicate authorities",
			msg: types.MsgRemoveMarketAuthorities{
				RemoveAddresses: []string{sampleAuth, sampleAuth},
				Admin:           sample.Address(rng),
			},
			expectPass: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if !tc.expectPass {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestValidateBasicMsgRemoveMarkets(t *testing.T) {
	rng := sample.Rand()

	tcs := []struct {
		name       string
		msg        types.MsgRemoveMarkets
		expectPass bool
	}{
		{
			"if the Authority is not an acc-address - fail",
			types.MsgRemoveMarkets{
				Admin: "invalid",
			},
			false,
		},
		{
			name: "invalid message (no markets) - fail",
			msg: types.MsgRemoveMarkets{
				Markets: nil,
				Admin:   sample.Address(rng),
			},
			expectPass: false,
		},
		{
			name: "valid message - single market",
			msg: types.MsgRemoveMarkets{
				Markets: []string{"USDT/USD"},
				Admin:   sample.Address(rng),
			},
			expectPass: true,
		},
		{
			name: "valid message - multiple markets",
			msg: types.MsgRemoveMarkets{
				Markets: []string{"USDT/USD", "ETH/USD"},
				Admin:   sample.Address(rng),
			},
			expectPass: true,
		},
		{
			name: "invalid message (duplicate markets",
			msg: types.MsgRemoveMarkets{
				Markets: []string{"USDT/USD", "USDT/USD"},
				Admin:   sample.Address(rng),
			},
			expectPass: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if !tc.expectPass {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}
