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
		msg        types.MsgCreateMarkets
		expectPass bool
	}{
		{
			"if the signer is not an acc-address - fail",
			types.MsgCreateMarkets{
				Signer: "invalid",
			},
			false,
		},
		{
			"invalid ticker - fail",
			types.MsgCreateMarkets{
				Signer: sample.Address(sample.Rand()),
				CreateMarkets: []types.Market{
					{
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
				},
			},
			false,
		},
		{
			"invalid num providers - fail",
			types.MsgCreateMarkets{
				Signer: sample.Address(sample.Rand()),
				CreateMarkets: []types.Market{
					{
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
				},
			},
			false,
		},
		{
			"invalid empty offchain ticker - fail",
			types.MsgCreateMarkets{
				Signer: sample.Address(sample.Rand()),
				CreateMarkets: []types.Market{
					{
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
				},
			},
			false,
		},
		{
			"invalid no paths - fail",
			types.MsgCreateMarkets{
				Signer: sample.Address(sample.Rand()),
				CreateMarkets: []types.Market{
					{
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
				},
			},
			false,
		},
		{
			"invalid path - fail",
			types.MsgCreateMarkets{
				Signer: sample.Address(sample.Rand()),
				CreateMarkets: []types.Market{
					{
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
				},
			},
			false,
		},
		{
			"valid message",
			types.MsgCreateMarkets{
				Signer: sample.Address(sample.Rand()),
				CreateMarkets: []types.Market{
					{
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
			"if the signer is not an acc-address - fail",
			types.MsgUpdateMarkets{
				Signer: "invalid",
			},
			false,
		},
		{
			"invalid ticker - fail",
			types.MsgUpdateMarkets{
				Signer: sample.Address(sample.Rand()),
				UpdateMarkets: []types.Market{
					{
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
				},
			},
			false,
		},
		{
			"invalid num providers - fail",
			types.MsgUpdateMarkets{
				Signer: sample.Address(sample.Rand()),
				UpdateMarkets: []types.Market{
					{
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
				},
			},
			false,
		},
		{
			"invalid empty offchain ticker - fail",
			types.MsgUpdateMarkets{
				Signer: sample.Address(sample.Rand()),
				UpdateMarkets: []types.Market{
					{
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
				},
			},
			false,
		},
		{
			"invalid no paths - fail",
			types.MsgUpdateMarkets{
				Signer: sample.Address(sample.Rand()),
				UpdateMarkets: []types.Market{
					{
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
				},
			},
			false,
		},
		{
			"invalid path - fail",
			types.MsgUpdateMarkets{
				Signer: sample.Address(sample.Rand()),
				UpdateMarkets: []types.Market{
					{
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
				},
			},
			false,
		},
		{
			"valid message",
			types.MsgUpdateMarkets{
				Signer: sample.Address(sample.Rand()),
				UpdateMarkets: []types.Market{
					{
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
