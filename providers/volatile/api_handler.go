package volatile

import (
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"

	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// Name is the name of the provider.
	Name = "volatile-exchange-provider"
	Type = types.ConfigType
	// Offset is the average of the returned value from GetVolatilePrice.
	Offset = float64(100)
	// Amplitude is the magnitude of price variation of the cosine function used in GetVolatilePrice.
	Amplitude = float64(0.95)
	// Frequency sizes the repetition of the price curve in GetVolatilePrice.
	Frequency = float64(1)
)

// DefaultAPIConfig is required by the oracle to run. Most fields are ignored.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Enabled:          true,
	MaxQueries:       1,
	Timeout:          500 * time.Millisecond,
	Interval:         500 * time.Millisecond,
	ReconnectTimeout: 500 * time.Millisecond,
	URL:              Name,
}

var DefaultProviderConfig = config.ProviderConfig{
	Name: Name,
	API:  DefaultAPIConfig,
	Type: Type,
}

// APIHandler implements the APIHandler interface.
type APIHandler struct{}

// NewAPIHandler is invoked by the API factory to create the volatile api handler.
func NewAPIHandler() types.PriceAPIDataHandler {
	return &APIHandler{}
}

// CreateURL is a no-op in volatile api handler.
func (v *APIHandler) CreateURL(_ []types.ProviderTicker) (string, error) {
	return "volatile-exchange-url", nil
}

// ParseResponse returns the same volatile price for each ticker.
func (v *APIHandler) ParseResponse(tickers []types.ProviderTicker, _ *http.Response) types.PriceResponse {
	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	volatilePrice := GetVolatilePrice(time.Now, Amplitude, Offset, Frequency)
	for _, ticker := range tickers {
		resolved[ticker] = types.NewPriceResult(volatilePrice, time.Now().UTC())
	}

	return types.NewPriceResponse(resolved, unresolved)
}
