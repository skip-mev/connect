package raydium

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
)

const (
	// Name is the name of the Raydium API.
	Name = "raydium_api"

	// NormalizedTokenAmountExponent.
	NormalizedTokenAmountExponent = 18
)

// metadataCache is a synchronous data-structure that holds the metadata for each ticker.
type metadataCache struct {
	// metaDataPerTicker is a map from ticker to metadata.
	metaDataPerTicker map[string]TickerMetadata

	// RWMutex is used to synchronize access to the metadata cache.
	mtx sync.RWMutex
}

func newMetadataCache() *metadataCache {
	return &metadataCache{
		metaDataPerTicker: make(map[string]TickerMetadata),
	}
}

func (mc *metadataCache) updateMetaDataCache(ticker types.ProviderTicker) (TickerMetadata, error) {
	mc.mtx.Lock()
	defer mc.mtx.Unlock()
	if metadata, ok := mc.metaDataPerTicker[ticker.String()]; ok {
		return metadata, nil
	}

	metadata, err := unmarshalMetadataJSON(ticker.GetJSON())
	if err != nil {
		return TickerMetadata{}, fmt.Errorf("error unmarshalling metadata for ticker %s: %w", ticker.String(), err)
	}
	if err := metadata.ValidateBasic(); err != nil {
		return TickerMetadata{}, fmt.Errorf("metadata for ticker %s is invalid: %w", ticker.String(), err)
	}
	mc.metaDataPerTicker[ticker.String()] = metadata

	return metadata, nil
}

// getMetadataPerTicker returns the metadata for the given ticker.
func (mc *metadataCache) getMetadataPerTicker(ticker types.ProviderTicker) (TickerMetadata, bool) {
	mc.mtx.RLock()
	defer mc.mtx.RUnlock()

	metaData, ok := mc.metaDataPerTicker[ticker.String()]
	return metaData, ok
}

// TickerMetadata represents the metadata associated with a ticker's corresponding
// raydium pool.
type TickerMetadata struct {
	// BaseTokenVault is the metadata associated with the base token's token vault
	BaseTokenVault AMMTokenVaultMetadata `json:"base_token_vault"`

	// QuoteTokenVault is the metadata associated with the quote token's token vault
	QuoteTokenVault AMMTokenVaultMetadata `json:"quote_token_vault"`

	// AMMInfoAddress is the address of the AMMInfo account for this raydium pool
	AMMInfoAddress string `json:"amm_info_address"`

	// OpenOrdersAddress is the address of the open orders account for this raydium pool
	OpenOrdersAddress string `json:"open_orders_address"`
}

// ValidateBasic checks that the solana token vault addresses are valid.
func (metadata TickerMetadata) ValidateBasic() error {
	if _, err := solana.PublicKeyFromBase58(metadata.BaseTokenVault.TokenVaultAddress); err != nil {
		return err
	}

	if _, err := solana.PublicKeyFromBase58(metadata.QuoteTokenVault.TokenVaultAddress); err != nil {
		return err
	}

	if _, err := solana.PublicKeyFromBase58(metadata.AMMInfoAddress); err != nil {
		return err
	}

	if _, err := solana.PublicKeyFromBase58(metadata.OpenOrdersAddress); err != nil {
		return err
	}

	return nil
}

// AMMTokenVaultMetadata represents the metadata associated with a raydium AMM pool's
// token vault. Specifically, we require the token vault address and the token decimals
// for the token that the vault is associated with.
type AMMTokenVaultMetadata struct {
	// QuoteTokenAddress is the base58 encoded address of the serum token corresponding
	// to this market's quote address
	TokenVaultAddress string `json:"token_vault_address"`

	// TokenDecimals is the number of decimals used for the token, we use this for
	// normalizing the balance of tokens at the designated vault address
	TokenDecimals uint64 `json:"token_decimals"`
}

// unmarshalMetadataJSON unmarshals the given metadata string into a TickerMetadata,
// this method assumes that the metadata string is valid json, otherwise an error is returned.
func unmarshalMetadataJSON(metadata string) (TickerMetadata, error) {
	// unmarshal the metadata string into a TickerMetadata
	var tickerMetadata TickerMetadata
	if err := json.Unmarshal([]byte(metadata), &tickerMetadata); err != nil {
		return TickerMetadata{}, err
	}

	return tickerMetadata, nil
}

// NoRaydiumMetadataForTickerError is returned when there is no metadata associated with a given ticker.
func NoRaydiumMetadataForTickerError(ticker string) error {
	return fmt.Errorf("no raydium metadata for ticker: %s", ticker)
}

// SolanaJSONRPCError is returned when there is an error querying the solana JSON-RPC client.
func SolanaJSONRPCError(err error) error {
	return fmt.Errorf("solana json-rpc error: %s", err.Error())
}

// DefaultAPIConfig is the default configuration for the Raydium API price fetcher.
var DefaultAPIConfig = config.APIConfig{
	Enabled:          true,
	Name:             Name,
	Timeout:          2 * time.Second,
	Interval:         500 * time.Millisecond,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       10,
	Atomic:           false,
	BatchSize:        25, // maximal # of accounts in getMultipleAccounts query is 100
	Endpoints: []config.Endpoint{
		{
			URL: "https://api.mainnet-beta.solana.com",
		},
	},
	MaxBlockHeightAge: 30 * time.Second,
}
