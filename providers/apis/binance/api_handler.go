package binance

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/skip-mev/slinky/pkg/math"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ types.PriceAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the PriceAPIDataHandler interface for Binance.
// for more information about the Binance API, refer to the following link:
// https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#public-api-endpoints
type APIHandler struct {
	// market is the config for the Binance API.
	market mmtypes.MarketConfig
	// api is the config for the Binance API.
	api config.APIConfig
}

// NewAPIHandler returns a new Binance Price API Data Handler.
func NewAPIHandler(
	market mmtypes.MarketConfig,
	api config.APIConfig,
) (types.PriceAPIDataHandler, error) {
	if err := market.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid market config for %s: %w", Name, err)
	}

	if market.Name != Name {
		return nil, fmt.Errorf("expected market config name %s, got %s", Name, market.Name)
	}

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
		market: market,
		api:    api,
	}, nil
}

// CreateURL returns the URL that is used to fetch data from the Binance API for the
// given tickers.
func (h *APIHandler) CreateURL(
	tickers []mmtypes.Ticker,
) (string, error) {
	var tickerStrings string
	for _, ticker := range tickers {
		market, ok := h.market.TickerConfigs[ticker.String()]
		if !ok {
			return "", fmt.Errorf("ticker %s not found in market config", ticker.String())
		}

		tickerStrings += fmt.Sprintf("%s%s%s%s", Quotation, market.OffChainTicker, Quotation, Separator)
	}

	if len(tickerStrings) == 0 {
		return "", fmt.Errorf("empty url created. invalid or no ticker were provided")
	}

	return fmt.Sprintf(
		h.api.URL,
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
) types.PriceResponse {
	// Parse the response into a BinanceResponse.
	result, err := Decode(resp)
	if err != nil {
		return providertypes.NewGetResponseWithErr[mmtypes.Ticker, *big.Int](tickers, err)
	}

	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	inverted := h.market.Invert()
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

		resolved[market.Ticker] = providertypes.NewResult(price, time.Now())
	}

	// Add currency pairs that received no response to the unresolved map.
	for _, ticker := range tickers {
		_, resolvedOk := resolved[ticker]
		_, unresolvedOk := unresolved[ticker]

		if !resolvedOk && !unresolvedOk {
			unresolved[ticker] = fmt.Errorf("no response")
		}
	}

	return providertypes.NewGetResponse(resolved, unresolved)
}
