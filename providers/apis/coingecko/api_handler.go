package coingecko

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ types.PriceAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the PriceAPIDataHandler interface for CoinGecko.
type APIHandler struct {
	// marketCfg is the config for the CoinGecko API.
	market mmtypes.MarketConfig
	// apiCfg is the config for the CoinGecko API.
	api config.APIConfig
}

// NewAPIHandler returns a new CoinGecko PriceAPIDataHandler.
func NewAPIHandler(
	market mmtypes.MarketConfig,
	api config.APIConfig,
) (types.PriceAPIDataHandler, error) {
	if err := market.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid market config for %s: %w", Name, err)
	}

	if market.Name != Name {
		return nil, fmt.Errorf("expected market config name %s, got %s", Name, market.Name)
	}

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
		market: market,
		api:    api,
	}, nil
}

// CreateURL returns the URL that is used to fetch data from the CoinGecko API for the
// given tickers. The CoinGecko API supports fetching spot prices for multiple tickers
// in a single request.
func (h *APIHandler) CreateURL(
	tickers []mmtypes.Ticker,
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
	return fmt.Sprintf("%s%s", h.api.URL, finalEndpoint), nil
}

// ParseResponse parses the response from the CoinGecko API. The response is expected
// to match every base currency with every quote currency. As such, we need to filter
// out the responses that are not expected. Note that the response will only return
// a response for the inputted tickers.
func (h *APIHandler) ParseResponse(
	tickers []mmtypes.Ticker,
	resp *http.Response,
) types.PriceResponse {
	// Parse the response.
	var result CoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return providertypes.NewGetResponseWithErr[mmtypes.Ticker, *big.Int](tickers, err)
	}

	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	// Filter out the responses that are not expected.
	inverted := h.market.Invert()
	for base, quotes := range result {
		for quote, price := range quotes {
			// The ticker is represented as base/quote.
			ticker := fmt.Sprintf("%s%s%s", base, TickerSeparator, quote)

			// If the ticker is not configured, we skip it.
			market, ok := inverted[ticker]
			if !ok {
				continue
			}

			// Resolve the price.
			price := math.Float64ToBigInt(price, market.Ticker.Decimals)
			resolved[market.Ticker] = providertypes.NewResult(price, time.Now())
		}
	}

	// Add all of the expected tickers that did not return a response to the unresolved
	// map.
	for _, ticker := range tickers {
		if _, resolvedOk := resolved[ticker]; !resolvedOk {
			unresolved[ticker] = fmt.Errorf("no response")
		}
	}

	return providertypes.NewGetResponse(resolved, unresolved)
}
