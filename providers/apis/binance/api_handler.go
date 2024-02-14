package binance

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/skip-mev/slinky/pkg/math"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/api/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ handlers.APIDataHandler[mmtypes.Ticker, *big.Int] = (*APIHandler)(nil)

// APIHandler implements the APIHandler interface for Binance.
// for more information about the Binance API, refer to the following link:
// https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#public-api-endpoints
type APIHandler struct {
	// marketCfg is the config for the Binance API.
	marketCfg mmtypes.MarketConfig
	// apiCfg is the config for the Binance API.
	apiCfg config.APIConfig
}

// NewAPIHandler returns a new Binance API handler.
func NewAPIHandler(
	marketCfg mmtypes.MarketConfig,
	apiCfg config.APIConfig,
) (handlers.APIDataHandler[mmtypes.Ticker, *big.Int], error) {
	if err := marketCfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid provider config %w", err)
	}

	if marketCfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, marketCfg.Name)
	}

	if apiCfg.Name != Name {
		return nil, fmt.Errorf("expected api config name %s, got %s", Name, apiCfg.Name)
	}

	if !apiCfg.Enabled {
		return nil, fmt.Errorf("api config for %s is not enabled", Name)

	}

	if err := apiCfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config %w", err)
	}

	return &APIHandler{
		marketCfg: marketCfg,
		apiCfg:    apiCfg,
	}, nil
}

// CreateURL returns the URL that is used to fetch data from the Binance API for the
// given currency pairs.
func (h *APIHandler) CreateURL(
	cps []mmtypes.Ticker,
) (string, error) {
	var cpStrings string

	for _, cp := range cps {
		market, ok := h.marketCfg.TickerConfigs[cp.String()]
		if !ok {
			return "", fmt.Errorf("currency pair %s not found in market config", cp.String())
		}

		cpStrings += fmt.Sprintf("%s%s%s%s", Quotation, market.OffChainTicker, Quotation, Separator)
	}

	if len(cpStrings) == 0 {
		return "", fmt.Errorf("empty url created. invalid or no currency pairs were provided")
	}

	// remove last comma from list
	cpStrings = strings.TrimSuffix(cpStrings, Separator)
	return fmt.Sprintf(h.apiCfg.URL, LeftBracket, cpStrings, RightBracket), nil
}

func (h *APIHandler) ParseResponse(
	cps []mmtypes.Ticker,
	resp *http.Response,
) providertypes.GetResponse[mmtypes.Ticker, *big.Int] {
	// Parse the response into a BinanceResponse.
	result, err := Decode(resp)
	if err != nil {
		return providertypes.NewGetResponseWithErr[mmtypes.Ticker, *big.Int](cps, err)
	}

	var (
		resolved   = make(map[mmtypes.Ticker]providertypes.Result[*big.Int])
		unresolved = make(map[mmtypes.Ticker]error)
	)

	// Filter out the responses that are not expected.
	inverted := h.marketCfg.Invert()
	for _, data := range result {
		market, ok := inverted[data.Symbol]
		if !ok {
			continue
		}

		cp := market.Ticker
		price, err := math.Float64StringToBigInt(data.Price, cp.Decimals)
		if err != nil {
			unresolved[cp] = fmt.Errorf("failed to convert price %s to big.Int: %w", data.Price, err)
			continue
		}

		resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now())
	}

	for _, cp := range cps {
		_, resolvedOk := resolved[cp]
		_, unresolvedOk := unresolved[cp]
		if !resolvedOk && !unresolvedOk {
			unresolved[cp] = fmt.Errorf("no response")
		}
	}

	return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unresolved)
}
