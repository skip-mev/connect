package geckoterminal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ types.PriceAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the PriceAPIDataHandler interface for GeckoTerminal.
type APIHandler struct {
	// marketCfg is the config for the GeckoTerminal API.
	market types.ProviderMarketMap
	// apiCfg is the config for the GeckoTerminal API.
	api config.APIConfig
}

// NewAPIHandler returns a new GeckoTerminal PriceAPIDataHandler.
func NewAPIHandler(
	market types.ProviderMarketMap,
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

// CreateURL returns the URL that is used to fetch data from the GeckoTerminal API for the
// given tickers. Note that the GeckoTerminal API supports fetching multiple spot prices
// iff they are all on the same chain.
func (h *APIHandler) CreateURL(
	tickers []mmtypes.Ticker,
) (string, error) {
	if len(tickers) > MaxNumberOfTickers {
		return "", fmt.Errorf("expected at most %d tickers, got %d", MaxNumberOfTickers, len(tickers))
	}

	addresses := make([]string, len(tickers))
	for i, ticker := range tickers {
		cfg, ok := h.market.TickerConfigs[ticker]
		if !ok {
			return "", fmt.Errorf("no config for ticker %s", ticker.String())
		}

		addresses[i] = cfg.OffChainTicker
	}

	return fmt.Sprintf(h.api.URL, strings.Join(addresses, ",")), nil
}

// ParseResponse parses the response from the GeckoTerminal API. The response is expected
// to contain multiple spot prices for a given token address. Note that all of the tokens
// are shared on the same chain.
func (h *APIHandler) ParseResponse(
	tickers []mmtypes.Ticker,
	resp *http.Response,
) types.PriceResponse {
	// Parse the response.
	var result GeckoTerminalResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return types.NewPriceResponseWithErr(tickers, err)
	}

	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	data := result.Data
	if data.Type != ExpectedResponseType {
		return types.NewPriceResponseWithErr(tickers, fmt.Errorf("expected type %s, got %s", ExpectedResponseType, data.Type))
	}

	// Filter out the responses that are not expected.
	attributes := data.Attributes
	for address, price := range attributes.TokenPrices {
		ticker, ok := h.market.OffChainMap[address]
		if !ok {
			return types.NewPriceResponseWithErr(tickers, fmt.Errorf("no ticker for address %s", address))
		}

		// Convert the price to a big.Int.
		price, err := math.Float64StringToBigInt(price, ticker.Decimals)
		if err != nil {
			unresolved[ticker] = fmt.Errorf("failed to convert price to big.Int: %w", err)
			continue
		}

		resolved[ticker] = types.NewPriceResult(price, time.Now())
	}

	// Add all of the expected tickers that did not return a response to the unresolved
	// map.
	for _, ticker := range tickers {
		if _, resolvedOk := resolved[ticker]; !resolvedOk {
			unresolved[ticker] = fmt.Errorf("received no price response")
		}
	}

	return types.NewPriceResponse(resolved, unresolved)
}
