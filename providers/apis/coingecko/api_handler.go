package coingecko

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/math"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/base/api/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var _ handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int] = (*APIHandler)(nil)

// APIHandler implements the Base Provider API handler interface for CoinGecko.
// This provider can be configured to support API based fetching, however, the provider
// does not require it.
type APIHandler struct {
	// cfg is the provider config.
	cfg config.ProviderConfig
}

// NewAPIHandler returns a new CoinGecko API handler.
func NewAPIHandler(
	cfg config.ProviderConfig,
) (handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int], error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid provider config %w", err)
	}

	if !cfg.API.Enabled {
		return nil, fmt.Errorf("api is not enabled for provider %s", cfg.Name)
	}

	if cfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, cfg.Name)
	}

	return &APIHandler{
		cfg: cfg,
	}, nil
}

// CreateURL returns the URL that is used to fetch data from the CoinGecko API for the
// given currency pairs. The CoinGecko API supports fetching spot prices for multiple
// currency pairs in a single request. The URL that is generated automatically populates
// the API key if it is set.
func (h *APIHandler) CreateURL(
	cps []slinkytypes.CurrencyPair,
) (string, error) {
	// Create a list of base currencies and quote currencies.
	bases, quotes, err := h.getUniqueBaseAndQuoteDenoms(cps)
	if err != nil {
		return "", err
	}

	// This creates the endpoint that needs to be requested regardless of whether
	// an API key is set.
	pricesEndPoint := fmt.Sprintf(PairPriceEndpoint, bases, quotes)
	finalEndpoint := fmt.Sprintf("%s%s", pricesEndPoint, Precision)

	// Otherwise, we just return the base url with the endpoint.
	return fmt.Sprintf("%s%s", h.cfg.API.URL, finalEndpoint), nil
}

// ParseResponse parses the response from the CoinGecko API. The response is expected
// to match every base currency with every quote currency. As such, we need to filter
// out the responses that are not expected. Note that the response will only return
// a response for the inputted currency pairs.
func (h *APIHandler) ParseResponse(
	cps []slinkytypes.CurrencyPair,
	resp *http.Response,
) providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int] {
	// Parse the response.
	var result CoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return providertypes.NewGetResponseWithErr[slinkytypes.CurrencyPair, *big.Int](cps, err)
	}

	var (
		resolved   = make(map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int])
		unresolved = make(map[slinkytypes.CurrencyPair]error)
	)

	// Map each of the currency pairs for easy lookup.
	configCPs := config.NewMarketConfig()
	for _, cp := range cps {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.String()]
		if !ok {
			continue
		}

		configCPs.CurrencyPairToMarketConfigs[cp.String()] = market
	}

	// Filter out the responses that are not expected.
	for base, quotes := range result {
		for quote, price := range quotes {
			// The ticker is represented as base/quote.
			ticker := fmt.Sprintf("%s%s%s", base, TickerSeparator, quote)

			// If the ticker is not configured, we skip it.
			market, ok := h.cfg.Market.TickerToMarketConfigs[ticker]
			if !ok {
				continue
			}

			// Resolve the price.
			cp := market.CurrencyPair
			price := math.Float64ToBigInt(price, cp.Decimals())
			resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now())
			delete(configCPs.CurrencyPairToMarketConfigs, cp.String())
		}
	}

	// If there are any currency pairs that were not resolved, we need to add them
	// to the unresolved map.
	for _, market := range configCPs.CurrencyPairToMarketConfigs {
		unresolved[market.CurrencyPair] = fmt.Errorf("currency pair %s did not get a response", market.CurrencyPair.String())
	}

	return providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](resolved, unresolved)
}
