package coinbase

import (
	"context"
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of the provider.
	Name = "coinbase"
)

var _ base.APIDataHandler[oracletypes.CurrencyPair, *big.Int] = (*CoinBaseAPIHandler)(nil)

// CoinBaseAPIHandler implements the Base Provider API handler interface for Coinbase.
// This provider is a very simple implementation that fetches spot prices from the Coinbase API.
type CoinBaseAPIHandler struct { //nolint
	logger *zap.Logger

	// pairs is a list of currency pairs that the provider should fetch
	// prices for.
	pairs []oracletypes.CurrencyPair

	// config is the Coinbase config.
	config Config
}

// NewCoinBaseAPIHandler returns a new Coinbase API handler.
func NewCoinBaseAPIHandler(
	logger *zap.Logger,
	pairs []oracletypes.CurrencyPair,
	providerCfg config.ProviderConfig,
) (*CoinBaseAPIHandler, error) {
	if providerCfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, providerCfg.Name)
	}

	cfg, err := ReadCoinbaseConfigFromFile(providerCfg.Path)
	if err != nil {
		return nil, err
	}

	logger = logger.With(zap.String("api_handler", Name))
	logger.Info("done initializing api handler")

	return &CoinBaseAPIHandler{
		logger: logger,
		pairs:  pairs,
		config: cfg,
	}, nil
}

// Get fetches the latest prices from the Coinbase API.
func (h *CoinBaseAPIHandler) Get(ctx context.Context) (map[oracletypes.CurrencyPair]*big.Int, error) {
	resp := make(map[oracletypes.CurrencyPair]*big.Int)

	for _, currencyPair := range h.pairs {
		spotPrice, err := h.getPriceForPair(ctx, currencyPair)
		if err != nil {
			h.logger.Error(
				Name,
				zap.String("failed to get price for pair", currencyPair.ToString()),
				zap.Error(err),
			)
			continue
		}

		h.logger.Info(
			"got price for pair",
			zap.String("pair", currencyPair.ToString()),
			zap.String("price", spotPrice.String()),
		)

		resp[currencyPair] = spotPrice
	}

	return resp, nil
}
