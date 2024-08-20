package geckoterminal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

var _ types.PriceAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the PriceAPIDataHandler interface for GeckoTerminal.
type APIHandler struct {
	// apiCfg is the config for the GeckoTerminal API.
	api config.APIConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewAPIHandler returns a new GeckoTerminal PriceAPIDataHandler.
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

// CreateURL returns the URL that is used to fetch data from the GeckoTerminal API for the
// given tickers. Note that the GeckoTerminal API supports fetching multiple spot prices
// iff they are all on the same chain.
func (h *APIHandler) CreateURL(
	tickers []types.ProviderTicker,
) (string, error) {
	addresses := make([]string, len(tickers))
	for i, ticker := range tickers {
		addresses[i] = ticker.GetOffChainTicker()
		h.cache.Add(ticker)
	}

	return fmt.Sprintf(h.api.Endpoints[0].URL, strings.Join(addresses, ",")), nil
}

// ParseResponse parses the response from the GeckoTerminal API. The response is expected
// to contain multiple spot prices for a given token address. Note that all of the tokens
// are shared on the same chain.
func (h *APIHandler) ParseResponse(
	tickers []types.ProviderTicker,
	resp *http.Response,
) types.PriceResponse {
	// Parse the response.
	var result GeckoTerminalResponse
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

	data := result.Data
	if data.Type != ExpectedResponseType {
		err := fmt.Errorf("expected type %s, got %s", ExpectedResponseType, data.Type)
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorInvalidResponse),
		)
	}

	// Filter out the responses that are not expected.
	attributes := data.Attributes
	for address, price := range attributes.TokenPrices {
		ticker, ok := h.cache.FromOffChainTicker(address)
		if !ok {
			err := fmt.Errorf("no ticker for address %s", address)
			return types.NewPriceResponseWithErr(
				tickers,
				providertypes.NewErrorWithCode(err, providertypes.ErrorUnknownPair),
			)
		}

		// Convert the price to a big.Float.
		price, err := math.Float64StringToBigFloat(price)
		if err != nil {
			wErr := fmt.Errorf("failed to convert price to big.Float: %w", err)
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(
					wErr,
					providertypes.ErrorFailedToParsePrice,
				),
			}

			continue
		}

		resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
	}

	// Add all expected tickers that did not return a response to the unresolved
	// map.
	for _, ticker := range tickers {
		_, resolvedOk := resolved[ticker]
		_, unresolvedOk := unresolved[ticker]

		if !resolvedOk && !unresolvedOk {
			err := fmt.Errorf("received no price response")
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(err, providertypes.ErrorNoResponse),
			}
		}
	}

	return types.NewPriceResponse(resolved, unresolved)
}
