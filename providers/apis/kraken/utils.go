package kraken

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

// NOTE: All documentation for this file can be located on the Kraken docs.
// API documentation: https://docs.kraken.com/rest/. This
// API does not require a subscription to use (i.e. No API key is required).

const (
	// Name is the name of the Kraken API provider.
	Name = "kraken_api"

	// URL is the base URL of the Kraken API. This includes the base and quote
	// currency pairs that need to be inserted into the URL.
	URL = "https://api.kraken.com/0/public/Ticker?pair=%s"

	// Separator is the character that separates tickers in the query URL.
	Separator = ","
)

// DefaultAPIConfig is the default configuration for the Kraken API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           true,
	Enabled:          true,
	Timeout:          3000 * time.Millisecond,
	Interval:         600 * time.Millisecond,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	Endpoints:        []config.Endpoint{{URL: URL}},
}

// TickerResult is the result of a Kraken API call for a single ticker.
//
// https://api.kraken.com/0/public/Ticker
type TickerResult struct {
	pair            string
	ClosePriceStats []string `json:"c"`
}

func (ktr *TickerResult) LastPrice() string {
	return ktr.ClosePriceStats[0]
}

// ResponseBody returns a list of tickers for the response.  If there is an error, it will be included,
// and all Tickers will be undefined.
type ResponseBody struct {
	Errors  []string                `json:"error" validate:"omitempty"`
	Tickers map[string]TickerResult `json:"result"`
}

// Decode decodes the given http response into a TickerResult.
func Decode(resp *http.Response) (ResponseBody, error) {
	// Parse the response into a ResponseBody.
	var result ResponseBody
	err := json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
