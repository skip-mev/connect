package raydium_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	oracleconfig "github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/oracle/types"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/apis/defi/raydium"
	"github.com/skip-mev/slinky/providers/apis/defi/raydium/mocks"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

const (
	USDCVaultAddress = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh6"
	BTCVaultAddress  = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh5"
	ETHVaultAddress  = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh4"
	USDTVaultAddress = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh3"
	MOGVaultAddress  = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh2"
	SOLVaultAddress  = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh1"
)

// Test Provider init.
func TestProviderInit(t *testing.T) {
	t.Run("config fails validate basic", func(t *testing.T) {
		cfg := oracleconfig.APIConfig{
			Enabled:    true,
			MaxQueries: 0,
		}

		_, err := raydium.NewAPIPriceFetcher(
			oracletypes.ProviderMarketMap{},
			cfg,
			zap.NewNop(),
		)

		require.True(t, strings.Contains(err.Error(), "config for raydium is invalid"))
	})

	t.Run("market config fails validate basic", func(t *testing.T) {
		// valid config
		cfg := oracleconfig.APIConfig{
			Enabled:          false,
			MaxQueries:       2,
			Interval:         1 * time.Second,
			Timeout:          2 * time.Second,
			ReconnectTimeout: 2 * time.Second,
		}
		market := oracletypes.ProviderMarketMap{
			Name: raydium.Name,
			OffChainMap: map[string]mmtypes.Ticker{
				"BTC/USDC": {
					CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USDC"),
				},
			},
		}

		_, err := raydium.NewAPIPriceFetcher(
			market,
			cfg,
			zap.NewNop(),
		)
		require.True(t, strings.Contains(err.Error(), "market config for raydium is invalid"))
	})

	t.Run("incorrect provider name", func(t *testing.T) {
		cfg := oracleconfig.APIConfig{
			Enabled:          true,
			MaxQueries:       2,
			Interval:         1 * time.Second,
			Timeout:          2 * time.Second,
			ReconnectTimeout: 2 * time.Second,
		}
		market := oracletypes.ProviderMarketMap{
			Name: raydium.Name + "a",
		}

		_, err := raydium.NewAPIPriceFetcher(
			market,
			cfg,
			zap.NewNop(),
		)
		require.Error(t, err, fmt.Sprintf("config.Name is not %s", raydium.Name))

		cfg = oracleconfig.APIConfig{
			Enabled:          true,
			MaxQueries:       2,
			Interval:         1 * time.Second,
			Timeout:          2 * time.Second,
			ReconnectTimeout: 2 * time.Second,
			Name:             raydium.Name + "a",
		}
		market = oracletypes.ProviderMarketMap{
			Name: raydium.Name,
		}

		_, err = raydium.NewAPIPriceFetcher(
			market,
			cfg,
			zap.NewNop(),
		)
		require.Error(t, err, fmt.Sprintf("market.Name is not %s", raydium.Name))
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
		market := oracletypes.ProviderMarketMap{
			Name: raydium.Name,
		}

		_, err := raydium.NewAPIPriceFetcher(
			market,
			cfg,
			zap.NewNop(),
		)
		require.Error(t, err, "config is not enabled")
	})

	t.Run("unmarshalling metadata json for tickers fails", func(t *testing.T) {
		cfg := oracleconfig.APIConfig{
			Enabled:          true,
			MaxQueries:       2,
			Interval:         1 * time.Second,
			Timeout:          2 * time.Second,
			ReconnectTimeout: 2 * time.Second,
			Name:             raydium.Name,
			URL:              "https://raydium.io",
		}
		market := oracletypes.ProviderMarketMap{
			Name: raydium.Name,
			TickerConfigs: oracletypes.TickerToProviderConfig{
				mmtypes.Ticker{
					CurrencyPair:     slinkytypes.NewCurrencyPair("BTC", "USDC"),
					Decimals:         8,
					MinProviderCount: 1,
					Metadata_JSON:    "{}",
				}: {
					OffChainTicker: "BTC/USDC",
					Name:           raydium.Name,
				},
			},
			OffChainMap: map[string]mmtypes.Ticker{
				"BTC/USDC": {
					CurrencyPair:     slinkytypes.NewCurrencyPair("BTC", "USDC"),
					Decimals:         8,
					MinProviderCount: 1,
					Metadata_JSON:    "{}",
				},
			},
		}

		_, err := raydium.NewAPIPriceFetcher(
			market,
			cfg,
			zap.NewNop(),
		)
		t.Log(err)
		require.True(t, strings.Contains(err.Error(), "metadata for ticker BTC/USDC is invalid"))
	})

	t.Run("correctly unmarshals metadata json for ticker", func(t *testing.T) {
		cfg := oracleconfig.APIConfig{
			Enabled:          true,
			MaxQueries:       2,
			Interval:         1 * time.Second,
			Timeout:          2 * time.Second,
			ReconnectTimeout: 2 * time.Second,
			Name:             raydium.Name,
			URL:              "https://raydium.io",
		}
		market := oracletypes.ProviderMarketMap{
			Name: raydium.Name,
			TickerConfigs: oracletypes.TickerToProviderConfig{
				mmtypes.Ticker{
					CurrencyPair:     slinkytypes.NewCurrencyPair("BTC", "USDC"),
					Decimals:         8,
					MinProviderCount: 1,
					Metadata_JSON: `{
						"base_token_vault": {
							"token_vault_address": "` + USDCVaultAddress + `",
							"token_vault_decimals": 6
						},
						"quote_token_vault": {
							"token_vault_address": "` + BTCVaultAddress + `",
							"token_vault_decimals": 8
						}
					}`,
				}: {
					OffChainTicker: "BTC/USDC",
					Name:           raydium.Name,
				},
			},
			OffChainMap: map[string]mmtypes.Ticker{
				"BTC/USDC": {
					CurrencyPair:     slinkytypes.NewCurrencyPair("BTC", "USDC"),
					Decimals:         8,
					MinProviderCount: 1,
					Metadata_JSON: `{
						"base_token_vault": {
							"token_vault_address": "` + USDCVaultAddress + `",
							"token_vault_decimals": 6
						},
						"quote_token_vault": {
							"token_vault_address": "` + BTCVaultAddress + `",
							"token_vault_decimals": 8
						}
					}`,
				},
			},
		}

		_, err := raydium.NewAPIPriceFetcher(
			market,
			cfg,
			zap.NewNop(),
		)
		t.Log(err)
		require.NoError(t, err)
	})
}

