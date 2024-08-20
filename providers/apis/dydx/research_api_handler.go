package dydx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/pkg/arrays"
	"github.com/skip-mev/connect/v2/providers/apis/coinmarketcap"
	dydxtypes "github.com/skip-mev/connect/v2/providers/apis/dydx/types"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	"github.com/skip-mev/connect/v2/service/clients/marketmap/types"
)

var _ types.MarketMapAPIDataHandler = (*ResearchAPIHandler)(nil)

// NewResearchAPIHandler returns a new MarketMap MarketMapAPIDataHandler.
func NewResearchAPIHandler(
	logger *zap.Logger,
	api config.APIConfig,
) (*ResearchAPIHandler, error) {
	if api.Name != ResearchAPIHandlerName && api.Name != ResearchCMCAPIHandlerName {
		return nil, fmt.Errorf("expected api config name %s or %s, got %s", ResearchAPIHandlerName, ResearchCMCAPIHandlerName, api.Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api config for %s is not enabled", ResearchAPIHandlerName)
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config for %s: %w", ResearchAPIHandlerName, err)
	}

	return &ResearchAPIHandler{
		APIHandler: APIHandler{
			api:    api,
			logger: logger,
		},
		url:     api.Endpoints[1].URL, // We assume the first URL is the endpoint of the dydx mainnet
		onlyCMC: api.Name == ResearchCMCAPIHandlerName,
	}, nil
}

// ResearchAPIHandler is a subclass of the dydx_chain.ResearchAPIHandler. It interprets the dydx ResearchJSON
// as a market-map.
type ResearchAPIHandler struct {
	APIHandler

	// url is the URL to query for the market map.
	url string
	// onlyCMC is a flag that indicates whether the handler should only return CoinMarketCap markets.
	onlyCMC bool
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
	respMarketParams, err := h.researchJSONToQueryAllMarketsParamsResponse(research)
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
	marketMap, err := ConvertMarketParamsToMarketMap(respMarketParams)
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
func (h *ResearchAPIHandler) researchJSONToQueryAllMarketsParamsResponse(research dydxtypes.ResearchJSON) (dydxtypes.QueryAllMarketParamsResponse, error) {
	// iterate over all entries in the research json + unmarshal it's market-params
	resp := dydxtypes.QueryAllMarketParamsResponse{}
	for _, market := range research {
		if h.onlyCMC && market.CMCID < 0 {
			continue
		}

		// convert the dydx research json market-param into a MarketParam struct
		marketParam, err := h.marketParamFromResearchJSONMarketParam(market)
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
func (h *ResearchAPIHandler) marketParamFromResearchJSONMarketParam(marketParam dydxtypes.Params) (dydxtypes.MarketParam, error) {
	var exchangeConfigJSON dydxtypes.ExchangeConfigJson
	if !h.onlyCMC {
		exchangeConfigJSON = dydxtypes.ExchangeConfigJson{
			Exchanges: marketParam.ExchangeConfigJSON,
		}
	} else {
		exchange := dydxtypes.ExchangeMarketConfigJson{
			ExchangeName: coinmarketcap.Name,
			Ticker:       fmt.Sprintf("%d", marketParam.CMCID),
		}
		exchangeConfigJSON = dydxtypes.ExchangeConfigJson{
			Exchanges: []dydxtypes.ExchangeMarketConfigJson{exchange},
		}
	}
	// marshal to a json string
	exchangeConfigJSONBz, err := json.Marshal(exchangeConfigJSON)
	if err != nil {
		return dydxtypes.MarketParam{}, err
	}

	var minExchanges uint32
	if h.onlyCMC {
		minExchanges = 1
	} else {
		minExchanges = marketParam.MinExchanges
	}

	return dydxtypes.MarketParam{
		Id:                 marketParam.ID,
		Pair:               marketParam.Pair,
		Exponent:           int32(marketParam.Exponent),
		MinExchanges:       minExchanges,
		MinPriceChangePpm:  marketParam.MinPriceChangePpm,
		ExchangeConfigJson: string(exchangeConfigJSONBz),
	}, nil
}
