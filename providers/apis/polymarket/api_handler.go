package polymarket

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

const (
	// Name is the name of the Polymarket provider.
	Name = "polymarket_api"

	// URL is the default base URL of the Polymarket CLOB API. It uses the midpoint endpoint with a given token ID.
	URL = "https://clob.polymarket.com/midpoint?token_id=%s"

	// priceAdjustmentMax is the value the price gets set to in the event of price == 1.00.
	priceAdjustmentMax = .9999
	priceAdjustmentMin = .00001
)

var (
	_ types.PriceAPIDataHandler = (*APIHandler)(nil)

	// valueExtractorFromEndpoint maps a URL path to a function that can extract the returned data from the response of that endpoint.
	valueExtractorFromEndpoint = map[string]valueExtractor{
		"/midpoint": dataFromMidpoint,
		"/price":    dataFromPrice,
	}
)

// APIHandler implements the PriceAPIDataHandler interface for Polymarket, which can be used
// by a base provider. The handler fetches data from either the `/midpoint` or `/price` endpoint.
type APIHandler struct {
	api config.APIConfig
}

// NewAPIHandler returns a new Polymarket PriceAPIDataHandler.
func NewAPIHandler(api config.APIConfig) (types.PriceAPIDataHandler, error) {
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

// CreateURL returns the URL that is used to fetch data from the Polymarket API for the
// given ticker. Since the midpoint endpoint is automatically denominated in USD, only one ID is expected to be passed
// into this method.
func (h APIHandler) CreateURL(ids []types.ProviderTicker) (string, error) {
	if len(ids) != 1 {
		return "", fmt.Errorf("expected 1 ticker, got %d", len(ids))
	}
	return fmt.Sprintf(h.api.Endpoints[0].URL, ids[0].GetOffChainTicker()), nil
}

// midpointResponseBody is the response structure for the `/midpoint` endpoint of the Polymarket API.
type midpointResponseBody struct {
	Mid string `json:"mid"`
}

// priceResponseBody is the response structure for the `/price` endpoint of the Polymarket API.
type priceResponseBody struct {
	Price string `json:"price"`
}

// valueExtractor is a function that can extract (price, midpoint) from a http response body.
// This function is expected to return a sting representation of a float.
type valueExtractor func(io.ReadCloser) (string, error)

// dataFromPrice unmarshalls data from the /price endpoint.
func dataFromPrice(reader io.ReadCloser) (string, error) {
	var result priceResponseBody
	err := json.NewDecoder(reader).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.Price, nil
}

// dataFromMidpoint unmarshalls data from the /midpoint endpoint.
func dataFromMidpoint(reader io.ReadCloser) (string, error) {
	var result midpointResponseBody
	err := json.NewDecoder(reader).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.Mid, nil
}

// ParseResponse parses the HTTP response from either the `/price` or `/midpoint` endpoint of the Polymarket API endpoint and returns
// the resulting data.
func (h APIHandler) ParseResponse(ids []types.ProviderTicker, response *http.Response) types.PriceResponse {
	if len(ids) != 1 {
		return types.NewPriceResponseWithErr(
			ids,
			providertypes.NewErrorWithCode(
				fmt.Errorf("expected 1 ticker, got %d", len(ids)),
				providertypes.ErrorInvalidResponse,
			),
		)
	}

	// get the extractor function for this endpoint.
	extractor, ok := valueExtractorFromEndpoint[response.Request.URL.Path]
	if !ok {
		return types.NewPriceResponseWithErr(
			ids,
			providertypes.NewErrorWithCode(fmt.Errorf("unknown request path %q", response.Request.URL.Path), providertypes.ErrorFailedToDecode),
		)
	}

	// extract the value. it should be a string representation of a float.
	val, err := extractor(response.Body)
	if err != nil {
		return types.NewPriceResponseWithErr(
			ids,
			providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToDecode),
		)
	}

	price, ok := new(big.Float).SetString(val)
	if !ok {
		return types.NewPriceResponseWithErr(
			ids,
			providertypes.NewErrorWithCode(fmt.Errorf("failed to convert %q to float", val), providertypes.ErrorFailedToDecode),
		)
	}
	if err := validatePrice(price); err != nil {
		return types.NewPriceResponseWithErr(
			ids,
			providertypes.NewErrorWithCode(err, providertypes.ErrorInvalidResponse),
		)
	}

	// set price to priceAdjustmentMax if its 1.00
	if big.NewFloat(1.00).Cmp(price) == 0 {
		price = new(big.Float).SetFloat64(priceAdjustmentMax)
	}
	// switch price to priceAdjustmentMin if its 0.00.
	if big.NewFloat(0.00).Cmp(price) == 0 {
		price = new(big.Float).SetFloat64(priceAdjustmentMin)
	}

	resolved := types.ResolvedPrices{
		ids[0]: types.NewPriceResult(price, time.Now().UTC()),
	}

	return types.NewPriceResponse(resolved, nil)
}

// validatePrice ensures the price is between [1.00 and 0.00].
func validatePrice(price *big.Float) error {
	if sign := price.Sign(); sign == -1 {
		return fmt.Errorf("price must be greater than 0.00")
	}

	maxPriceFloat := 1.00
	maxPrice := big.NewFloat(maxPriceFloat)
	diff := new(big.Float).Sub(maxPrice, price)
	if diff.Sign() == -1 {
		return fmt.Errorf("price exceeded %.2f", maxPriceFloat)
	}

	return nil
}
