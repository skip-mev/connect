package geckoterminal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var _ types.PriceAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the PriceAPIDataHandler interface for GeckoTerminal.
type APIHandler struct {
	// apiCfg is the config for the GeckoTerminal API.
	api config.APIConfig
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
		api: api,
	}, nil
}

// CreateURL returns the URL that is used to fetch data from the GeckoTerminal API for the
// given tickers. Note that the GeckoTerminal API supports fetching multiple spot prices
// iff they are all on the same chain.
func (h *APIHandler) CreateURL(
	tickers []types.ProviderTicker,
) (string, error) {
	if len(tickers) > MaxNumberOfTickers {
		return "", fmt.Errorf("expected at most %d tickers, got %d", MaxNumberOfTickers, len(tickers))
	}

	addresses := make([]string, len(tickers))
	for i, ticker := range tickers {
		addresses[i] = ticker.OffChainTicker()
	}

	return fmt.Sprintf(h.api.URL, strings.Join(addresses, ",")), nil
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
		return types.NewPriceResponseWithErr(tickers, providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToDecode))
	}

	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	data := result.Data
	if data.Type != ExpectedResponseType {
		err := fmt.Errorf("expected type %s, got %s", ExpectedResponseType, data.Type)
		return types.NewPriceResponseWithErr(tickers, providertypes.NewErrorWithCode(err, providertypes.ErrorInvalidResponse))
	}

	// Filter out the responses that are not expected.
	attributes := data.Attributes
	providerTickers := types.NewProviderTickers(tickers...)
	for address, price := range attributes.TokenPrices {
		ticker, ok := providerTickers.FromOffChain(address)
		err := fmt.Errorf("no ticker for address %s", address)
		if !ok {
			return types.NewPriceResponseWithErr(tickers, providertypes.NewErrorWithCode(err, providertypes.ErrorUnknownPair))
		}

		// Convert the price to a big.Int.
		price, err := math.Float64StringToBigInt(price, ticker.Decimals())
		if err != nil {
			wErr := fmt.Errorf("failed to convert price to big.Int: %w", err)
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(wErr, providertypes.ErrorFailedToParsePrice),
			}
			continue
		}

		resolved[ticker] = types.NewPriceResult(price, time.Now())
	}

	// Add all expected tickers that did not return a response to the unresolved
	// map.
	for _, ticker := range tickers {
		if _, resolvedOk := resolved[ticker]; !resolvedOk {
			err := fmt.Errorf("received no price response")
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(err, providertypes.ErrorNoResponse),
			}
		}
	}

	return types.NewPriceResponse(resolved, unresolved)
}
