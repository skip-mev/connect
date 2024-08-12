package polymarket

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

const (
	// Name is the name of the Polymarket provider.
	Name = "polymarket_api"

	// URL is the default base URL of the Polymarket CLOB API. It uses the midpoint endpoint with a given token ID.
	URL = "https://clob.polymarket.com/markets/%s"

	// priceAdjustmentMax is the value the price gets set to in the event of price == 1.00.
	priceAdjustmentMin = 0.00001
)

var _ types.PriceAPIDataHandler = (*APIHandler)(nil)

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

	if len(api.Endpoints) != 1 {
		return nil, fmt.Errorf("invalid polymarket endpoint config: expected 1 endpoint got %d", len(api.Endpoints))
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
	marketID, _, err := getMarketAndTokenFromTicker(ids[0])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(h.api.Endpoints[0].URL, marketID), nil
}

type TokenData struct {
	TokenID string  `json:"token_id"`
	Outcome string  `json:"outcome"`
	Price   float64 `json:"price"`
}

type MarketData struct {
	ConditionID   string      `json:"condition_id"`
	Question      string      `json:"question"`
	MarketSlug    string      `json:"market_slug"`
	EventSlug     string      `json:"event_slug"`
	Image         string      `json:"image"`
	Tokens        []TokenData `json:"tokens"`
	RewardsConfig []struct {
		AssetAddress string `json:"asset_address"`
		StartDate    string `json:"start_date"`
		EndDate      string `json:"end_date"`
		ID           int    `json:"id"`
		RatePerDay   int    `json:"rate_per_day"`
		TotalRewards int    `json:"total_rewards"`
		TotalDays    int    `json:"total_days"`
	} `json:"rewards_config"`
	RewardsMaxSpread float64 `json:"rewards_max_spread"`
	RewardsMinSize   int     `json:"rewards_min_size"`
}

type MarketsResponse struct {
	Data       []MarketData `json:"data"`
	NextCursor string       `json:"next_cursor"`
	Limit      int          `json:"limit"`
	Count      int          `json:"count"`
}

// ParseResponse parses the HTTP response from either the `/price` or `/midpoint` endpoint of the Polymarket API endpoint and returns
// the resulting data.
func (h APIHandler) ParseResponse(ids []types.ProviderTicker, response *http.Response) types.PriceResponse {
	if len(ids) != 1 {
		return priceResponseError(ids, fmt.Errorf("expected 1 ticker, got %d", len(ids)), providertypes.ErrorInvalidResponse)
	}

	var result MarketsResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return priceResponseError(ids, fmt.Errorf("failed to decode market response: %w", err), providertypes.ErrorFailedToDecode)
	}

	if len(result.Data) != 1 {
		return priceResponseError(
			ids,
			fmt.Errorf("expected 1 market in response, got %d", len(result.Data)),
			providertypes.ErrorInvalidResponse,
		)
	}

	_, tokenID, err := getMarketAndTokenFromTicker(ids[0])
	if err != nil {
		return priceResponseError(ids, err, providertypes.ErrorAPIGeneral)
	}

	var tokenData *TokenData
	for _, token := range result.Data[0].Tokens {
		if token.TokenID == tokenID {
			tokenData = &token
			break
		}
	}

	if tokenData == nil {
		return priceResponseError(ids, fmt.Errorf("token ID %s not found in response", tokenID), providertypes.ErrorInvalidResponse)
	}

	price := new(big.Float).SetFloat64(tokenData.Price)

	// switch price to priceAdjustmentMin if its 0.00.
	if big.NewFloat(0.00).Cmp(price) == 0 {
		price = new(big.Float).SetFloat64(priceAdjustmentMin)
	}

	resolved := types.ResolvedPrices{
		ids[0]: types.NewPriceResult(price, time.Now().UTC()),
	}

	return types.NewPriceResponse(resolved, nil)
}

func priceResponseError(ids []types.ProviderTicker, err error, code providertypes.ErrorCode) providertypes.GetResponse[types.ProviderTicker, *big.Float] {
	return types.NewPriceResponseWithErr(
		ids,
		providertypes.NewErrorWithCode(err, code),
	)
}

func getMarketAndTokenFromTicker(t types.ProviderTicker) (marketID string, tokenID string, err error) {
	split := strings.Split(t.GetOffChainTicker(), "/")
	if len(split) != 2 {
		return "", "", fmt.Errorf("expected ticker format market_id/token_id, got: %s", t.GetOffChainTicker())
	}
	return split[0], split[1], nil
}
