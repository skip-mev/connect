package volatile

import (
	"fmt"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"

	"github.com/skip-mev/slinky/oracle/types"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

const (
	// Name is the name of the provider.
	Name = "volatile-exchange-provider"
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

// APIHandler implements the APIHandler interface.
type APIHandler struct {
	tickers map[mmtypes.Ticker]bool
}

// NewAPIHandler is invoked by the API factory to create the volatile api handler.
func NewAPIHandler(market types.ProviderMarketMap) (types.PriceAPIDataHandler, error) {
	if market.Name != Name {
		return nil, fmt.Errorf("expected market config name to be %s, got %s", Name, market.Name)
	}
	v := APIHandler{
		tickers: make(map[mmtypes.Ticker]bool),
	}

	for ticker := range market.TickerConfigs {
		v.tickers[ticker] = true
	}
	return &v, nil
}

// CreateURL is a no-op in volatile api handler.
func (v *APIHandler) CreateURL(_ []mmtypes.Ticker) (string, error) {
	return "volatile-exchange-url", nil
}

// ParseResponse returns the same volatile price for each ticker.
func (v *APIHandler) ParseResponse(tickers []mmtypes.Ticker, _ *http.Response) types.PriceResponse {
	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	volatilePrice := GetVolatilePrice(time.Now, Amplitude, Offset, Frequency)

	for _, ticker := range tickers {
		if _, ok := v.tickers[ticker]; ok {
			resolved[ticker] = types.NewPriceResult(volatilePrice, time.Now())
		} else {
			err := fmt.Errorf("failed to resolve ticker %s", ticker)
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(err, providertypes.ErrorUnknownPair),
			}
		}
	}

	return types.NewPriceResponse(resolved, unresolved)
}
