package coinbase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
)

var _ types.PriceAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the PriceAPIDataHandler interface for Coinbase, which can be used
// by a base provider. The DataHandler fetches data from the spot price Coinbase API. It is
// atomic in that it must request data from the Coinbase API sequentially for each ticker.
type APIHandler struct {
	// api is the config for the Coinbase API.
	api config.APIConfig
}

// NewAPIHandler returns a new Coinbase PriceAPIDataHandler.
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

// CreateURL returns the URL that is used to fetch data from the Coinbase API for the
// given tickers. Since the Coinbase API only supports fetching spot prices for a single
// ticker at a time, this function will return an error if the ticker slice contains more
// than one ticker.
func (h *APIHandler) CreateURL(
	tickers []types.ProviderTicker,
) (string, error) {
	if len(tickers) != 1 {
		return "", fmt.Errorf("expected 1 ticker, got %d", len(tickers))
	}
	return fmt.Sprintf(h.api.Endpoints[0].URL, tickers[0].GetOffChainTicker()), nil
}

// ParseResponse parses the spot price HTTP response from the Coinbase API and returns
// the resulting price. Note that this can only parse a single ticker at a time.
func (h *APIHandler) ParseResponse(
	tickers []types.ProviderTicker,
	resp *http.Response,
) types.PriceResponse {
	if len(tickers) != 1 {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(
				fmt.Errorf("expected 1 ticker, got %d", len(tickers)),
				providertypes.ErrorInvalidResponse,
			),
		)
	}

	// Parse the response into a CoinBaseResponse.
	var result CoinBaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToDecode),
		)
	}

	// Convert the float64 price into a big.Float.
	ticker := tickers[0]
	price, err := math.Float64StringToBigFloat(result.Data.Amount)
	if err != nil {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToParsePrice),
		)
	}

	return types.NewPriceResponse(
		types.ResolvedPrices{
			ticker: types.NewPriceResult(price, time.Now().UTC()),
		},
		nil,
	)
}
