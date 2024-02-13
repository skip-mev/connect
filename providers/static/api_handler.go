package static

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/base/api/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var _ handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int] = (*MockAPIHandler)(nil)

const (
	// Name is the name of the provider.
	Name = "static-mock-provider"
)

// MockAPIHandler implements a mock API handler that returns static data.
type MockAPIHandler struct {
	exchangeRates map[slinkytypes.CurrencyPair]*big.Int
	currencyPairs []slinkytypes.CurrencyPair
}

// NewAPIHandler returns a new MockAPIHandler. This constructs a
// new static mock provider from the config. Notice this method expects the
// market configuration map to be populated w/ entries of the form CurrencyPair.ToString():
// big.NewInt(price).
func NewAPIHandler(
	cfg config.ProviderConfig,
) (*MockAPIHandler, error) {
	if cfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name to be static-mock-provider, got %s", cfg.Name)
	}

	s := MockAPIHandler{
		exchangeRates: make(map[slinkytypes.CurrencyPair]*big.Int),
		currencyPairs: make([]slinkytypes.CurrencyPair, 0),
	}

	for cpString, market := range cfg.Market.CurrencyPairToMarketConfigs {
		cp, err := slinkytypes.CurrencyPairFromString(cpString)
		if err != nil {
			continue
		}

		price, converted := big.NewInt(0).SetString(market.Ticker, 10)
		if !converted {
			return nil, fmt.Errorf("failed to parse price %s for currency pair %s", price, cpString)
		}

		s.exchangeRates[cp] = price
		s.currencyPairs = append(s.currencyPairs, cp)
	}

	return &s, nil
}

// CreateURL is a no-op.
func (s *MockAPIHandler) CreateURL(_ []slinkytypes.CurrencyPair) (string, error) {
	return "static-url", nil
}

// ParseResponse is a no-op. This simply returns the price of the currency pairs configured
// timestamped with the current time.
func (s *MockAPIHandler) ParseResponse(
	cps []slinkytypes.CurrencyPair,
	_ *http.Response,
) providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int] {
	var (
		resolved   = make(map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int])
		unresolved = make(map[slinkytypes.CurrencyPair]error)
	)

	for _, cp := range cps {
		if price, ok := s.exchangeRates[cp]; ok {
			resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now())
		} else {
			unresolved[cp] = fmt.Errorf("failed to resolve currency pair %s", cp)
		}
	}

	return providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](resolved, unresolved)
}
