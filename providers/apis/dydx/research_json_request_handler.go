package dydx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"encoding/json"

	dydxtypes "github.com/skip-mev/slinky/providers/apis/dydx/types"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
)

var _ apihandlers.RequestHandler = &ResearchJSONRequestHandler{}

// NewResearchJSONRequestHandler returns a new ResearchJSONRequestHandler.
func NewResearchJSONRequestHandler(rh apihandlers.RequestHandler) *ResearchJSONRequestHandler {
	return &ResearchJSONRequestHandler{
		RequestHandler: rh,
	}
}

// ResearchJSONRequestHandler handles the logic of fetching a dydx research
// json from github.
type ResearchJSONRequestHandler struct {
	apihandlers.RequestHandler
}

// Do calls the underlying RequestHandler's Do method, and interprets the http-response
// body as a dydx research json.
func (r *ResearchJSONRequestHandler) Do(ctx context.Context, url string) (*http.Response, error) {
	// make the request via the underlying RequestHandler
	resp, err := r.RequestHandler.Do(ctx, url)
	if err != nil {
		return nil, err
	}
	// we may be modifying the response body, so we need to close a reference to the body
	body := resp.Body
	defer body.Close()

	// determine if the http-response is ok, propagate the http-response
	// so that we can report status-codes via prometheus
	if isBadHTTPResponse(resp) {
		return resp, fmt.Errorf("bad http response: %s", resp.Status)
	}

	// unmarshal the response body into a dydx research json
	var research dydxtypes.ResearchJSON
	if err := json.NewDecoder(resp.Body).Decode(&research); err != nil {
		return nil, fmt.Errorf("error decoding response body as a research-json: %w", err)
	}

	allMarketParams, err := researchJSONToQueryAllMarketsParamsResponse(research)
	if err != nil {
		return nil, fmt.Errorf("error converting research json to QueryAllMarketsParamsResponse: %w", err)
	}

	bz, err := json.Marshal(allMarketParams)
	if err != nil {
		return nil, err
	}

	resp.Body = io.NopCloser(bytes.NewReader(bz))
	resp.ContentLength = int64(len(bz))
	return resp, nil
}

// isBadHTTPResponse returns true if the http-response is ok.
func isBadHTTPResponse(resp *http.Response) bool {
	return resp.StatusCode != http.StatusOK
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
		Id:                marketParam.Id,
		Pair: 			   marketParam.Pair,
		Exponent: 		   int32(marketParam.Exponent),
		MinExchanges: marketParam.MinExchanges,
		MinPriceChangePpm: marketParam.MinPriceChangePpm,
		ExchangeConfigJson: string(exchangeConfigJSONBz),
	}, nil
}
