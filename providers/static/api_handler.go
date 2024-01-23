package static

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/api/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var _ handlers.APIDataHandler[oracletypes.CurrencyPair, *big.Int] = (*StaticMockAPIHandler)(nil)

const (
	// Name is the name of the provider.
	Name = "static-mock-provider"
)

// StaticMockAPIHandler implements a mock API handler that returns static data.
type StaticMockAPIHandler struct { //nolint
	exchangeRates map[oracletypes.CurrencyPair]*big.Int
	currencyPairs []oracletypes.CurrencyPair
}

// NewStaticMockAPIHandler returns a new StaticMockAPIHandler. This constructs a
// new static mock provider from the config. Notice this method expects the
// TokenNameToSymbol map to be populated w/ entries of the form CurrencyPair.String():
// big.NewInt(price).
func NewStaticMockAPIHandler(
	providerCfg config.ProviderConfig,
) (*StaticMockAPIHandler, error) {
	if providerCfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name to be static-mock-provider, got %s", providerCfg.Name)
	}

	config, err := ReadStaticMockProviderConfigFromFile(providerCfg.Path)
	if err != nil {
		return nil, err
	}

	s := StaticMockAPIHandler{
		exchangeRates: make(map[oracletypes.CurrencyPair]*big.Int),
		currencyPairs: make([]oracletypes.CurrencyPair, 0),
	}

	for cpString, price := range config.TokenPrices {
		cp, err := oracletypes.CurrencyPairFromString(cpString)
		if err != nil {
			continue
		}

		price, converted := big.NewInt(0).SetString(price, 10)
		if !converted {
			return nil, fmt.Errorf("failed to parse price %s for currency pair %s", price, cpString)
		}

		s.exchangeRates[cp] = price
		s.currencyPairs = append(s.currencyPairs, cp)
	}

	return &s, nil
}

// CreateURL is a no-op.
func (s *StaticMockAPIHandler) CreateURL(_ []oracletypes.CurrencyPair) (string, error) {
	return "static-url", nil
}

// Atomic returns true as the static mock provider is atomic i.e. returns the price of all
// currency pairs in a single request.
func (s *StaticMockAPIHandler) Atomic() bool {
	return true
}

// ParseResponse is a no-op. This simply returns the price of the currency pairs configured
// timestamped with the current time.
func (s *StaticMockAPIHandler) ParseResponse(
	cps []oracletypes.CurrencyPair,
	_ *http.Response,
) providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int] {
	var (
		resolved   = make(map[oracletypes.CurrencyPair]providertypes.Result[*big.Int])
		unresolved = make(map[oracletypes.CurrencyPair]error)
	)

	for _, cp := range cps {
		if price, ok := s.exchangeRates[cp]; ok {
			resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now())
		} else {
			unresolved[cp] = fmt.Errorf("failed to resolve currency pair %s", cp)
		}
	}

	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved)
}

// Name returns the name of the provider.
func (s *StaticMockAPIHandler) Name() string {
	return Name
}
