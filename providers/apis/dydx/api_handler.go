package dydx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	dydxtypes "github.com/skip-mev/slinky/providers/apis/dydx/types"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/service/clients/marketmap/types"
)

var _ types.MarketMapAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the MarketMapAPIDataHandler interface for the dYdX prices module, which can be used
// by a base provider. This is specifically for fetching market data from the dYdX prices module, which is
// then translated to a market map.
type APIHandler struct {
	// api is the api config for the dYdX market params API.
	api config.APIConfig

	// logger is the logger for the API handler.
	logger *zap.Logger
}

// NewAPIHandler returns a new MarketMap MarketMapAPIDataHandler.
func NewAPIHandler(
	api config.APIConfig,
	logger *zap.Logger,
) (types.MarketMapAPIDataHandler, error) {
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

	return fmt.Sprintf(Endpoint, h.api.URL), nil
}

// ParseResponse parses the response from the x/prices API and returns the resolved and
// unresolved market map data. The response from the MarketMap API is expected to be a
// a single market map object that was converted from the dYdX market params response.
func (h *APIHandler) ParseResponse(
	chains []types.Chain,
	resp *http.Response,
) types.MarketMapResponse {
	if len(chains) != 1 {
		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("expected one chain, got %d", len(chains)),
				providertypes.ErrorInvalidAPIChains,
			),
		)
	}

	if resp == nil {
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
		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("failed to parse dydx market params response: %w", err),
				providertypes.ErrorFailedToDecode,
			),
		)
	}

	// Convert the dydx market params to a market map.
	marketResp, err := ConvertMarketParamsToMarketMap(params, h.logger)
	if err != nil {
		return types.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("failed to convert dydx market params to market map: %w", err),
				providertypes.ErrorUnknown,
			),
		)
	}

	resolved := make(types.ResolvedMarketMap)
	resolved[chains[0]] = types.NewMarketMapResult(&marketResp, time.Now())
	return types.NewMarketMapResponse(resolved, nil)
}
