package coingecko

import (
	"fmt"
	"strings"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// NOTE: All documentation for this file can be located on the CoinGecko
// API documentation: https://www.coingecko.com/api/documentation. The CoinGecko
// API can be configured to be API based or not.

const (
	// Name is the name of the Coingecko provider.
	Name = "coingecko"

	// URL is the base URL for the CoinGecko API. This URL does not require
	// an API key but may be rate limited.
	URL = "https://api.coingecko.com/api/v3"

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

var (
	// DefaultAPIConfig is the default configuration for the CoinGecko API.
	DefaultAPIConfig = config.APIConfig{
		Name:       Name,
		Atomic:     true,
		Enabled:    true,
		Timeout:    500 * time.Millisecond,
		Interval:   15 * time.Second, // Coingecko has a very low rate limit.
		MaxQueries: 1,
		URL:        URL,
	}

	// DefaultMarketConfig is the default market configuration for CoinGecko.
	DefaultMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"ATOM/USD": {
				Ticker:       "cosmos/usd",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
			},
			"BITCOIN/USD": {
				Ticker:       "bitcoin/usd",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"CELESTIA/USD": {
				Ticker:       "celestia/usd",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD"),
			},
			"DYDX/USD": {
				Ticker:       "dydx-chain/usd",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ethereum/btc",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ethereum/usd",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"OSMOSIS/USD": {
				Ticker:       "osmosis/usd",
				CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USD"),
			},
			"SOLANA/USD": {
				Ticker:       "solana/usd",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
			},
		},
	}
)

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
// from a list of currency pairs. Note that this function will only return the
// denoms that are configured for the handler. If any of the currency pairs are
// not configured, they will not be fetched.
func (h *APIHandler) getUniqueBaseAndQuoteDenoms(pairs []oracletypes.CurrencyPair) (string, string, error) {
	if len(pairs) == 0 {
		return "", "", fmt.Errorf("no currency pairs specified")
	}

	// Create a map of unique base and quote denoms.
	seenBases := make(map[string]struct{})
	bases := make([]string, 0)

	seenQuotes := make(map[string]struct{})
	quotes := make([]string, 0)

	// Iterate through every currency pair and add the base and quote to the
	// unique bases and quotes list as long as they are supported.
	for _, cp := range pairs {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.String()]
		if !ok {
			continue
		}

		// Split the market ticker into the base and quote currencies.
		split := strings.Split(market.Ticker, TickerSeparator)
		if len(split) != 2 {
			continue
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
	}

	// If there are no bases or quotes, then none of the currency pairs are
	// supported.
	if len(bases) == 0 {
		return "", "", fmt.Errorf("none of the base currencies are supported")
	}

	if len(quotes) == 0 {
		return "", "", fmt.Errorf("none of the quote currencies are supported")
	}

	return strings.Join(bases, ","), strings.Join(quotes, ","), nil
}
