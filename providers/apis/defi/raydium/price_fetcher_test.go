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

	oracleconfig "github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/defi/raydium"
	"github.com/skip-mev/slinky/providers/apis/defi/raydium/mocks"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
)

const (
	USDCVaultAddress = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh6"
	BTCVaultAddress  = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh5"
	ETHVaultAddress  = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh4"
	USDTVaultAddress = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh3"
	MOGVaultAddress  = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh2"
	SOLVaultAddress  = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh1"
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
		ethVaultPk := solana.MustPublicKeyFromBase58(ETHVaultAddress)
		usdtVaultPk := solana.MustPublicKeyFromBase58(USDTVaultAddress)
		client.On("GetMultipleAccountsWithOpts", mock.Anything, []solana.PublicKey{
			btcVaultPk, usdcVaultPk, ethVaultPk, usdtVaultPk,
		}, &rpc.GetMultipleAccountsOpts{
			Commitment: rpc.CommitmentFinalized,
		}).Return(
			&rpc.GetMultipleAccountsResult{}, nil,
		).Once()

		ts := defaultTickersToProviderTickers(tickers[:2])
		resp := pf.Fetch(ctx, ts)
		// expect a failed response
		require.Equal(t, len(resp.Resolved), 0)
		require.Equal(t, len(resp.UnResolved), 2)

		for _, result := range resp.UnResolved {
			require.True(t, strings.Contains(result.Error(), "expected 4 accounts, got 0"))
		}
	})

	t.Run("failing accounts query", func(t *testing.T) {
		ctx := context.Background()
		err := fmt.Errorf("error")
		btcVaultPk := solana.MustPublicKeyFromBase58(BTCVaultAddress)
		usdcVaultPk := solana.MustPublicKeyFromBase58(USDCVaultAddress)
		ethVaultPk := solana.MustPublicKeyFromBase58(ETHVaultAddress)
		usdtVaultPk := solana.MustPublicKeyFromBase58(USDTVaultAddress)
		client.On("GetMultipleAccountsWithOpts", mock.Anything, []solana.PublicKey{
			btcVaultPk, usdcVaultPk, ethVaultPk, usdtVaultPk,
		}, &rpc.GetMultipleAccountsOpts{
			Commitment: rpc.CommitmentFinalized,
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
		mogVaultPk := solana.MustPublicKeyFromBase58(MOGVaultAddress)
		solVaultPk := solana.MustPublicKeyFromBase58(SOLVaultAddress)
		ethVaultPk := solana.MustPublicKeyFromBase58(ETHVaultAddress)
		usdtVaultPk := solana.MustPublicKeyFromBase58(USDTVaultAddress)

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

		solVaultBz := new(bytes.Buffer)
		solEnc := bin.NewBinEncoder(solVaultBz)
		solTokenVaultMetadata := token.Account{
			Amount: 1e9,
		}
		solTokenVaultMetadata.MarshalWithEncoder(solEnc)

		client.On("GetMultipleAccountsWithOpts", mock.Anything, []solana.PublicKey{
			btcVaultPk, usdcVaultPk, ethVaultPk, usdtVaultPk, mogVaultPk, solVaultPk,
		}, &rpc.GetMultipleAccountsOpts{
			Commitment: rpc.CommitmentFinalized,
		}).Return(
			&rpc.GetMultipleAccountsResult{
				Value: []*rpc.Account{
					{
						Data: nil, // btc/usdc shld be unresolved
					},
					{
						Data: nil,
					},
					{
						Data: rpc.DataBytesOrJSONFromBytes(ethVaultBz.Bytes()),
					},
					{
						Data: rpc.DataBytesOrJSONFromBytes(usdtVaultBz.Bytes()),
					},
					{
						Data: rpc.DataBytesOrJSONFromBytes(solVaultBz.Bytes()),
					},
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

		client.On("GetMultipleAccountsWithOpts", mock.Anything, []solana.PublicKey{
			btcVaultPk, usdcVaultPk,
		}, &rpc.GetMultipleAccountsOpts{
			Commitment: rpc.CommitmentFinalized,
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
