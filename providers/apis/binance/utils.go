package binance

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
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

	// DefaultUSMarketConfig is the default market configuration for Binance.
	DefaultUSMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"BITCOIN/USD": {
				Ticker:       "BTCUSD",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"BITCOIN/USDT": {
				Ticker:       "BTCUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			"BITCOIN/USDC": {
				Ticker:       "BTCUSDC",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDC"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ETHUSD",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"ETHEREUM/USDT": {
				Ticker:       "ETHUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
			},
			"ETHEREUM/USDC": {
				Ticker:       "ETHUSDC",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDC"),
			},
			"ATOM/USD": {
				Ticker:       "ATOMUSD",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
			},
			"ATOM/USDT": {
				Ticker:       "ATOMUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USDT"),
			},
			"AVAX/USD": {
				Ticker:       "AVAXUSD",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
			},
			"AVAX/USDT": {
				Ticker:       "AVAXUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USDT"),
			},
			"SOLANA/USD": {
				Ticker:       "SOLUSD",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
			},
			"SOLANA/USDT": {
				Ticker:       "SOLUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDT"),
			},
			"SOLANA/USDC": {
				Ticker:       "SOLUSDC",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDC"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ETHBTC",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			"USDC/USD": {
				Ticker:       "USDCUSD",
				CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USD"),
			},
			"USDT/USD": {
				Ticker:       "USDTUSD",
				CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
			},
			"USDC/USDT": {
				Ticker:       "USDCUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USDT"),
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