// Test getting prices.
func TestProviderFetchPrices(t *testing.T) {
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

	tickers := []mmtypes.Ticker{
		{
			CurrencyPair:     slinkytypes.NewCurrencyPair("BTC", "USDC"),
			Decimals:         8,
			MinProviderCount: 1,
			Metadata_JSON:    marshalDataToJSON(btcUSDCMetadata),
		},
		{
			CurrencyPair:     slinkytypes.NewCurrencyPair("ETH", "USDT"),
			Decimals:         8,
			MinProviderCount: 1,
			Metadata_JSON:    marshalDataToJSON(ethUSDTMetadata),
		},
		{
			CurrencyPair:     slinkytypes.NewCurrencyPair("MOG", "SOL"),
			Decimals:         18,
			MinProviderCount: 1,
			Metadata_JSON:    marshalDataToJSON(mogSOLMetadata),
		},
	}

	client := mocks.NewSolanaJSONRPCClient(t)
	pf, err := newPriceFetcherFromTickers(tickers, client)
	require.NoError(t, err)

	t.Run("accounts resp returns len(tickers) * 2 accounts", func(t *testing.T) {
		ctx := context.Background()
		btcVaultPk := solana.MustPublicKeyFromBase58(BTCVaultAddress)
		usdcVaultPk := solana.MustPublicKeyFromBase58(USDCVaultAddress)
		ethVaultPk := solana.MustPublicKeyFromBase58(ETHVaultAddress)
		usdtVaultPk := solana.MustPublicKeyFromBase58(USDTVaultAddress)
		client.On("GetMultipleAccountsWithOpts", ctx, []solana.PublicKey{
			btcVaultPk, usdcVaultPk, ethVaultPk, usdtVaultPk,
		}, &rpc.GetMultipleAccountsOpts{
			Commitment: rpc.CommitmentFinalized,
		}).Return(
			&rpc.GetMultipleAccountsResult{}, nil,
		).Once()

		resp := pf.FetchPrices(ctx, tickers[:2])
		// expect a failed response
		require.Equal(t, len(resp.Resolved), 0)
		require.Equal(t, len(resp.UnResolved), 2)

		for _, result := range resp.UnResolved {
			require.True(t, strings.Contains(result.Error(), "expected 4 accounts, got 0"))
		}
	})

	t.Run("nil accounts are handled gracefully (skipped + added to unresolved)", func(t *testing.T) {
		ctx := context.Background()
		btcVaultPk := solana.MustPublicKeyFromBase58(BTCVaultAddress)
		usdcVaultPk := solana.MustPublicKeyFromBase58(USDCVaultAddress)
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

		client.On("GetMultipleAccountsWithOpts", ctx, []solana.PublicKey{
			btcVaultPk, usdcVaultPk, ethVaultPk, usdtVaultPk,
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
				},
			}, nil,
		)

		resp := pf.FetchPrices(ctx, tickers[:2])
		t.Log(resp)
		// expect a failed response
		require.Equal(t, len(resp.Resolved), 1)
		require.Equal(t, len(resp.UnResolved), 1)

		require.True(t, strings.Contains(resp.UnResolved[tickers[0]].Error(), "solana json-rpc error"))
		result := resp.Resolved[tickers[1]]
		require.Equal(t, result.Value.Uint64(), uint64(3e8))
	})
}

func marshalDataToJSON(obj interface{}) string {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func newPriceFetcherFromTickers(tickers []mmtypes.Ticker, client *mocks.SolanaJSONRPCClient) (*raydium.APIPriceFetcher, error) {
	cfg := oracleconfig.APIConfig{
		Enabled:          true,
		MaxQueries:       2,
		Interval:         1 * time.Second,
		Timeout:          2 * time.Second,
		ReconnectTimeout: 2 * time.Second,
		Name:             raydium.Name,
		URL:              "https://raydium.io",
	}
	market := oracletypes.ProviderMarketMap{
		Name:          raydium.Name,
		TickerConfigs: make(oracletypes.TickerToProviderConfig),
		OffChainMap:   make(map[string]mmtypes.Ticker),
	}

	for _, ticker := range tickers {
		market.TickerConfigs[ticker] = mmtypes.ProviderConfig{
			Name:           raydium.Name,
			OffChainTicker: ticker.String(),
		}
		market.OffChainMap[ticker.String()] = ticker
	}

	return raydium.NewAPIPriceFetcherWithClient(
		market,
		cfg,
		client,
		zap.NewExample(),
	)
}
