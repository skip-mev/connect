package coinmarketcap

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers"
	"github.com/skip-mev/slinky/providers/base"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of the provider.
	Name = "coinmarketcap"
)

var _ base.APIDataHandler[oracletypes.CurrencyPair, *big.Int] = (*CoinMarketCapAPIHandler)(nil)

// CoinMarketCapAPIHandler implements the Base Provider API handler interface for CoinMarketCap.
// This provider is a very simple implementation that fetches spot prices from the CoinMarketCap API.
type CoinMarketCapAPIHandler struct { //nolint
	logger *zap.Logger

	// pairs is a list of currency pairs that the provider should fetch
	// prices for.
	pairs []oracletypes.CurrencyPair

	// config is the coinmarketcap config.
	config Config
}

// NewCoinMarketCapAPIHandler returns a new CoinMarketCap API handler.
func NewCoinMarketCapAPIHandler(
	logger *zap.Logger,
	pairs []oracletypes.CurrencyPair,
	providerCfg config.ProviderConfig,
) (*CoinMarketCapAPIHandler, error) {
	if providerCfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, providerCfg.Name)
	}

	cfg, err := ReadCoinMarketCapConfigFromFile(providerCfg.Path)
	if err != nil {
		return nil, err
	}

	logger = logger.With(zap.String("api_handler", Name))
	logger.Info("done initializing api handler")

	return &CoinMarketCapAPIHandler{
		pairs:  pairs,
		config: cfg,
		logger: logger,
	}, nil
}

// Get fetches the latest prices from the CoinMarketCap API. This method starts all price
// requests concurrently, and waits for them all to finish, or for the context to be
// cancelled, at which point it aggregates the responses and returns.
func (h *CoinMarketCapAPIHandler) Get(ctx context.Context) (map[oracletypes.CurrencyPair]*big.Int, error) {
	type priceData struct {
		price *big.Int
		cp    oracletypes.CurrencyPair
	}

	// create response channel
	responses := make(chan priceData, len(h.pairs))

	wg := sync.WaitGroup{}
	wg.Add(len(h.pairs))

	// fan-out requests to coinmarketcap api
	for _, currencyPair := range h.pairs {
		go func(pair oracletypes.CurrencyPair) {
			defer wg.Done()

			// get price
			qp, err := h.getPriceForPair(ctx, pair)
			if err != nil {
				h.logger.Error(
					"failed to get price for pair",
					zap.String("pair", pair.ToString()),
					zap.Error(err),
				)
			} else {
				h.logger.Info(
					"got price for pair",
					zap.String("pair", pair.ToString()),
					zap.String("price", qp.String()),
				)

				// send price to response channel
				responses <- priceData{
					qp,
					pair,
				}
			}
		}(currencyPair)
	}

	// close response channel when all requests have been processed, or if context is cancelled
	go func() {
		defer close(responses)

		select {
		case <-ctx.Done():
			return
		case <-providers.Finish(&wg):
			return
		}
	}()

	// fan-in
	prices := make(map[oracletypes.CurrencyPair]*big.Int)
	for resp := range responses {
		prices[resp.cp] = resp.price
	}

	return prices, nil
}
