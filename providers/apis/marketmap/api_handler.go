package marketmap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	providertypes "github.com/skip-mev/slinky/providers/types"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ types.MarketMapAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the MarketMapAPIDataHandler interface for MarketMap, which can be used
// by a base provider. This is specifically for fetching market map data from the x/marketmap module.
type APIHandler struct {
	// api is the config for the MarketMap API.
	api config.APIConfig
}

// NewAPIHandler returns a new MarketMap MarketMapAPIDataHandler.
func NewAPIHandler(
	api config.APIConfig,
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
		api: api,
	}, nil
}

// CreateURL returns the URL that is used to fetch the latest market map data from the
// MarketMap API. Effectively, this will likely be querying the x/marketmap module.
func (h *APIHandler) CreateURL(chains []types.Chain) (string, error) {
	if len(chains) != 1 {
		return "", fmt.Errorf("expected one chain, got %d", len(chains))
	}

	return h.api.URL, nil
}

// ParseResponse parses the response from the MarketMap API and returns the resolved and
// unresolved market map data. The response from the MarketMap API is expected to be a
// a single market map object.
func (h *APIHandler) ParseResponse(
	chains []types.Chain,
	resp *http.Response,
) types.MarketMapResponse {
	if len(chains) != 1 {
		return types.NewMarketMapResponseWithErr(chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("expected one chain, got %d", len(chains)),
				providertypes.ErrorInvalidAPIChains,
			),
		)
	}

	if resp == nil {
		return types.NewMarketMapResponseWithErr(chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("nil response"),
				providertypes.ErrorNoResponse,
			),
		)
	}

	// Parse the response body into a market map object.
	var market mmtypes.GetMarketMapResponse
	if err := json.NewDecoder(resp.Body).Decode(&market); err != nil {
		return types.NewMarketMapResponseWithErr(chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("failed to parse market map response: %w", err),
				providertypes.ErrorFailedToDecode,
			),
		)
	}

	// Validate the market map response.
	if err := market.MarketMap.ValidateBasic(); err != nil {
		return types.NewMarketMapResponseWithErr(chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("invalid market map response: %w", err),
				providertypes.ErrorInvalidResponse,
			),
		)
	}

	// Ensure the chain id in the response matches the chain id in the request.
	chain := chains[0]
	if market.ChainId != chain.ChainID {
		return types.NewMarketMapResponseWithErr(chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("expected chain id %s, got %s", chain.ChainID, market.ChainId),
				providertypes.ErrorInvalidChainID,
			),
		)
	}

	resolved := make(types.ResolvedMarketMap)
	resolved[chain] = types.NewMarketMapResult(&market, time.Now())
	return types.NewMarketMapResponse(resolved, nil)
}
