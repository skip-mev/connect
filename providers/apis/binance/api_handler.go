package binance

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/skip-mev/slinky/pkg/math"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/constants"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ constants.PriceAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the APIHandler interface for Binance.
// for more information about the Binance API, refer to the following link:
// https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#public-api-endpoints
type APIHandler struct {
	// marketCfg is the config for the Binance API.
	marketCfg mmtypes.MarketConfig
	// apiCfg is the config for the Binance API.
	apiCfg config.APIConfig
}

// NewAPIHandler returns a new Binance Price API Data Handler.
func NewAPIHandler(
	marketCfg mmtypes.MarketConfig,
	apiCfg config.APIConfig,
) (constants.PriceAPIDataHandler, error) {
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
// given tickers.
func (h *APIHandler) CreateURL(
	tickers []mmtypes.Ticker,
) (string, error) {
	var tickerStrings string
	for _, ticker := range tickers {
		market, ok := h.marketCfg.TickerConfigs[ticker.String()]
		if !ok {
			return "", fmt.Errorf("ticker %s not found in market config", ticker.String())
		}

		tickerStrings += fmt.Sprintf("%s%s%s%s", Quotation, market.OffChainTicker, Quotation, Separator)
	}

	if len(tickerStrings) == 0 {
		return "", fmt.Errorf("empty url created. invalid or no ticker were provided")
	}

	return fmt.Sprintf(
		h.apiCfg.URL,
		LeftBracket,
		strings.TrimSuffix(tickerStrings, Separator),
		RightBracket,
	), nil
}

// ParseResponse parses the response from the Binance API and returns a GetResponse. Each
// of the tickers supplied will get a response or an error.
func (h *APIHandler) ParseResponse(
	tickers []mmtypes.Ticker,
	resp *http.Response,
) providertypes.GetResponse[mmtypes.Ticker, *big.Int] {
	// Parse the response into a BinanceResponse.
	result, err := Decode(resp)
	if err != nil {
		return providertypes.NewGetResponseWithErr[mmtypes.Ticker, *big.Int](tickers, err)
	}

	var (
		resolved   = make(map[mmtypes.Ticker]providertypes.Result[*big.Int])
		unresolved = make(map[mmtypes.Ticker]error)
	)

	inverted := h.marketCfg.Invert()
	for _, data := range result {
		// Filter out the responses that are not expected.
		market, ok := inverted[data.Symbol]
		if !ok {
			continue
		}

		price, err := math.Float64StringToBigInt(data.Price, market.Ticker.Decimals)
		if err != nil {
			unresolved[market.Ticker] = fmt.Errorf("failed to convert price %s to big.Int: %w", data.Price, err)
			continue
		}

		resolved[market.Ticker] = providertypes.NewResult[*big.Int](price, time.Now())
	}

	// Add currency pairs that received no response to the unresolved map.
	for _, ticker := range tickers {
		_, resolvedOk := resolved[ticker]
		_, unresolvedOk := unresolved[ticker]

		if !resolvedOk && !unresolvedOk {
			unresolved[ticker] = fmt.Errorf("no response")
		}
	}

	return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unresolved)
}
