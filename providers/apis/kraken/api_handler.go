package kraken

import (
	"errors"
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

// APIHandler implements the PriceAPIDataHandler interface for Kraken.
// for more information about the Kraken API, refer to the following link:
// https://docs.kraken.com/rest/
type APIHandler struct {
	// api is the config for the Kraken API.
	api config.APIConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewAPIHandler returns a new Kraken PriceAPIDataHandler.
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

// CreateURL returns the URL that is used to fetch data from the Kraken API for the
// given tickers.
func (h *APIHandler) CreateURL(
	tickers []types.ProviderTicker,
) (string, error) {
	var tickerStrings string
	for _, ticker := range tickers {
		tickerStrings += fmt.Sprintf("%s%s", ticker.GetOffChainTicker(), Separator)
		h.cache.Add(ticker)
	}

	if len(tickerStrings) == 0 {
		return "", fmt.Errorf("empty url created. invalid or no ticker were provided")
	}

	return fmt.Sprintf(
		h.api.Endpoints[0].URL,
		strings.TrimSuffix(tickerStrings, Separator),
	), nil
}

// ParseResponse parses the response from the Kraken API and returns a GetResponse. Each
// of the tickers supplied will get a response or an error.
func (h *APIHandler) ParseResponse(
	tickers []types.ProviderTicker,
	resp *http.Response,
) types.PriceResponse {
	// Parse the response into a ResponseBody.
	result, err := Decode(resp)
	if err != nil {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToDecode),
		)
	}

	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	// The Kraken API will return an empty list of errors with an API result containing valid tickers. However, it's
	// easier for us to validate that there were no errors if this field is set to nil whenever it's empty.
	if len(result.Errors) == 0 {
		result.Errors = nil
	}

	if len(result.Errors) > 0 {
		err := fmt.Errorf(
			"kraken API call error: %w", errors.New(strings.Join(result.Errors, ", ")),
		)
		return types.NewPriceResponseWithErr(tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorInvalidResponse),
		)
	}

	for pair, resultTicker := range result.Tickers {
		resultTicker.pair = pair
		ticker, ok := h.cache.FromOffChainTicker(pair)
		if !ok {
			continue
		}

		price, err := math.Float64StringToBigFloat(resultTicker.LastPrice())
		if err != nil {
			wErr := fmt.Errorf(
				"failed to convert price %s to big.Float: %w",
				resultTicker.LastPrice(),
				err,
			)

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

	// Add currency pairs that received no response to the unresolved map.
	for _, ticker := range tickers {
		_, resolvedOk := resolved[ticker]
		_, unresolvedOk := unresolved[ticker]

		if !resolvedOk && !unresolvedOk {
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(
					fmt.Errorf("no response"),
					providertypes.ErrorNoResponse,
				),
			}
		}
	}

	return types.NewPriceResponse(resolved, unresolved)
}
