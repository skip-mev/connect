package static

import (
	"context"
	"fmt"
	"math/big"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"go.uber.org/zap"
)

var _ base.APIDataHandler[oracletypes.CurrencyPair, *big.Int] = (*StaticMockAPIHandler)(nil)

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
// TokenNameToSymbol map to be populated w/ entries of the form CurrencyPair.ToString():
// big.NewInt(price).
func NewStaticMockAPIHandler(
	_ *zap.Logger,
	_ []oracletypes.CurrencyPair,
	providerCfg config.ProviderConfig,
) (*StaticMockAPIHandler, error) {
	if providerCfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name to be static-mock-provider, got %s", providerCfg)
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

// Get returns the latest exchange rates for the given currency pairs.
func (s *StaticMockAPIHandler) Get(_ context.Context) (map[oracletypes.CurrencyPair]*big.Int, error) {
	return s.exchangeRates, nil
}
