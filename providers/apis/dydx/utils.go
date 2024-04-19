package dydx

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/skip-mev/slinky/providers/apis/defi/uniswapv3"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// Name is the name of the MarketMap provider.
	Name = "dydx_api"

	// ChainID is the chain ID for the dYdX market map provider.
	ChainID = "dydx-node"

	// Endpoint is the endpoint for the dYdX market map API.
	Endpoint = "%s/dydxprotocol/prices/params/market?limit=10000"

	// Delimeter is the delimeter used to separate the base and quote assets in a pair.
	Delimeter = "-"

	// UniswapV3TickerFields is the number of fields to expect to parse from a UniswapV3 ticker.
	UniswapV3TickerFields = 3

	// UniswapV3TickerSeparator is the separator for fields contained within a ticker for a uniswapv3_api provider.
	UniswapV3TickerSeparator = "-"
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
	URL:              "localhost:1317",
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
