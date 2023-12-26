package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers"
	"github.com/skip-mev/slinky/providers/base"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of the provider.
	Name = "coingecko"
)

var _ base.APIDataHandler[oracletypes.CurrencyPair, *big.Int] = (*CoinGeckoAPIHandler)(nil)

// CoinGeckoAPIHandler implements the Base Provider API handler interface for CoinGecko.
// This provider is a very simple implementation that fetches spot prices from the
// CoinGecko API.
type CoinGeckoAPIHandler struct { //nolint
	logger *zap.Logger

	// pairs is a list of currency pairs that the provider should fetch
	// prices for.
	pairs []oracletypes.CurrencyPair

	// bases is a list of base currencies that the provider should fetch
	// prices for.
	bases string

	// quotes is a list of quote currencies that the provider should fetch
	// prices for.
	quotes string

	// config is the CoinGecko config.
	config Config
}

// NewCoinGeckoAPIHandler returns a new CoinGecko API handler.
func NewCoinGeckoAPIHandler(
	logger *zap.Logger,
	pairs []oracletypes.CurrencyPair,
	providerCfg config.ProviderConfig,
) (*CoinGeckoAPIHandler, error) {
	if providerCfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, providerCfg.Name)
	}

	cfg, err := ReadCoinGeckoConfigFromFile(providerCfg.Path)
	if err != nil {
		return nil, err
	}

	bases, quotes := getUniqueBaseAndQuoteDenoms(pairs)
	logger = logger.With(zap.String("api_handler", Name))
	logger.Info("done initializing api handler")

	return &CoinGeckoAPIHandler{
		bases:  bases,
		quotes: quotes,
		pairs:  pairs,
		logger: logger,
		config: cfg,
	}, nil
}

// Get fetches the latest prices from the CoinGecko API. The price is fetched
// from the CoinGecko API in a single request for all pairs. Since the CoinGecko
// response will match some base denoms to quote denoms that should not be supported,
// we filter out pairs that are not supported by the API handler.
//
// Response format:
//
//	{
//	  "cosmos": {
//	    "usd": 11.35
//	  },
//	  "bitcoin": {
//	    "usd": 10000
//	  }
//	}
func (h *CoinGeckoAPIHandler) Get(ctx context.Context) (map[oracletypes.CurrencyPair]*big.Int, error) {
	url := h.getPriceEndpoint(h.bases, h.quotes)

	// make the request to url and unmarshal the response into respMap
	respMap := make(map[string]map[string]float64)

	// if an API key is set, add it to the request
	var reqFn providers.ReqFn
	if h.config.APIKey != "" {
		reqFn = func(req *http.Request) {
			req.Header.Set(apiKeyHeader, h.config.APIKey)
		}
	}

	if err := providers.GetWithContextAndHeader(ctx, url, func(body []byte) error {
		return json.Unmarshal(body, &respMap)
	}, reqFn); err != nil {
		h.logger.Error(
			"failed to get prices for pairs",
			zap.Error(err),
		)

		return nil, err
	}

	prices := make(map[oracletypes.CurrencyPair]*big.Int)

	// Filter out pairs that are not supported by the API handler.
	for _, pair := range h.pairs {
		base := strings.ToLower(pair.Base)
		quote := strings.ToLower(pair.Quote)

		if _, ok := respMap[base]; !ok {
			continue
		}

		if _, ok := respMap[base][quote]; !ok {
			continue
		}

		price := providers.Float64ToBigInt(respMap[base][quote], pair.Decimals())
		prices[pair] = price

		h.logger.Info(
			"got price for pair",
			zap.String("pair", pair.ToString()),
			zap.String("price", price.String()),
		)
	}

	return prices, nil
}
