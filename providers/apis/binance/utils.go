package binance

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// NOTE: All documentation for this file can be located on the Binance GitHub
// API documentation: https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#symbol-price-ticker. This
// API does not require a subscription to use (i.e. No API key is required).

const (
	// Name is the name of the Binance provider.
	Name = "binance"

	// URL is the base URL of the Binance API. This includes the base and quote
	// currency pairs that need to be inserted into the URL. This URL should be utilized
	// by Non-US users.
	URL = "https://api.binance.com/api/v3/ticker/price?symbols=%s%s%s"

	// US_URL is the base URL of the Binance US API. This includes the base and quote
	// currency pairs that need to be inserted into the URL. This URL should be utilized
	// by US users. Note that the US URL does not support all the currency pairs that
	// the Non-US URL supports.
	US_URL = "https://api.binance.us/api/v3/ticker/price?symbols=%s%s%s" //nolint

	Quotation    = "%22"
	Separator    = ","
	LeftBracket  = "%5B"
	RightBracket = "%5D"
)

var (
	// DefaultUSAPIConfig is the default configuration for the Binance API.
	DefaultUSAPIConfig = config.APIConfig{
		Name:       Name,
		Atomic:     true,
		Enabled:    true,
		Timeout:    500 * time.Millisecond,
		Interval:   1 * time.Second,
		MaxQueries: 1,
		URL:        US_URL,
	}

	// DefaultNonUSAPIConfig is the default configuration for the Binance API.
	DefaultNonUSAPIConfig = config.APIConfig{
		Name:       Name,
		Atomic:     true,
		Enabled:    true,
		Timeout:    500 * time.Millisecond,
		Interval:   1 * time.Second,
		MaxQueries: 1,
		URL:        URL,
	}

	// DefaultUSMarketConfig is the default US market configuration for Binance.
	DefaultUSMarketConfig = mmtypes.MarketConfig{
		Name: Name,
		TickerConfigs: map[string]mmtypes.TickerConfig{
			"ATOM/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "ATOMUSDT",
			},
			"AVAX/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "AVAXUSDT",
			},
			"BITCOIN/USDC": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDC"),
					Decimals:     8,
				},
				OffChainTicker: "BTCUSDC",
			},
			"BITCOIN/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "BTCUSDT",
			},
			"ETHEREUM/BITCOIN": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
					Decimals:     8,
				},
				OffChainTicker: "ETHBTC",
			},
			"ETHEREUM/USDC": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDC"),
					Decimals:     8,
				},
				OffChainTicker: "ETHUSDC",
			},
			"ETHEREUM/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "ETHUSDT",
			},
			"SOLANA/USDC": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDC"),
					Decimals:     8,
				},
				OffChainTicker: "SOLUSDC",
			},
			"SOLANA/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "SOLUSDT",
			},
			"USDC/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "USDCUSDT",
			},
		},
	}

	// DefaultNonUSMarketConfig is the default market configuration for Binance.
	DefaultNonUSMarketConfig = mmtypes.MarketConfig{
		Name: Name,
		TickerConfigs: map[string]mmtypes.TickerConfig{
			"ATOM/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "ATOMUSDT",
			},
			"AVAX/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "AVAXUSDT",
			},
			"BITCOIN/USDC": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDC"),
					Decimals:     8,
				},
				OffChainTicker: "BTCUSDC",
			},
			"BITCOIN/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "BTCUSDT",
			},
			"ETHEREUM/BITCOIN": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
					Decimals:     8,
				},
				OffChainTicker: "ETHBTC",
			},
			"ETHEREUM/USDC": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDC"),
					Decimals:     8,
				},
				OffChainTicker: "ETHUSDC",
			},
			"ETHEREUM/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "ETHUSDT",
			},
			"SOLANA/USDC": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDC"),
					Decimals:     8,
				},
				OffChainTicker: "SOLUSDC",
			},
			"SOLANA/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "SOLUSDT",
			},
			"USDC/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "USDCUSDT",
			},
		},
	}
)

type (
	// Response is the expected response returned by the Binance API.
	// The response is json formatted.
	// Response format:
	//
	//	[
	//  {
	//    "symbol": "LTCBTC",
	//    "price": "4.00000200"
	//  },
	//  {
	//    "symbol": "ETHBTC",
	//    "price": "0.07946600"
	//  }
	// ].
	Response []Data

	// Data BinanceData is the data returned by the Binance API.
	Data struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
)

// Decode decodes the given http response into a BinanceResponse.
func Decode(resp *http.Response) (Response, error) {
	// Parse the response into a BinanceResponse.
	var result Response
	err := json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
