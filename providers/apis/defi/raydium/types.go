package raydium

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// Name is the name of the Raydium API.
	Name = "raydium_api"

	// NormalizedTokenAmountExponent.
	NormalizedTokenAmountExponent = 18
)

// updateMetaDataCache unmarshals the metadata JSON for each ticker and adds it to the
// metadata map.
func (pf *APIPriceFetcher) updateMetaDataCache(ticker types.ProviderTicker) (TickerMetadata, error) {
	if metadata, ok := pf.metaDataPerTicker[ticker.String()]; ok {
		return metadata, nil
	}

	metadata, err := unmarshalMetadataJSON(ticker.GetJSON())
	if err != nil {
		return TickerMetadata{}, fmt.Errorf("error unmarshalling metadata for ticker %s: %w", ticker.String(), err)
	}

	if err := metadata.ValidateBasic(); err != nil {
		return TickerMetadata{}, fmt.Errorf("metadata for ticker %s is invalid: %w", ticker.String(), err)
	}
	pf.metaDataPerTicker[ticker.String()] = metadata

	return metadata, nil
}

// TickerMetadata represents the metadata associated with a ticker's corresponding
// raydium pool.
type TickerMetadata struct {
	// BaseTokenVault is the metadata associated with the base token's token vault
	BaseTokenVault AMMTokenVaultMetadata `json:"base_token_vault"`

	// QuoteTokenVault is the metadata associated with the quote token's token vault
	QuoteTokenVault AMMTokenVaultMetadata `json:"quote_token_vault"`
}

// ValidateBasic checks that the solana token vault addresses are valid.
func (metadata TickerMetadata) ValidateBasic() error {
	if _, err := solana.PublicKeyFromBase58(metadata.BaseTokenVault.TokenVaultAddress); err != nil {
		return err
	}

	if _, err := solana.PublicKeyFromBase58(metadata.QuoteTokenVault.TokenVaultAddress); err != nil {
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

var DefaultAPIConfig = config.APIConfig{
	Enabled:          true,
	Name:             Name,
	Timeout:          2 * time.Second,
	Interval:         500 * time.Millisecond,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       10,
	Atomic:           false,
	BatchSize:        50, // maximal # of accounts in getMultipleAccounts query is 100
}
