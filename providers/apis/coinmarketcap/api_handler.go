package coinmarketcap

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
)

var _ types.PriceAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the PriceAPIDataHandler interface for CoinMarketCap, which can be used
// by a base provider. The DataHandler fetches data from the Quote Latest V2 CoinMarketCap API.
// Requests for prices are fulfilled in a single request.
type APIHandler struct {
	// api is the config for the CoinMarketCap API.
	api config.APIConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewAPIHandler returns a new CoinMarketCap PriceAPIDataHandler.
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

// CreateURL returns the URL that is used to fetch data from the CoinMarketCap API for the
// given tickers.
func (h *APIHandler) CreateURL(
	tickers []types.ProviderTicker,
) (string, error) {
	var ids []string //nolint:prealloc
	for _, ticker := range tickers {
		ids = append(ids, ticker.GetOffChainTicker())
		h.cache.Add(ticker)
	}

	if len(ids) == 0 {
		return "", fmt.Errorf("no tickers provided")
	}

	query := strings.Join(ids, ",")
	return fmt.Sprintf(Endpoint, h.api.Endpoints[0].URL, query), nil
}

// ParseResponse parses the spot price HTTP response from the CoinMarketCap API and returns
// the resulting price(s).
func (h *APIHandler) ParseResponse(
	tickers []types.ProviderTicker,
	resp *http.Response,
) types.PriceResponse {
	// Parse the response into a CoinBaseResponse.
	var result CoinMarketCapResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToDecode),
		)
	}

	if result.Status.ErrorCode != 0 {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(
				fmt.Errorf("coinmarketcap error: %s", result.Status.ErrorMessage),
				providertypes.ErrorAPIGeneral,
			),
		)
	}

	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	// Convert the float64 price into a big.Float.
	for id, tickerResponse := range result.Data {
		ticker, exists := h.cache.FromOffChainTicker(id)
		if !exists {
			continue
		}

		if len(tickerResponse.Quote) == 0 {
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(
					fmt.Errorf("no quote for ticker %s", ticker.GetOffChainTicker()),
					providertypes.ErrorNoResponse,
				),
			}
			continue
		}

		quote, exists := tickerResponse.Quote[DefaultQuoteDenom]
		if !exists {
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(
					fmt.Errorf("no USD quote for ticker %s", ticker.GetOffChainTicker()),
					providertypes.ErrorNoResponse,
				),
			}
			continue
		}

		resolved[ticker] = types.NewPriceResult(big.NewFloat(quote.Price), time.Now().UTC())
	}

	return types.NewPriceResponse(resolved, unresolved)
}
