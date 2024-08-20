package static

import (
	"math/big"
	"net/http"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"github.com/skip-mev/connect/v2/oracle/types"
)

var _ types.PriceAPIDataHandler = (*MockAPIHandler)(nil)

const (
	// Name is the name of the provider.
	Name = "static-mock-provider"
)

// MockAPIHandler implements a mock API handler that returns static data.
type MockAPIHandler struct{}

// NewAPIHandler returns a new MockAPIHandler. This constructs a new static mock provider from
// the config. Notice this method expects the market configuration map to the offchain ticker
// to the desired price.
func NewAPIHandler() types.PriceAPIDataHandler {
	return &MockAPIHandler{}
}

// CreateURL is a no-op.
func (s *MockAPIHandler) CreateURL(_ []types.ProviderTicker) (string, error) {
	return "static-url", nil
}

// ParseResponse is a no-op. This simply returns the price of the tickers configured,
// timestamped with the current time.
func (s *MockAPIHandler) ParseResponse(
	tickers []types.ProviderTicker,
	_ *http.Response,
) types.PriceResponse {
	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	for _, ticker := range tickers {
		var metaData MetaData
		if err := metaData.FromJSON(ticker.GetJSON()); err == nil {
			resolved[ticker] = types.NewPriceResult(
				big.NewFloat(metaData.Price),
				time.Now().UTC(),
			)
		} else {
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(
					err,
					providertypes.ErrorFailedToParsePrice,
				),
			}
		}
	}

	return types.NewPriceResponse(resolved, unresolved)
}
