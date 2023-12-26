package erc4626sharepriceoracle

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/common"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers"
	"github.com/skip-mev/slinky/providers/base"
	"github.com/skip-mev/slinky/providers/evm"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of this provider
	Name = "erc4626-share-price-oracle"
)

var _ base.APIDataHandler[oracletypes.CurrencyPair, *big.Int] = (*ERC4626SharePriceAPIHandler)(nil)

// ERC4626SharePriceAPIHandler implements the Base Provider API handler interface for instances of the
// ERC4626SharePriceOracle.sol contract.
type ERC4626SharePriceAPIHandler struct {
	logger *zap.Logger

	// pairs is a list of currency pairs that the provider should fetch
	// prices for.
	pairs []oracletypes.CurrencyPair

	// config is the ERC4626SharePriceOracle config.
	config evm.Config

	// rpcEndpoint is the endpoint of the ethereum rpc node to use for querying. This
	// is used to make RPC calls to the Ethereum node with a configured API key.
	rpcEndpoint string
}

// NewERC4626SharePriceAPIHandler returns a new ERC4626 Share Price API handler.
func NewERC4626SharePriceAPIHandler(
	logger *zap.Logger,
	pairs []oracletypes.CurrencyPair,
	providerCfg config.ProviderConfig,
) (*ERC4626SharePriceAPIHandler, error) {
	if providerCfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name to be %s, got %s", Name, providerCfg.Name)
	}

	cfg, err := evm.ReadEVMConfigFromFile(providerCfg.Path)
	if err != nil {
		return nil, err
	}

	filteredPairs := make([]oracletypes.CurrencyPair, 0)
	for _, pair := range pairs {
		if metadata, ok := cfg.TokenNameToMetadata[pair.Quote]; ok {
			if !common.IsHexAddress(metadata.Symbol) {
				return nil, fmt.Errorf("invalid contract address: %s", metadata.Symbol)
			}

			filteredPairs = append(filteredPairs, pair)
		}
	}

	logger = logger.With(zap.String("api_handler", Name))
	logger.Info("done initializing api handler")

	return &ERC4626SharePriceAPIHandler{
		logger:      logger,
		rpcEndpoint: getRPCEndpoint(cfg),
		config:      cfg,
		pairs:       filteredPairs,
	}, nil
}

// Get returns the prices of the given pairs.
func (h *ERC4626SharePriceAPIHandler) Get(ctx context.Context) (map[oracletypes.CurrencyPair]*big.Int, error) {
	type priceData struct {
		price *big.Int
		cp    oracletypes.CurrencyPair
	}

	// create response channel
	responses := make(chan priceData, len(h.pairs))

	wg := sync.WaitGroup{}
	wg.Add(len(h.pairs))

	// fan-out requests to RPC provider
	for _, currencyPair := range h.pairs {
		go func(pair oracletypes.CurrencyPair) {
			defer wg.Done()

			// get price
			qp, err := h.getPriceForPair(pair)
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
