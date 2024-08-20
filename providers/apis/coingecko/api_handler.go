package coingecko

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
)

var _ types.PriceAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the PriceAPIDataHandler interface for CoinGecko.
type APIHandler struct {
	// apiCfg is the config for the CoinGecko API.
	api config.APIConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewAPIHandler returns a new CoinGecko PriceAPIDataHandler.
func NewAPIHandler(
	api config.APIConfig,
) (types.PriceAPIDataHandler, error) {
	if api.Name != Name {
		return nil, fmt.Errorf("expected api config name %s, got %s", Name, api.Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api config for %s is not enabled", Name)
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config for %s: %w", Name, err)
	}

	return &APIHandler{
		api:   api,
		cache: types.NewProviderTickers(),
	}, nil
}

// CreateURL returns the URL that is used to fetch data from the CoinGecko API for the
// given tickers. The CoinGecko API supports fetching spot prices for multiple tickers
// in a single request.
func (h *APIHandler) CreateURL(
	tickers []types.ProviderTicker,
) (string, error) {
	// Create a list of base currencies and quote currencies.
	bases, quotes, err := h.getUniqueBaseAndQuoteDenoms(tickers)
	if err != nil {
		return "", err
	}

	// This creates the endpoint that needs to be requested regardless of whether
	// an API key is set.
	pricesEndPoint := fmt.Sprintf(PairPriceEndpoint, bases, quotes)
	finalEndpoint := fmt.Sprintf("%s%s", pricesEndPoint, Precision)

	// Otherwise, we just return the base url with the endpoint.
	return fmt.Sprintf("%s%s", h.api.Endpoints[0].URL, finalEndpoint), nil
}

// ParseResponse parses the response from the CoinGecko API. The response is expected
// to match every base currency with every quote currency. As such, we need to filter
// out the responses that are not expected. Note that the response will only return
// a response for the inputted tickers.
func (h *APIHandler) ParseResponse(
	tickers []types.ProviderTicker,
	resp *http.Response,
) types.PriceResponse {
	// Parse the response.
	var result CoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToDecode),
		)
	}

	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	// Filter out the responses that are not expected.
	for base, quotes := range result {
		for quote, price := range quotes {
			// The ticker is represented as base/quote.
			offChainTicker := fmt.Sprintf("%s%s%s", base, TickerSeparator, quote)

			// If the ticker is not configured, we skip it.
			ticker, ok := h.cache.FromOffChainTicker(offChainTicker)
			if !ok {
				continue
			}

			// Resolve the price.
			resolved[ticker] = types.NewPriceResult(
				big.NewFloat(price),
				time.Now().UTC(),
			)
		}
	}

	// Add all expected tickers that did not return a response to the unresolved
	// map.
	for _, ticker := range tickers {
		if _, resolvedOk := resolved[ticker]; !resolvedOk {
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("no response"), providertypes.ErrorNoResponse),
			}
		}
	}

	return types.NewPriceResponse(resolved, unresolved)
}
