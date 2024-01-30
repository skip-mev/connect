package binance

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// NOTE: All the documentation for this file can be located on the Binance github
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

	// DefaultMarketConfig is the default market configuration for Binance.
	DefaultMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"BITCOIN/USD": {
				Ticker:       "BTCUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ETHUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"ATOM/USD": {
				Ticker:       "ATOMUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
			},
			"SOLANA/USD": {
				Ticker:       "SOLUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ETHBTC",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
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
	//]
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
