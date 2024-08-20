package dydx

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/apis/defi/raydium"
	"github.com/skip-mev/connect/v2/providers/apis/defi/uniswapv3"
)

const (
	// Name is the name of the MarketMap provider.
	Name = "dydx_api"

	// SwitchOverAPIHandlerName is the name of the dYdX switch over API.
	SwitchOverAPIHandlerName = "dydx_migration_api"

	// ResearchAPIHandlerName is the name of the dYdX research json API.
	ResearchAPIHandlerName = "dydx_research_json_api"

	// ResearchCMCAPIHandlerName is the name of the dYdX research json API that only returns CoinMarketCap markets.
	ResearchCMCAPIHandlerName = "dydx_research_coinmarketcap_api"

	// ChainID is the chain ID for the dYdX market map provider.
	ChainID = "dydx-node"

	// Endpoint is the endpoint for the dYdX market map API.
	Endpoint = "%s/dydxprotocol/prices/params/market?pagination.limit=10000"

	// Delimiter is the delimiter used to separate the base and quote assets in a pair.
	Delimiter = "-"

	// UniswapV3TickerFields is the number of fields to expect to parse from a UniswapV3 ticker.
	UniswapV3TickerFields = 3

	// UniswapV3TickerSeparator is the separator for fields contained within a ticker for a uniswapv3_api provider.
	UniswapV3TickerSeparator = Delimiter

	// RaydiumTickerFields is the minimum number of fields to expect the raydium exchange ticker to have.
	RaydiumTickerFields = 8

	// RaydiumTickerSeparator is the separator for fields contained within a ticker for the raydium provider.
	RaydiumTickerSeparator = Delimiter
)

// DefaultAPIConfig returns the default configuration for the dYdX market map API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           true,
	Enabled:          true,
	Timeout:          20 * time.Second, // Set a high timeout to account for slow API responses in the case where many markets are queried.
	Interval:         10 * time.Second,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	Endpoints:        []config.Endpoint{{URL: "http://localhost:1317"}},
}

// DefaultSwitchOverAPIConfig returns the default configuration for the dYdX switch over API provider.
var DefaultSwitchOverAPIConfig = config.APIConfig{
	Name:             SwitchOverAPIHandlerName,
	Atomic:           true,
	Enabled:          true,
	Timeout:          20 * time.Second, // Set a high timeout to account for slow API responses in the case where many markets are queried.
	Interval:         10 * time.Second,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	Endpoints: []config.Endpoint{
		{
			URL: "http://localhost:1317", // REST endpoint (HTTP/HTTPS prefix)
		},
		{
			URL: "localhost:9090", // gRPC endpoint (NO HTTP/HTTPS prefix)
		},
	},
}

// DefaultResearchAPIConfig returns the default configuration for the dYdX market map API.
var DefaultResearchAPIConfig = config.APIConfig{
	Name:             ResearchAPIHandlerName,
	Atomic:           true,
	Enabled:          true,
	Timeout:          20 * time.Second, // Set a high timeout to account for slow API responses in the case where many markets are queried.
	Interval:         10 * time.Second,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	Endpoints: []config.Endpoint{
		{
			URL: "https://dydx-api.lavenderfive.com",
		},
		{
			URL: "https://raw.githubusercontent.com/dydxprotocol/v4-web/main/public/configs/otherMarketData.json",
		},
	},
}

// DefaultResearchCMCAPIConfig returns the default configuration for the dYdX market map API that only returns CoinMarketCap markets.
var DefaultResearchCMCAPIConfig = config.APIConfig{
	Name:             ResearchCMCAPIHandlerName,
	Atomic:           true,
	Enabled:          true,
	Timeout:          20 * time.Second, // Set a high timeout to account for slow API responses in the case where many markets are queried.
	Interval:         10 * time.Second,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	Endpoints: []config.Endpoint{
		{
			URL: "https://dydx-api.lavenderfive.com",
		},
		{
			URL: "https://raw.githubusercontent.com/dydxprotocol/v4-web/main/public/configs/otherMarketData.json",
		},
	},
}

// UniswapV3MetadataFromTicker returns the metadataJSON string for uniswapv3_api according to the dYdX encoding.
// This is PoolAddress-DecimalsBase-DecimalsQuote.
func UniswapV3MetadataFromTicker(ticker string, invert bool) (string, error) {
	fields := strings.Split(ticker, UniswapV3TickerSeparator)
	if len(fields) != UniswapV3TickerFields {
		return "", fmt.Errorf("expected %d fields, got %d", UniswapV3TickerFields, len(fields))
	}

	baseDecimals, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse base decimals: %w", err)
	}

	quoteDecimals, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse quote decimals: %w", err)
	}

	parsedConfig := uniswapv3.PoolConfig{
		Address:       fields[0],
		BaseDecimals:  baseDecimals,
		QuoteDecimals: quoteDecimals,
		Invert:        invert,
	}

	if err = parsedConfig.ValidateBasic(); err != nil {
		return "", err
	}

	cfgBytes, err := json.Marshal(parsedConfig)
	if err != nil {
		return "", err
	}

	return string(cfgBytes), nil
}

// RaydiumMetadataFromTicker extracts json-metadata from a ticker for Raydium.
// All raydium tickers on dydx will be formatted as follows
// (BASE-QUOTE-BASE_VAULT-BASE_DECIMALS-QUOTE_VAULT-QUOTE_DECIMALS-OPEN_ORDERS_ADDRESS-AMM_INFO_ADDRESS).
func RaydiumMetadataFromTicker(ticker string) (string, error) {
	// split fields by separator and expect there to be at least 6 values
	fields := strings.Split(ticker, RaydiumTickerSeparator)
	if len(fields) < RaydiumTickerFields {
		return "", fmt.Errorf("expected at least 6 fields, got %d for ticker: %s", len(fields), ticker)
	}

	// check that vault addresses are valid solana addresses
	baseTokenVault := fields[2]
	if _, err := solana.PublicKeyFromBase58(baseTokenVault); err != nil {
		return "", fmt.Errorf("failed to parse base token vault: %w", err)
	}

	quoteTokenVault := fields[4]
	if _, err := solana.PublicKeyFromBase58(quoteTokenVault); err != nil {
		return "", fmt.Errorf("failed to parse quote token vault: %w", err)
	}

	// check that decimals are valid
	baseDecimals, err := strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse base decimals: %w", err)
	}

	quoteDecimals, err := strconv.ParseUint(fields[5], 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse quote decimals: %w", err)
	}

	// expect the open-orders address to be valid
	if _, err := solana.PublicKeyFromBase58(fields[6]); err != nil {
		return "", fmt.Errorf("failed to parse open orders address: %w", err)
	}

	// expect the amm id address to be valid
	if _, err := solana.PublicKeyFromBase58(fields[7]); err != nil {
		return "", fmt.Errorf("failed to parse amm id address: %w", err)
	}

	// create the Raydium metadata
	parsedConfig := raydium.TickerMetadata{
		BaseTokenVault: raydium.AMMTokenVaultMetadata{
			TokenVaultAddress: baseTokenVault,
			TokenDecimals:     baseDecimals,
		},
		QuoteTokenVault: raydium.AMMTokenVaultMetadata{
			TokenVaultAddress: quoteTokenVault,
			TokenDecimals:     quoteDecimals,
		},
		OpenOrdersAddress: fields[6],
		AMMInfoAddress:    fields[7],
	}
	// convert the metadata to json
	cfgBytes, err := json.Marshal(parsedConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal %s provider metadata for ticker %s: %w", raydium.Name, ticker, err)
	}

	return string(cfgBytes), nil
}
