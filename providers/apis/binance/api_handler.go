package binance

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/skip-mev/slinky/pkg/math"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/api/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of the provider.
	Name = "binance"
)

var _ handlers.APIDataHandler[oracletypes.CurrencyPair, *big.Int] = (*APIHandler)(nil)

// APIHandler implements the APIHandler interface for Binance.
// for more information about the Binance API, refer to the following link:
// https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#public-api-endpoints
type APIHandler struct {
	Config
	BaseURL string
}

// NewBinanceAPIHandler returns a new Binance API handler.
func NewBinanceAPIHandler(
	providerCfg config.ProviderConfig,
) (*APIHandler, error) {
	if providerCfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, providerCfg.Name)
	}

	cfg, err := ReadBinanceConfigFromFile(providerCfg.Path)
	if err != nil {
		return nil, err
	}

	return &APIHandler{
		Config:  cfg,
		BaseURL: BaseURL,
	}, nil
}

// CreateURL returns the URL that is used to fetch data from the Binance API for the
// given currency pairs.
func (h *APIHandler) CreateURL(
	cps []oracletypes.CurrencyPair,
) (string, error) {
	var cpStrings string

	for _, cp := range cps {
		base, ok := h.SupportedBases[cp.Base]
		if !ok {
			continue
		}

		quote, ok := h.SupportedQuotes[cp.Quote]
		if !ok {
			continue
		}

		cpStrings += fmt.Sprintf("%s%s%s%s%s", Quotation, base, quote, Quotation, Separator)
	}

	if len(cpStrings) == 0 {
		return "", fmt.Errorf("empty url created. invalid or no currency pairs were provided")
	}

	// remove last comma from list
	cpStrings = strings.TrimSuffix(cpStrings, Separator)
	return fmt.Sprintf(h.BaseURL, LeftBracket, cpStrings, RightBracket), nil
}

func (h *APIHandler) ParseResponse(
	cps []oracletypes.CurrencyPair,
	resp *http.Response,
) providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int] {
	// Parse the response into a BinanceResponse.
	result, err := h.Decode(resp)
	if err != nil {
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

		cpMap[fmt.Sprintf("%s%s", base, quote)] = cp
	}

	// Filter out the responses that are not expected.
	for _, data := range result {
		cp, ok := cpMap[data.Symbol]
		if !ok {
			continue
		}

		price, err := math.Float64StringToBigInt(data.Price, cp.Decimals())
		if err != nil {
			return providertypes.NewGetResponseWithErr[oracletypes.CurrencyPair, *big.Int](cps, err)
		}

		resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now())
		delete(cpMap, data.Symbol)
	}

	// If there are any currency pairs that were not resolved, we need to add them
	// to the unresolved map.
	for _, cp := range cpMap {
		unresolved[cp] = fmt.Errorf("currency pair %s did not get a response", cp.String())
	}

	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved)
}

func (h *APIHandler) Atomic() bool {
	return true
}

// Name returns the name of the handler.
func (h *APIHandler) Name() string {
	return Name
}

// Decode decodes the given http response into a BinanceResponse.
func (h *APIHandler) Decode(resp *http.Response) (Response, error) {
	// Parse the response into a BinanceResponse.
	var result Response
	err := json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
