package types_test

import (
	"testing"

	"github.com/skip-mev/chaintestutil/sample"
	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/mm2/types"
)

func TestValidateBasicMsgCreateMarket(t *testing.T) {
	validCurrencyPair := slinkytypes.CurrencyPair{
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
	validCurrencyPair := slinkytypes.CurrencyPair{
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
					Version:           0,
				},
				Authority: sample.Address(sample.Rand()),
			},
			expectPass: false,
		},
		{
			name: "valid params",
			msg: types.MsgParams{
				Params: types.Params{
					MarketAuthorities: []string{sample.Address(sample.Rand())},
					Version:           0,
				},
				Authority: sample.Address(sample.Rand()),
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
