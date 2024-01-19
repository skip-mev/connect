package coingecko

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/skip-mev/slinky/x/oracle/types"
// )

// // NOTE: All of the documentation for this file can be located on the CoinGecko
// // API documentation: https://www.coingecko.com/api/documentation. The CoinGecko
// // API can be configured to be API based or not.

// const (
// 	// APIKeyHeader is the header used to send the API key to the CoinGecko API.
// 	APIKeyHeader = "&x-cg-pro-api-key=" //nolint:gosec

// 	// BaseURL is the base URL for the CoinGecko API. This URL does not require
// 	// an API key but may be rate limited.
// 	BaseURL = "https://api.coingecko.com/api/v3"

// 	// APIURL is the api base URL for the CoinGecko API. This API requires an
// 	// API key and is rate limited according to the CoinGecko plan.
// 	APIURL = "https://pro-api.coingecko.com/api/v3"

// 	// PairPriceEndpoint is the URL used to fetch the price of a list of currency
// 	// pairs. The ids are the base currencies and the vs_currencies are the quote
// 	// currencies. Note that the IDs and vs_currencies are comma separated but are
// 	// not 1:1 in their representation.
// 	PairPriceEndpoint = "/simple/price?ids=%s&vs_currencies=%s"

// 	// Precision is the precision of the price returned by the CoinGecko API. All
// 	// results are returned with 18 decimal places and are expected to be converted
// 	// to the appropriate precision by the parser.
// 	Precision = "&precision=18"
// )

// type (
// 	// CoinGeckoResponse is the response returned by the CoinGecko API. The response
// 	// format looks like the following:
// 	// {
// 	// 		"bitcoin": {
// 	// 			"usd": 43808.30302432908,
// 	// 			"btc": 1
// 	// 		},
// 	// 		"ethereum": {
// 	// 			"usd": 2240.4139379890357,
// 	//			"btc": 0.05113686971792297
// 	// 		}
// 	// 	}
// 	CoinGeckoResponse map[string]map[string]float64 //nolint
// )

// // getUniqueBaseAndQuoteDenoms returns a list of unique base and quote denoms
// // from a list of currency pairs. Note that this function will only return the
// // denoms that are configured for the handler. If any of the currency pairs are
// // not configured, they will not be fetched.
// func (h *CoinGeckoAPIHandler) getUniqueBaseAndQuoteDenoms(pairs []types.CurrencyPair) (string, string, error) {
// 	if len(pairs) == 0 {
// 		return "", "", fmt.Errorf("no currency pairs specified")
// 	}

// 	// Create a map of unique base and quote denoms.
// 	seenBases := make(map[string]struct{})
// 	bases := make([]string, 0)

// 	seenQuotes := make(map[string]struct{})
// 	quotes := make([]string, 0)

// 	// Iterate through every currency pair and add the base and quote to the
// 	// unique bases and quotes list as long as they are supported.
// 	for _, cp := range pairs {
// 		if _, ok := seenBases[cp.Base]; !ok {
// 			if b, ok := h.SupportedBases[cp.Base]; ok {
// 				bases = append(bases, b)
// 			}

// 			seenBases[cp.Base] = struct{}{}
// 		}

// 		if _, ok := seenQuotes[cp.Quote]; !ok {
// 			if q, ok := h.SupportedQuotes[cp.Quote]; ok {
// 				quotes = append(quotes, q)
// 			}

// 			seenQuotes[cp.Quote] = struct{}{}
// 		}
// 	}

// 	// If there are no bases or quotes, then none of the currency pairs are
// 	// supported.
// 	if len(bases) == 0 {
// 		return "", "", fmt.Errorf("none of the base currencies are supported")
// 	}

// 	if len(quotes) == 0 {
// 		return "", "", fmt.Errorf("none of the quote currencies are supported")
// 	}

// 	return strings.Join(bases, ","), strings.Join(quotes, ","), nil
// }
