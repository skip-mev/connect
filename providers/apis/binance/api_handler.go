package binance

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/skip-mev/slinky/pkg/math"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/base/api/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var _ handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int] = (*APIHandler)(nil)

// APIHandler implements the APIHandler interface for Binance.
// for more information about the Binance API, refer to the following link:
// https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#public-api-endpoints
type APIHandler struct {
	// cfg is the config for the Binance API.
	cfg config.ProviderConfig
}

// NewAPIHandler returns a new Binance API handler.
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

// CreateURL returns the URL that is used to fetch data from the Binance API for the
// given currency pairs.
func (h *APIHandler) CreateURL(
	cps []slinkytypes.CurrencyPair,
) (string, error) {
	var cpStrings string

	for _, cp := range cps {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.String()]
		if !ok {
			continue
		}

		cpStrings += fmt.Sprintf("%s%s%s%s", Quotation, market.Ticker, Quotation, Separator)
	}

	if len(cpStrings) == 0 {
		return "", fmt.Errorf("empty url created. invalid or no currency pairs were provided")
	}

	// remove last comma from list
	cpStrings = strings.TrimSuffix(cpStrings, Separator)
	return fmt.Sprintf(h.cfg.API.URL, LeftBracket, cpStrings, RightBracket), nil
}

func (h *APIHandler) ParseResponse(
	cps []slinkytypes.CurrencyPair,
	resp *http.Response,
) providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int] {
	// Parse the response into a BinanceResponse.
	result, err := Decode(resp)
	if err != nil {
		return providertypes.NewGetResponseWithErr[slinkytypes.CurrencyPair, *big.Int](cps, err)
	}

	var (
		resolved   = make(map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int])
		unresolved = make(map[slinkytypes.CurrencyPair]error)
	)

	// Determine of the provided currency pairs which are supported by the Binance API.
	configuredCps := config.NewMarketConfig()
	for _, cp := range cps {
		market, ok := h.cfg.Market.CurrencyPairToMarketConfigs[cp.String()]
		if !ok {
			continue
		}

		configuredCps.CurrencyPairToMarketConfigs[cp.String()] = market
	}

	// Filter out the responses that are not expected.
	for _, data := range result {
		market, ok := h.cfg.Market.TickerToMarketConfigs[data.Symbol]
		if !ok {
			continue
		}

		cp := market.CurrencyPair
		price, err := math.Float64StringToBigInt(data.Price, cp.Decimals())
		if err != nil {
			return providertypes.NewGetResponseWithErr[slinkytypes.CurrencyPair, *big.Int](cps, err)
		}

		resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now())
		delete(configuredCps.CurrencyPairToMarketConfigs, cp.String())
	}

	// If there are any currency pairs that were not resolved, return an error.
	for _, market := range configuredCps.CurrencyPairToMarketConfigs {
		cp := market.CurrencyPair
		unresolved[cp] = fmt.Errorf("currency pair %s did not get a response", cp.String())
	}

	return providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](resolved, unresolved)
}
