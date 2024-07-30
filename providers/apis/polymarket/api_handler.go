package polymarket

import (
	"encoding/json"
	"fmt"
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

	// URL is the base URL of the Polymarket CLOB API endpoint for the Price of a given token ID.
	URL = "https://clob.polymarket.com/price?token_id=%s&side=BUY"

	// priceAdjustment is the value the price gets set to in the event of price == 1.00
	priceAdjustment = 0.9999999
)

var _ types.PriceAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the PriceAPIDataHandler interface for Polymarket, which can be used
// by a base provider. The handler fetches data from the `/price` endpoint.
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
// given ticker. Since the price endpoint is automatically denominated in USD, only one ID is expected to be passed
// into this method.
func (h APIHandler) CreateURL(ids []types.ProviderTicker) (string, error) {
	if len(ids) != 1 {
		return "", fmt.Errorf("expected 1 ticker, got %d", len(ids))
	}
	return fmt.Sprintf(h.api.Endpoints[0].URL, ids[0].GetOffChainTicker()), nil
}

// ResponseBody is the response structure for the `/price` endpoint of the Polymarket API.
type ResponseBody struct {
	Price string `json:"price"`
}

// ParseResponse parses the HTTP response from the `/price` Polymarket API endpoint and returns
// the resulting price.
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

	var result ResponseBody
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return types.NewPriceResponseWithErr(
			ids,
			providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToDecode),
		)
	}

	price, ok := new(big.Float).SetString(result.Price)
	if !ok {
		return types.NewPriceResponseWithErr(
			ids,
			providertypes.NewErrorWithCode(fmt.Errorf("failed to convert %q to float", result.Price), providertypes.ErrorFailedToDecode),
		)
	}
	if err := validatePrice(price); err != nil {
		return types.NewPriceResponseWithErr(
			ids,
			providertypes.NewErrorWithCode(err, providertypes.ErrorInvalidResponse),
		)
	}

	// we don't ever want to return 1.00. Set to priceAdjustment.
	if big.NewFloat(1.00).Cmp(price) == 0 {
		price = new(big.Float).SetFloat64(priceAdjustment)
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
