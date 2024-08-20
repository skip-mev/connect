package dydx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	dydxtypes "github.com/skip-mev/connect/v2/providers/apis/dydx/types"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	"github.com/skip-mev/connect/v2/service/clients/marketmap/types"
)

var _ types.MarketMapAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the MarketMapAPIDataHandler interface for the dYdX prices module, which can be used
// by a base provider. This is specifically for fetching market data from the dYdX prices module, which is
// then translated to a market map.
type APIHandler struct {
	logger *zap.Logger

	// api is the api config for the dYdX market params API.
	api config.APIConfig
}

// NewAPIHandler returns a new MarketMap MarketMapAPIDataHandler.
func NewAPIHandler(
	logger *zap.Logger,
	api config.APIConfig,
) (*APIHandler, error) {
	if api.Name != Name && api.Name != SwitchOverAPIHandlerName {
		return nil, fmt.Errorf(
			"expected api config name %s or %s, got %s",
			SwitchOverAPIHandlerName,
			Name,
			api.Name,
		)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api config for %s is not enabled", Name)
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config for %s: %w", Name, err)
	}

	return &APIHandler{
		api:    api,
		logger: logger,
	}, nil
}

// CreateURL returns the URL that is used to fetch the latest market map data from the
// dYdX prices module.
func (h *APIHandler) CreateURL(chains []types.Chain) (string, error) {
	if len(chains) != 1 {
		return "", fmt.Errorf("expected one chain, got %d", len(chains))
	}

	return fmt.Sprintf(Endpoint, h.api.Endpoints[0].URL), nil
}

// ParseResponse parses the response from the x/prices API and returns the resolved and
// unresolved market map data. The response from the MarketMap API is expected to be a
// a single market map object that was converted from the dYdX market params response.
func (h *APIHandler) ParseResponse(
	chains []types.Chain,
	resp *http.Response,
) types.MarketMapResponse {
	if len(chains) != 1 {
		h.logger.Debug("expected one chain", zap.Any("chains", chains))
		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("expected one chain, got %d", len(chains)),
				providertypes.ErrorInvalidAPIChains,
			),
		)
	}

	if resp == nil {
		h.logger.Debug("got nil response from dydx market params API")
		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("nil response"),
				providertypes.ErrorNoResponse,
			),
		)
	}

	// Parse the response body into a dydx market params response object.
	var params dydxtypes.QueryAllMarketParamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&params); err != nil {
		h.logger.Debug("failed to parse dydx market params response", zap.Error(err))
		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("failed to parse dydx market params response: %w", err),
				providertypes.ErrorFailedToDecode,
			),
		)
	}

	// Convert the dydx market params to a market map.
	marketResp, err := ConvertMarketParamsToMarketMap(params)
	if err != nil {
		h.logger.Debug(
			"failed to convert dydx market params to market map",
			zap.Any("params", params),
			zap.Error(err),
		)

		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("failed to convert dydx market params to market map: %w", err),
				providertypes.ErrorUnknown,
			),
		)
	}

	// validate the market-map
	if err := marketResp.MarketMap.ValidateBasic(); err != nil {
		h.logger.Debug("failed to validate market map", zap.Error(err))
		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("failed to validate market map: %w", err),
				providertypes.ErrorUnknown,
			),
		)
	}

	resolved := make(types.ResolvedMarketMap)
	resolved[chains[0]] = types.NewMarketMapResult(&marketResp, time.Now())

	h.logger.Debug("successfully resolved market map", zap.Int("markets", len(marketResp.MarketMap.Markets)))
	return types.NewMarketMapResponse(resolved, nil)
}
