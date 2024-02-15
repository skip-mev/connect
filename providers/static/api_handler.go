package static

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/types"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ types.PriceAPIDataHandler = (*MockAPIHandler)(nil)

const (
	// Name is the name of the provider.
	Name = "static-mock-provider"
)

// MockAPIHandler implements a mock API handler that returns static data.
type MockAPIHandler struct {
	exchangeRates types.TickerPrices
	tickers       []mmtypes.Ticker
}

// NewAPIHandler returns a new MockAPIHandler. This constructs a new static mock provider from
// the config. Notice this method expects the market configuration map to the offchain ticker
// to the desired price.
func NewAPIHandler(
	market mmtypes.MarketConfig,
) (types.PriceAPIDataHandler, error) {
	if market.Name != Name {
		return nil, fmt.Errorf("expected market config name to be %s, got %s", Name, market.Name)
	}

	s := MockAPIHandler{
		exchangeRates: make(types.TickerPrices),
		tickers:       make([]mmtypes.Ticker, 0),
	}

	for cpString, market := range market.TickerConfigs {
		price, converted := big.NewInt(0).SetString(market.OffChainTicker, 10)
		if !converted {
			return nil, fmt.Errorf("failed to parse price %s for ticker %s", price, cpString)
		}

		s.exchangeRates[market.Ticker] = price
		s.tickers = append(s.tickers, market.Ticker)
	}

	return &s, nil
}

// CreateURL is a no-op.
func (s *MockAPIHandler) CreateURL(_ []mmtypes.Ticker) (string, error) {
	return "static-url", nil
}

// ParseResponse is a no-op. This simply returns the price of the tickers configured,
// timestamped with the current time.
func (s *MockAPIHandler) ParseResponse(
	tickers []mmtypes.Ticker,
	_ *http.Response,
) types.PriceResponse {
	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	for _, ticker := range tickers {
		if price, ok := s.exchangeRates[ticker]; ok {
			resolved[ticker] = providertypes.NewResult[*big.Int](price, time.Now())
		} else {
			unresolved[ticker] = fmt.Errorf("failed to resolve ticker %s", ticker)
		}
	}

	return providertypes.NewGetResponse(resolved, unresolved)
}
