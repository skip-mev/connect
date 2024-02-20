package types_test

import (
	"testing"

	"github.com/skip-mev/chaintestutil/sample"
	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/marketmap/types"
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
		msg        types.MsgCreateMarket
		expectPass bool
	}{
		{
			"if the signer is not an acc-address - fail",
			types.MsgCreateMarket{
				Signer: "invalid",
			},
			false,
		},
		{
			"invalid ticker - fail",
			types.MsgCreateMarket{
				Signer: sample.Address(sample.Rand()),
				Ticker: types.Ticker{
					CurrencyPair:     validCurrencyPair,
					Decimals:         0,
					MinProviderCount: 0,
				},
				Providers: types.Providers{
					Providers: []types.ProviderConfig{
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
				Paths: types.Paths{
					Paths: []types.Path{
						{
							Operations: []types.Operation{
								{
									CurrencyPair: validTicker.CurrencyPair,
									Invert:       false,
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"invalid num providers - fail",
			types.MsgCreateMarket{
				Signer: sample.Address(sample.Rand()),
				Ticker: validTicker,
				Providers: types.Providers{
					Providers: []types.ProviderConfig{
						{
							Name:           "kucoin",
							OffChainTicker: "btc-eth",
						},
					},
				},
				Paths: types.Paths{
					Paths: []types.Path{
						{
							Operations: []types.Operation{
								{
									CurrencyPair: validTicker.CurrencyPair,
									Invert:       false,
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"invalid empty offchain ticker - fail",
			types.MsgCreateMarket{
				Signer: sample.Address(sample.Rand()),
				Ticker: validTicker,
				Providers: types.Providers{
					Providers: []types.ProviderConfig{
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
				Paths: types.Paths{
					Paths: []types.Path{
						{
							Operations: []types.Operation{
								{
									CurrencyPair: validTicker.CurrencyPair,
									Invert:       false,
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"invalid no paths - fail",
			types.MsgCreateMarket{
				Signer: sample.Address(sample.Rand()),
				Ticker: validTicker,
				Providers: types.Providers{
					Providers: []types.ProviderConfig{
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
				Paths: types.Paths{},
			},
			false,
		},
		{
			"invalid path - fail",
			types.MsgCreateMarket{
				Signer: sample.Address(sample.Rand()),
				Ticker: validTicker,
				Providers: types.Providers{
					Providers: []types.ProviderConfig{
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
				Paths: types.Paths{
					Paths: make([]types.Path, 0),
				},
			},
			false,
		},
		{
			"valid message",
			types.MsgCreateMarket{
				Signer: sample.Address(sample.Rand()),
				Ticker: validTicker,
				Providers: types.Providers{
					Providers: []types.ProviderConfig{
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
				Paths: types.Paths{
					Paths: []types.Path{
						{
							Operations: []types.Operation{
								{
									CurrencyPair: validTicker.CurrencyPair,
									Invert:       false,
								},
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
