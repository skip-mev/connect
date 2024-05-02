package dydx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/arrays"
	dydxtypes "github.com/skip-mev/slinky/providers/apis/dydx/types"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/service/clients/marketmap/types"
)

var _ types.MarketMapAPIDataHandler = (*ResearchAPIHandler)(nil)

// NewResearchAPIHandler returns a new MarketMap MarketMapAPIDataHandler.
func NewResearchAPIHandler(
	logger *zap.Logger,
	api config.APIConfig,
) (*ResearchAPIHandler, error) {
	if api.Name != ResearchAPIHandlerName {
		return nil, fmt.Errorf("expected api config name %s, got %s", ResearchAPIHandlerName, api.Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api config for %s is not enabled", ResearchAPIHandlerName)
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config for %s: %w", ResearchAPIHandlerName, err)
	}

	// expect a single endpoint
	if len(api.Endpoints) != 1 {
		return nil, fmt.Errorf("expected one endpoint, got %d", len(api.Endpoints))
	}

	return &ResearchAPIHandler{
		APIHandler: APIHandler{
			api:    api,
			logger: logger,
		},
		url: api.Endpoints[0].URL,
	}, nil
}

// ResearchAPIHandler is a subclass of the dydx_chain.ResearchAPIHandler. It interprets the dydx ResearchJSON
// as a market-map.
type ResearchAPIHandler struct {
	APIHandler

	// url is the URL to query for the market map.
	url string
}

// CreateURL returns a static url (the url of the first configured endpoint). If the dydx chain is not
// configured in the request, an error is returned.
func (h *ResearchAPIHandler) CreateURL(chains []types.Chain) (string, error) {
	// expect at least one chain to be a dydx chain
	if _, ok := arrays.CheckEntryInArray(types.Chain{
		ChainID: ChainID,
	}, chains); !ok {
		return "", fmt.Errorf("dydx chain is not configured in request for the dydx research json")
	}

	return h.url, nil
}

// ParseResponse parses the response from the dydx ResearchJSON API into a MarketMap, and
// unmarshals the market-map in accordance with the underlying dydx ResearchAPIHandler.
func (h *ResearchAPIHandler) ParseResponse(
	chains []types.Chain,
	resp *http.Response,
) types.MarketMapResponse {
	// expect at least one chain to be a dydx chain
	chain, ok := arrays.CheckEntryInArray(types.Chain{
		ChainID: ChainID,
	}, chains)
	if !ok {
		h.logger.Error("dydx chain is not configured in request for the dydx research json")
		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("expected one chain, got %d", len(chains)),
				providertypes.ErrorInvalidAPIChains,
			),
		)
	}

	// parse the response
	// unmarshal the response body into a dydx research json
	var research dydxtypes.ResearchJSON
	if err := json.NewDecoder(resp.Body).Decode(&research); err != nil {
		h.logger.Error("failed to parse dydx research json response", zap.Error(err))
		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("failed to parse dydx research json response: %w", err),
				providertypes.ErrorFailedToDecode,
			),
		)
	}

	// convert the dydx research json into a QueryAllMarketsParamsResponse
	respMarketParams, err := researchJSONToQueryAllMarketsParamsResponse(research)
	if err != nil {
		h.logger.Error("failed to convert dydx research json into QueryAllMarketsParamsResponse", zap.Error(err))
		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("failed to convert dydx research json into QueryAllMarketsParamsResponse: %w", err),
				providertypes.ErrorFailedToDecode,
			),
		)
	}

	// convert the response to a market-map
	marketMap, err := h.ConvertMarketParamsToMarketMap(respMarketParams)
	if err != nil {
		h.logger.Error("failed to convert QueryAllMarketsParamsResponse into MarketMap", zap.Error(err))
		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("failed to convert QueryAllMarketsParamsResponse into MarketMap: %w", err),
				providertypes.ErrorFailedToDecode,
			),
		)
	}

	// resolve the response under the dydx chain
	resolved := make(types.ResolvedMarketMap)
	resolved[chain] = types.NewMarketMapResult(&marketMap, time.Now())

	h.logger.Debug("successfully resolved dydx research json into a market map", zap.Int("num_markets", len(marketMap.MarketMap.Markets)))
	return types.NewMarketMapResponse(resolved, nil)
}

// researchJSONToQueryAllMarketsParamsResponse converts a dydx research json into a
// QueryAllMarketsParamsResponse.
func researchJSONToQueryAllMarketsParamsResponse(research dydxtypes.ResearchJSON) (dydxtypes.QueryAllMarketParamsResponse, error) {
	// iterate over all entries in the research json + unmarshal it's market-params
	resp := dydxtypes.QueryAllMarketParamsResponse{}
	for _, market := range research {
		researchMarketParam, ok := market[dydxtypes.MarketParamIndex]
		if !ok {
			return dydxtypes.QueryAllMarketParamsResponse{}, fmt.Errorf("market %v does not have params", market)
		}

		// convert the dydx research json market-param into a MarketParam struct
		marketParam, err := marketParamFromResearchJSONMarketParam(researchMarketParam)
		if err != nil {
			return dydxtypes.QueryAllMarketParamsResponse{}, err
		}

		// unmarshal the market-params into a MarketParam struct
		resp.MarketParams = append(resp.MarketParams, marketParam)
	}

	return resp, nil
}

// marketParamFromResearchJSONMarketParam converts a dydx research json market-param
// into a MarketParam struct.
func marketParamFromResearchJSONMarketParam(marketParam dydxtypes.ResearchJSONMarketParam) (dydxtypes.MarketParam, error) {
	exchangeConfigJSON := dydxtypes.ExchangeConfigJson{
		Exchanges: marketParam.ExchangeConfigJSON,
	}
	// marshal to a json string
	exchangeConfigJSONBz, err := json.Marshal(exchangeConfigJSON)
	if err != nil {
		return dydxtypes.MarketParam{}, err
	}

	return dydxtypes.MarketParam{
		Id:                 marketParam.ID,
		Pair:               marketParam.Pair,
		Exponent:           int32(marketParam.Exponent),
		MinExchanges:       marketParam.MinExchanges,
		MinPriceChangePpm:  marketParam.MinPriceChangePpm,
		ExchangeConfigJson: string(exchangeConfigJSONBz),
	}, nil
}
