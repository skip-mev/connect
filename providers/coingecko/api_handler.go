package coingecko

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of the provider.
	Name = "coingecko"
)

var _ handlers.APIDataHandler[oracletypes.CurrencyPair, *big.Int] = (*CoinGeckoAPIHandler)(nil)

// CoinGeckoAPIHandler implements the Base Provider API handler interface for CoinGecko.
// This provider can be configured to support API based fetching, however, the provider
// does not require it.
type CoinGeckoAPIHandler struct { //nolint
	// Config is the CoinGecko config.
	Config
}

// NewCoinGeckoAPIHandler returns a new CoinGecko API handler.
func NewCoinGeckoAPIHandler(
	providerCfg config.ProviderConfig,
) (*CoinGeckoAPIHandler, error) {
	if providerCfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, providerCfg.Name)
	}

	cfg, err := ReadCoinGeckoConfigFromFile(providerCfg.Path)
	if err != nil {
		return nil, err
	}

	return &CoinGeckoAPIHandler{
		Config: cfg,
	}, nil
}

// CreateURL returns the URL that is used to fetch data from the CoinGecko API for the
// given currency pairs. The CoinGecko API supports fetching spot prices for multiple
// currency pairs in a single request. The URL that is generated automatically populates
// the API key if it is set.
func (h *CoinGeckoAPIHandler) CreateURL(
	cps []oracletypes.CurrencyPair,
) (string, error) {
	// Create a list of base currencies and quote currencies.
	bases, quotes, err := h.getUniqueBaseAndQuoteDenoms(cps)
	if err != nil {
		return "", err
	}

	// This creates the endpoint that needs to be requested regardless of whether or not
	// an API key is set.
	pricesEndPoint := fmt.Sprintf(PairPriceEndpoint, bases, quotes)
	finalEndpoint := fmt.Sprintf("%s%s", pricesEndPoint, Precision)

	// If the API key is set, we need append the API url with the API key header along
	// with the API key.
	if len(h.APIKey) != 0 {
		return fmt.Sprintf("%s%s%s%s", APIURL, finalEndpoint, APIKeyHeader, h.APIKey), nil
	}

	// Otherwise, we just return the base url with the endpoint.
	return fmt.Sprintf("%s%s", BaseURL, finalEndpoint), nil
}

// Atomic returns true as the CoinGecko API is atomic i.e. returns the price of all
// currency pairs in a single request.
func (h *CoinGeckoAPIHandler) Atomic() bool {
	return true
}

// ParseResponse parses the response from the CoinGecko API. The response is expected
// to match every base currency with every quote currency. As such, we need to filter
// out the responses that are not expected. Note that the response will only return
// a response for the inputted currency pairs.
func (h *CoinGeckoAPIHandler) ParseResponse(
	cps []oracletypes.CurrencyPair,
	resp *http.Response,
) providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int] {
	// Parse the response.
	var result CoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return providertypes.NewGetResponseWithErr[oracletypes.CurrencyPair, *big.Int](cps, err)
	}

	var (
		resolved   = make(map[oracletypes.CurrencyPair]providertypes.Result[*big.Int])
		unresolved = make(map[oracletypes.CurrencyPair]error)
	)

	// Map each of the currency pairs for easy lookup.
	cpMap := make(map[string]oracletypes.CurrencyPair)
	for _, cp := range cps {
		base, ok := h.SupportedBases[cp.Base]
		if !ok {
			unresolved[cp] = fmt.Errorf("unknown base currency %s", cp.Base)
			continue
		}

		quote, ok := h.SupportedQuotes[cp.Quote]
		if !ok {
			unresolved[cp] = fmt.Errorf("unknown quote currency %s", cp.Quote)
			continue
		}

		cpMap[fmt.Sprintf("%s-%s", base, quote)] = cp
	}

	// Filter out the responses that are not expected.
	for base, quotes := range result {
		for quote, price := range quotes {
			key := fmt.Sprintf("%s-%s", base, quote)
			cp, ok := cpMap[key]
			if !ok {
				continue
			}

			price := math.Float64ToBigInt(price, cp.Decimals())
			resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now())
			delete(cpMap, key)
		}
	}

	// If there are any currency pairs that were not resolved, we need to add them
	// to the unresolved map.
	for _, cp := range cpMap {
		unresolved[cp] = fmt.Errorf("currency pair %s did not get a response", cp.String())
	}

	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved)
}

// Name returns the name of the handler.
func (h *CoinGeckoAPIHandler) Name() string {
	return Name
}
