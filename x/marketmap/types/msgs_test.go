package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/chaintestutil/sample"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

func TestValidateBasicMsgCreateMarket(t *testing.T) {
	validTicker := types.Ticker{
		Base:             "BTC",
		Quote:            "ETH",
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
					Base:             "",
					Quote:            "",
					Decimals:         0,
					MinProviderCount: 0,
				},
				ProvidersToOffChainTickers: map[string]string{
					"kucoin": "btc-eth",
					"mexc":   "btceth",
				},
				Paths: []types.Path{
					{Operations: []types.Operation{
						{
							Ticker: validTicker,
							Invert: false,
						},
					}},
				},
			},
			false,
		},
		{
			"invalid num providers - fail",
			types.MsgCreateMarket{
				Signer: sample.Address(sample.Rand()),
				Ticker: validTicker,
				ProvidersToOffChainTickers: map[string]string{
					"kucoin": "btc-eth",
				},
				Paths: []types.Path{
					{Operations: []types.Operation{
						{
							Ticker: validTicker,
							Invert: false,
						},
					}},
				},
			},
			false,
		},
		{
			"invalid empty offchain ticker - fail",
			types.MsgCreateMarket{
				Signer: sample.Address(sample.Rand()),
				Ticker: validTicker,
				ProvidersToOffChainTickers: map[string]string{
					"kucoin": "btc-eth",
					"mexc":   "",
				},
				Paths: []types.Path{
					{Operations: []types.Operation{
						{
							Ticker: validTicker,
							Invert: false,
						},
					}},
				},
			},
			false,
		},
		{
			"invalid empty offchain ticker - fail",
			types.MsgCreateMarket{
				Signer: sample.Address(sample.Rand()),
				Ticker: validTicker,
				ProvidersToOffChainTickers: map[string]string{
					"kucoin": "btc-eth",
					"mexc":   "",
				},
				Paths: []types.Path{
					{Operations: []types.Operation{
						{
							Ticker: validTicker,
							Invert: false,
						},
					}},
				},
			},
			false,
		},
		{
			"invalid no paths - fail",
			types.MsgCreateMarket{
				Signer: sample.Address(sample.Rand()),
				Ticker: validTicker,
				ProvidersToOffChainTickers: map[string]string{
					"kucoin": "btc-eth",
					"mexc":   "btceth",
				},
				Paths: []types.Path{},
			},
			false,
		},
		{
			"invalid path - fail",
			types.MsgCreateMarket{
				Signer: sample.Address(sample.Rand()),
				Ticker: validTicker,
				ProvidersToOffChainTickers: map[string]string{
					"kucoin": "btc-eth",
					"mexc":   "btceth",
				},
				Paths: []types.Path{
					{},
				},
			},
			false,
		},
		{
			"valid message",
			types.MsgCreateMarket{
				Signer: sample.Address(sample.Rand()),
				Ticker: validTicker,
				ProvidersToOffChainTickers: map[string]string{
					"kucoin": "btc-eth",
					"mexc":   "btceth",
				},
				Paths: []types.Path{
					{Operations: []types.Operation{
						{
							Ticker: validTicker,
							Invert: false,
						},
					}},
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
