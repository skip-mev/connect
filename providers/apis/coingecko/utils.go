package coingecko

import (
	"fmt"
	"strings"
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
)

// NOTE: All documentation for this file can be located on the CoinGecko
// API documentation: https://www.coingecko.com/api/documentation. The CoinGecko
// API can be configured to be API based or not.

const (
	// Name is the name of the Coingecko provider.
	Name = "coingecko_api"

	// URL is the base URL for the CoinGecko API. This URL does not require
	// an API key but may be rate limited.
	URL = "https://api.coingecko.com/api/v3"

	// APIURL is the base URL for the CoinGecko API. This URL requires an API
	// key and is not rate limited.
	APIURL = "https://pro-api.coingecko.com/api/v3"

	// APIKeyHeader is the header used to pass the API key to the CoinGecko API.
	APIKeyHeader = "x-cg-pro-api-key" //nolint

	// PairPriceEndpoint is the URL used to fetch the price of a list of currency
	// pairs. The ids are the base currencies and the vs_currencies are the quote
	// currencies. Note that the IDs and vs_currencies are comma separated but are
	// not 1:1 in their representation.
	PairPriceEndpoint = "/simple/price?ids=%s&vs_currencies=%s"

	// Precision is the precision of the price returned by the CoinGecko API. All
	// results are returned with 18 decimal places and are expected to be converted
	// to the appropriate precision by the parser.
	Precision = "&precision=18"

	// TickerSeparator is the formatter of the ticker that is used to fetch the price
	// of a currency pair. The first currency is the base currency and the second
	// currency is the quote currency.
	TickerSeparator = "/"
)

// DefaultAPIConfig is the default configuration for the CoinGecko API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           true,
	Enabled:          true,
	Timeout:          3000 * time.Millisecond,
	Interval:         20 * time.Second, // Coingecko has a very low rate limit.
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	Endpoints:        []config.Endpoint{{URL: URL}},
}

type (
	// CoinGeckoResponse is the response returned by the CoinGecko API. The response
	// format looks like the following:
	// {
	// 		"bitcoin": {
	// 			"usd": 43808.30302432908,
	// 			"btc": 1
	// 		},
	// 		"ethereum": {
	// 			"usd": 2240.4139379890357,
	//			"btc": 0.05113686971792297
	// 		}
	// 	}
	CoinGeckoResponse map[string]map[string]float64 //nolint
)

// getUniqueBaseAndQuoteDenoms returns a list of unique base and quote denoms
// from a list of tickers. Note that this function will only return the denoms
// that are configured for the handler. If any of the tickers are not configured,
// they will not be fetched.
func (h *APIHandler) getUniqueBaseAndQuoteDenoms(tickers []types.ProviderTicker) (string, string, error) {
	if len(tickers) == 0 {
		return "", "", fmt.Errorf("no tickers specified")
	}

	// Create a map of unique base and quote denoms.
	seenBases := make(map[string]struct{})
	bases := make([]string, 0)

	seenQuotes := make(map[string]struct{})
	quotes := make([]string, 0)

	// Iterate through every currency pair and add the base and quote to the
	// unique bases and quotes list as long as they are supported.
	for _, ticker := range tickers {
		// Split the market ticker into the base and quote currencies.
		split := strings.Split(ticker.GetOffChainTicker(), TickerSeparator)
		if len(split) != 2 {
			return "", "", fmt.Errorf("ticker %s is not formatted correctly", ticker.String())
		}

		base := split[0]
		if _, ok := seenBases[base]; !ok {
			seenBases[base] = struct{}{}
			bases = append(bases, base)
		}

		quote := split[1]
		if _, ok := seenQuotes[quote]; !ok {
			seenQuotes[quote] = struct{}{}
			quotes = append(quotes, quote)
		}

		h.cache.Add(ticker)
	}

	// If there are no bases or quotes, then none of the tickers are supported.
	if len(bases) == 0 {
		return "", "", fmt.Errorf("none of the base currencies are supported")
	}

	if len(quotes) == 0 {
		return "", "", fmt.Errorf("none of the quote currencies are supported")
	}

	return strings.Join(bases, ","), strings.Join(quotes, ","), nil
}
