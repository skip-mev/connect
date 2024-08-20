package raydium_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/gagliardetto/solana-go/programs/serum"

	oracleconfig "github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/apis/defi/raydium"
	"github.com/skip-mev/connect/v2/providers/apis/defi/raydium/mocks"
	"github.com/skip-mev/connect/v2/providers/apis/defi/raydium/schema"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
)

const (
	USDCVaultAddress         = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh6"
	BTCVaultAddress          = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh5"
	USDCBTCAMMIDAddress      = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh7"
	USDCBTCOpenOrdersAddress = "9BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh1"
	ETHVaultAddress          = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh4"
	USDTVaultAddress         = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh3"
	ETHUSDTAMMIDAddress      = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh8"
	ETHUSDTOpenOrdersAddress = "9BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh2"
	MOGVaultAddress          = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh2"
	SOLVaultAddress          = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh1"
	MOGSOLAMMIDAddress       = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh9"
	MOGSOLOpenOrdersAddress  = "9BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh3"
)

func TestTickerMetadataValidateBasic(t *testing.T) {
	tcs := []struct {
		name string
		raydium.TickerMetadata
		expFail bool
	}{
		{
			name: "invalid base token vault address",
			TickerMetadata: raydium.TickerMetadata{
				BaseTokenVault: raydium.AMMTokenVaultMetadata{
					TokenVaultAddress: "",
					TokenDecimals:     6,
				},
				QuoteTokenVault: raydium.AMMTokenVaultMetadata{
					TokenVaultAddress: USDCVaultAddress,
					TokenDecimals:     6,
				},
			},
			expFail: true,
		},
		{
			name: "invalid quote token vault address",
			TickerMetadata: raydium.TickerMetadata{
				BaseTokenVault: raydium.AMMTokenVaultMetadata{
					TokenVaultAddress: USDCVaultAddress,
					TokenDecimals:     6,
				},
				QuoteTokenVault: raydium.AMMTokenVaultMetadata{
					TokenVaultAddress: "",
					TokenDecimals:     6,
				},
			},
			expFail: true,
		},
		{
			name: "invalid amm info address",
			TickerMetadata: raydium.TickerMetadata{
				BaseTokenVault: raydium.AMMTokenVaultMetadata{
					TokenVaultAddress: USDCVaultAddress,
					TokenDecimals:     6,
				},
				QuoteTokenVault: raydium.AMMTokenVaultMetadata{
					TokenVaultAddress: USDCVaultAddress,
					TokenDecimals:     6,
				},
			},
			expFail: true,
		},
		{
			name: "invalid open orders address",
			TickerMetadata: raydium.TickerMetadata{
				BaseTokenVault: raydium.AMMTokenVaultMetadata{
					TokenVaultAddress: USDCVaultAddress,
					TokenDecimals:     6,
				},
				QuoteTokenVault: raydium.AMMTokenVaultMetadata{
					TokenVaultAddress: USDCVaultAddress,
					TokenDecimals:     6,
				},
				AMMInfoAddress: USDCBTCAMMIDAddress,
			},
			expFail: true,
		},
		{
			name: "valid",
			TickerMetadata: raydium.TickerMetadata{
				BaseTokenVault: raydium.AMMTokenVaultMetadata{
					TokenVaultAddress: USDCVaultAddress,
					TokenDecimals:     6,
				},
				QuoteTokenVault: raydium.AMMTokenVaultMetadata{
					TokenVaultAddress: USDCVaultAddress,
					TokenDecimals:     6,
				},
				AMMInfoAddress:    USDCBTCAMMIDAddress,
				OpenOrdersAddress: USDCBTCOpenOrdersAddress,
			},
			expFail: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.TickerMetadata.ValidateBasic()
			if tc.expFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test Provider init.
func TestProviderInit(t *testing.T) {
	t.Run("config fails validate basic", func(t *testing.T) {
		cfg := oracleconfig.APIConfig{
			Enabled:    true,
			MaxQueries: 0,
		}

		_, err := raydium.NewAPIPriceFetcher(
			zap.NewNop(),
			cfg,
			metrics.NewNopAPIMetrics(),
		)

		require.Error(t, err)
	})

	t.Run("config has invalid endpoints", func(t *testing.T) {
		cfg := oracleconfig.APIConfig{
			Enabled:    true,
			MaxQueries: 0,
			Endpoints: []oracleconfig.Endpoint{
				{
					URL: "", // invalid url
				},
				{
					URL: "https://raydium.io",
				},
			},
		}

		_, err := raydium.NewAPIPriceFetcher(
			zap.NewNop(),
			cfg,
			metrics.NewNopAPIMetrics(),
		)

		require.Error(t, err)
	})

	t.Run("incorrect provider name", func(t *testing.T) {
		cfg := oracleconfig.APIConfig{
			Enabled:          true,
			MaxQueries:       2,
			Interval:         1 * time.Second,
			Timeout:          2 * time.Second,
			ReconnectTimeout: 2 * time.Second,
			Endpoints:        []oracleconfig.Endpoint{{URL: "https://raydium.io"}},
			Name:             raydium.Name + "a",
		}

		_, err := raydium.NewAPIPriceFetcher(
			zap.NewNop(),
			cfg,
			metrics.NewNopAPIMetrics(),
		)
		require.Error(t, err)
	})

	t.Run("api not enabled", func(t *testing.T) {
		cfg := oracleconfig.APIConfig{
			Enabled:          false,
			MaxQueries:       2,
			Interval:         1 * time.Second,
			Timeout:          2 * time.Second,
			ReconnectTimeout: 2 * time.Second,
			Name:             raydium.Name,
		}

		_, err := raydium.NewAPIPriceFetcher(
			zap.NewNop(),
			cfg,
			metrics.NewNopAPIMetrics(),
		)
		require.Error(t, err, "config is not enabled")
	})
}

// Test getting prices.
func TestProviderFetch(t *testing.T) {
	btcUSDCMetadata := raydium.TickerMetadata{
		BaseTokenVault: raydium.AMMTokenVaultMetadata{
			TokenVaultAddress: BTCVaultAddress,
			TokenDecimals:     8,
		},
		QuoteTokenVault: raydium.AMMTokenVaultMetadata{
			TokenVaultAddress: USDCVaultAddress,
			TokenDecimals:     6,
		},
		AMMInfoAddress:    USDCBTCAMMIDAddress,
		OpenOrdersAddress: USDCBTCOpenOrdersAddress,
	}
	ethUSDTMetadata := raydium.TickerMetadata{
		BaseTokenVault: raydium.AMMTokenVaultMetadata{
			TokenVaultAddress: ETHVaultAddress,
			TokenDecimals:     18,
		},
		QuoteTokenVault: raydium.AMMTokenVaultMetadata{
			TokenVaultAddress: USDTVaultAddress,
			TokenDecimals:     6,
		},
		AMMInfoAddress:    ETHUSDTAMMIDAddress,
		OpenOrdersAddress: ETHUSDTOpenOrdersAddress,
	}
	mogSOLMetadata := raydium.TickerMetadata{
		BaseTokenVault: raydium.AMMTokenVaultMetadata{
			TokenVaultAddress: MOGVaultAddress,
			TokenDecimals:     18,
		},
		QuoteTokenVault: raydium.AMMTokenVaultMetadata{
			TokenVaultAddress: SOLVaultAddress,
			TokenDecimals:     9,
		},
		AMMInfoAddress:    MOGSOLAMMIDAddress,
		OpenOrdersAddress: MOGSOLOpenOrdersAddress,
	}

	tickers := []types.DefaultProviderTicker{
		{
			OffChainTicker: "BTC/USDC",
			JSON:           marshalDataToJSON(btcUSDCMetadata),
		},
		{
			OffChainTicker: "ETH/USDT",
			JSON:           marshalDataToJSON(ethUSDTMetadata),
		},
		{
			OffChainTicker: "MOG/SOL",
			JSON:           marshalDataToJSON(mogSOLMetadata),
		},
	}

	client := mocks.NewSolanaJSONRPCClient(t)
	pf, err := newPriceFetcher(client)
	require.NoError(t, err)

	t.Run("accounts resp returns len(tickers) * 2 accounts", func(t *testing.T) {
		ctx := context.Background()

		btcVaultPk := solana.MustPublicKeyFromBase58(BTCVaultAddress)
		usdcVaultPk := solana.MustPublicKeyFromBase58(USDCVaultAddress)
		usdcBtcAMMIDPk := solana.MustPublicKeyFromBase58(USDCBTCAMMIDAddress)
		usdcBtcOpenOrdersPk := solana.MustPublicKeyFromBase58(USDCBTCOpenOrdersAddress)

		ethVaultPk := solana.MustPublicKeyFromBase58(ETHVaultAddress)
		usdtVaultPk := solana.MustPublicKeyFromBase58(USDTVaultAddress)
		ethUsdtAMMIDPk := solana.MustPublicKeyFromBase58(ETHUSDTAMMIDAddress)
		ETHUSDTOpenOrdersPk := solana.MustPublicKeyFromBase58(ETHUSDTOpenOrdersAddress)

		client.On("GetMultipleAccountsWithOpts", mock.Anything, []solana.PublicKey{
			btcVaultPk, usdcVaultPk, usdcBtcAMMIDPk, usdcBtcOpenOrdersPk,
			ethVaultPk, usdtVaultPk, ethUsdtAMMIDPk, ETHUSDTOpenOrdersPk,
		}, &rpc.GetMultipleAccountsOpts{
			Commitment: rpc.CommitmentConfirmed,
		}).Return(
			&rpc.GetMultipleAccountsResult{}, nil,
		).Once()

		ts := defaultTickersToProviderTickers(tickers[:2])
		resp := pf.Fetch(ctx, ts)
		// expect a failed response
		require.Equal(t, len(resp.Resolved), 0)
		require.Equal(t, len(resp.UnResolved), 2)

		for _, result := range resp.UnResolved {
			require.True(t, strings.Contains(result.Error(), "expected 8 accounts, got 0"))
		}
	})

	t.Run("failing accounts query", func(t *testing.T) {
		ctx := context.Background()
		err := fmt.Errorf("error")

		btcVaultPk := solana.MustPublicKeyFromBase58(BTCVaultAddress)
		usdcVaultPk := solana.MustPublicKeyFromBase58(USDCVaultAddress)
		usdcBtcAMMIDPk := solana.MustPublicKeyFromBase58(USDCBTCAMMIDAddress)
		usdcBtcOpenOrdersPk := solana.MustPublicKeyFromBase58(USDCBTCOpenOrdersAddress)

		ethVaultPk := solana.MustPublicKeyFromBase58(ETHVaultAddress)
		usdtVaultPk := solana.MustPublicKeyFromBase58(USDTVaultAddress)
		ethUsdtAMMIDPk := solana.MustPublicKeyFromBase58(ETHUSDTAMMIDAddress)
		ETHUSDTOpenOrdersPk := solana.MustPublicKeyFromBase58(ETHUSDTOpenOrdersAddress)

		client.On("GetMultipleAccountsWithOpts", mock.Anything, []solana.PublicKey{
			btcVaultPk, usdcVaultPk, usdcBtcAMMIDPk, usdcBtcOpenOrdersPk,
			ethVaultPk, usdtVaultPk, ethUsdtAMMIDPk, ETHUSDTOpenOrdersPk,
		}, &rpc.GetMultipleAccountsOpts{
			Commitment: rpc.CommitmentConfirmed,
		}).Return(
			&rpc.GetMultipleAccountsResult{}, err,
		).Once()

		ts := defaultTickersToProviderTickers(tickers[:2])
		resp := pf.Fetch(ctx, ts)
		// expect a failed response
		require.Equal(t, len(resp.Resolved), 0)
		require.Equal(t, len(resp.UnResolved), 2)

		for _, result := range resp.UnResolved {
			require.True(t, strings.Contains(result.Error(), raydium.SolanaJSONRPCError(err).Error()))
		}
	})

	t.Run("unexpected ticker in query", func(t *testing.T) {
		ctx := context.Background()

		mogtia := types.DefaultProviderTicker{
			OffChainTicker: "MOG/TIA",
			JSON:           "{}",
		}
		resp := pf.Fetch(ctx, []types.ProviderTicker{
			mogtia,
		})
		// expect a failed response
		require.Equal(t, len(resp.Resolved), 0)
		require.Equal(t, len(resp.UnResolved), 1)

		for _, result := range resp.UnResolved {
			t.Log(result.Error())
			require.True(t, strings.Contains(result.Error(), raydium.NoRaydiumMetadataForTickerError("MOG/TIA").Error()))
		}
	})

	t.Run("nil accounts are handled gracefully (skipped + added to unresolved)", func(t *testing.T) {
		ctx := context.Background()

		btcVaultPk := solana.MustPublicKeyFromBase58(BTCVaultAddress)
		usdcVaultPk := solana.MustPublicKeyFromBase58(USDCVaultAddress)
		usdcBtcAMMIDPk := solana.MustPublicKeyFromBase58(USDCBTCAMMIDAddress)
		usdcBtcOpenOrdersPk := solana.MustPublicKeyFromBase58(USDCBTCOpenOrdersAddress)

		mogVaultPk := solana.MustPublicKeyFromBase58(MOGVaultAddress)
		solVaultPk := solana.MustPublicKeyFromBase58(SOLVaultAddress)
		mogSolAMMIDPk := solana.MustPublicKeyFromBase58(MOGSOLAMMIDAddress)
		MOGSOLOpenOrdersPk := solana.MustPublicKeyFromBase58(MOGSOLOpenOrdersAddress)

		ethVaultPk := solana.MustPublicKeyFromBase58(ETHVaultAddress)
		usdtVaultPk := solana.MustPublicKeyFromBase58(USDTVaultAddress)
		ethUsdtAMMIDPk := solana.MustPublicKeyFromBase58(ETHUSDTAMMIDAddress)
		ETHUSDTOpenOrdersPk := solana.MustPublicKeyFromBase58(ETHUSDTOpenOrdersAddress)

		ethVaultBz := new(bytes.Buffer)
		ethEnc := bin.NewBinEncoder(ethVaultBz)
		ethVaultTokenMetadata := token.Account{
			Amount: uint64(1e18),
		}
		ethVaultTokenMetadata.MarshalWithEncoder(ethEnc)

		usdtVaultBz := new(bytes.Buffer)
		usdcEnc := bin.NewBinEncoder(usdtVaultBz)
		usdtTokenVaultMetadata := token.Account{
			Amount: 3 * (1e6),
		}
		usdtTokenVaultMetadata.MarshalWithEncoder(usdcEnc)

		ethUsdtAMMIDBz := new(bytes.Buffer)
		ethUsdtAMMIDEnc := bin.NewBinEncoder(ethUsdtAMMIDBz)
		ethUsdtAMMIDMetadata := schema.AmmInfo{
			OutPut: schema.OutPutData{
				NeedTakePnlCoin: uint64(6e17),
				NeedTakePnlPc:   uint64(16e5),
			},
		}
		ethUsdtAMMIDEnc.Encode(&ethUsdtAMMIDMetadata)

		ethUsdtOpenOrdersBz := new(bytes.Buffer)
		ethUsdtOpenOrdersEnc := bin.NewBinEncoder(ethUsdtOpenOrdersBz)
		ethUsdtOpenOrdersMetadata := serum.OpenOrders{
			NativeBaseTokenTotal:  bin.Uint64(1e17),
			NativeQuoteTokenTotal: bin.Uint64(0.1e6),
		}
		ethUsdtOpenOrdersEnc.Encode(&ethUsdtOpenOrdersMetadata)

		solVaultBz := new(bytes.Buffer)
		solEnc := bin.NewBinEncoder(solVaultBz)
		solTokenVaultMetadata := token.Account{
			Amount: 1e9,
		}
		solTokenVaultMetadata.MarshalWithEncoder(solEnc)

		client.On("GetMultipleAccountsWithOpts", mock.Anything, []solana.PublicKey{
			btcVaultPk, usdcVaultPk, usdcBtcAMMIDPk, usdcBtcOpenOrdersPk,
			ethVaultPk, usdtVaultPk, ethUsdtAMMIDPk, ETHUSDTOpenOrdersPk,
			mogVaultPk, solVaultPk, mogSolAMMIDPk, MOGSOLOpenOrdersPk,
		}, &rpc.GetMultipleAccountsOpts{
			Commitment: rpc.CommitmentConfirmed,
		}).Return(
			&rpc.GetMultipleAccountsResult{
				Value: []*rpc.Account{
					nil,
					{
						Data: nil,
					},
					{
						Data: nil,
					},
					nil,
					{
						Data: rpc.DataBytesOrJSONFromBytes(ethVaultBz.Bytes()),
					},
					{
						Data: rpc.DataBytesOrJSONFromBytes(usdtVaultBz.Bytes()),
					},
					{
						Data: rpc.DataBytesOrJSONFromBytes(ethUsdtAMMIDBz.Bytes()),
					},
					{
						Data: rpc.DataBytesOrJSONFromBytes(ethUsdtOpenOrdersBz.Bytes()),
					},
					{
						Data: rpc.DataBytesOrJSONFromBytes(solVaultBz.Bytes()),
					},
					nil,
					nil,
					nil,
				},
			}, nil,
		)

		ts := defaultTickersToProviderTickers(tickers[:3])
		resp := pf.Fetch(ctx, ts)

		// expect a failed response
		require.Equal(t, len(resp.Resolved), 1)
		require.Equal(t, len(resp.UnResolved), 2)

		require.True(t, strings.Contains(resp.UnResolved[tickers[0]].Error(), "solana json-rpc error"))
		result := resp.Resolved[tickers[1]]
		require.Equal(t, result.Value.SetPrec(30), big.NewFloat(3).SetPrec(30))
	})

	t.Run("incorrectly encoded accounts are handled gracefully", func(t *testing.T) {
		ctx := context.Background()

		btcVaultPk := solana.MustPublicKeyFromBase58(BTCVaultAddress)
		usdcVaultPk := solana.MustPublicKeyFromBase58(USDCVaultAddress)
		usdcBtcAMMIDPk := solana.MustPublicKeyFromBase58(USDCBTCAMMIDAddress)
		usdcBtcOpenOrdersPk := solana.MustPublicKeyFromBase58(USDCBTCOpenOrdersAddress)

		client.On("GetMultipleAccountsWithOpts", mock.Anything, []solana.PublicKey{
			btcVaultPk, usdcVaultPk, usdcBtcAMMIDPk, usdcBtcOpenOrdersPk,
		}, &rpc.GetMultipleAccountsOpts{
			Commitment: rpc.CommitmentConfirmed,
		}).Return(
			&rpc.GetMultipleAccountsResult{
				Value: []*rpc.Account{
					{
						Data: rpc.DataBytesOrJSONFromBytes([]byte{1, 2, 3}), // btc/usdc shld be unresolved
					},
					{
						Data: nil,
					},
				},
			}, nil,
		)

		ts := defaultTickersToProviderTickers(tickers[:1])
		resp := pf.Fetch(ctx, ts)

		// expect a failed response
		require.Equal(t, len(resp.Resolved), 0)
		require.Equal(t, len(resp.UnResolved), 1)

		require.True(t, strings.Contains(resp.UnResolved[tickers[0]].Error(), "solana json-rpc error"))
	})
}

func marshalDataToJSON(obj interface{}) string {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func newPriceFetcher(client *mocks.SolanaJSONRPCClient) (*raydium.APIPriceFetcher, error) {
	cfg := oracleconfig.APIConfig{
		Enabled:          true,
		MaxQueries:       2,
		Interval:         1 * time.Second,
		Timeout:          2 * time.Second,
		ReconnectTimeout: 2 * time.Second,
		Name:             raydium.Name,
		Endpoints:        []oracleconfig.Endpoint{{URL: "https://raydium.io"}},
	}

	return raydium.NewAPIPriceFetcherWithClient(
		zap.NewExample(),
		cfg,
		client,
	)
}

func defaultTickersToProviderTickers(tickers []types.DefaultProviderTicker) []types.ProviderTicker {
	providerTickers := make([]types.ProviderTicker, len(tickers))
	for i, ticker := range tickers {
		providerTickers[i] = ticker
	}
	return providerTickers
}
