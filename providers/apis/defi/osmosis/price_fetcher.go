package osmosis

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	oracletypes "github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

var _ oracletypes.PriceAPIFetcher = &APIPriceFetcher{}

type APIPriceFetcher struct {
	// config is the APIConfiguration for this provider
	api config.APIConfig

	// client is the osmosis client used to query the API.
	client Client

	// metaDataPerTicker is a map of ticker.String() -> TickerMetadata
	metaDataPerTicker *metadataCache

	// logger
	logger *zap.Logger
}

// NewAPIPriceFetcher returns a new APIPriceFetcher. This method constructs the
// default Osmosis client in accordance with the config's endpoints.
func NewAPIPriceFetcher(
	logger *zap.Logger,
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
) (*APIPriceFetcher, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config: %w", err)
	}

	if api.Name != Name {
		return nil, fmt.Errorf("invalid api name; expected %s, got %s", Name, api.Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api is not enabled")
	}

	if apiMetrics == nil {
		return nil, fmt.Errorf("metrics cannot be nil")
	}

	client, err := NewMultiClientFromEndpoints(logger, api, apiMetrics)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client: %w", err)
	}

	return NewAPIPriceFetcherWithClient(logger, api, apiMetrics, client)
}

// NewAPIPriceFetcherWithClient returns a new APIPriceFetcher. This method constructs the
// osmosis client with the given client.
func NewAPIPriceFetcherWithClient(
	logger *zap.Logger,
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
	client Client,
) (*APIPriceFetcher, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config: %w", err)
	}

	if api.Name != Name {
		return nil, fmt.Errorf("invalid api name; expected %s, got %s", Name, api.Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api is not enabled")
	}

	if apiMetrics == nil {
		return nil, fmt.Errorf("metrics cannot be nil")
	}

	return &APIPriceFetcher{
		api:               api,
		client:            client,
		logger:            logger.With(zap.String("fetcher", Name)),
		metaDataPerTicker: newMetadataCache(),
	}, nil
}

// Fetch fetches prices from the osmosis API for the given currency-pairs. Specifically
// for each currency-pair,
//   - Query the spot price.
func (pf *APIPriceFetcher) Fetch(
	ctx context.Context,
	tickers []oracletypes.ProviderTicker,
) oracletypes.PriceResponse {
	resolved := make(oracletypes.ResolvedPrices)
	unresolved := make(oracletypes.UnResolvedPrices)

	g, ctx := errgroup.WithContext(ctx)
	unresolvedMtx := sync.Mutex{}
	resolveMtx := sync.Mutex{}
	g.SetLimit(pf.api.MaxQueries)

	// setup callbacks for writing to maps in parallel
	unresolvedTickerCallback := func(ticker oracletypes.ProviderTicker, err providertypes.ErrorWithCode) {
		unresolvedMtx.Lock()
		defer unresolvedMtx.Unlock()
		unresolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: err,
		}
	}

	resolvedTickerCallback := func(ticker oracletypes.ProviderTicker, price *big.Float) {
		resolveMtx.Lock()
		defer resolveMtx.Unlock()
		resolved[ticker] = oracletypes.NewPriceResult(price, time.Now().UTC())
	}

	pf.logger.Info("fetching for tickers", zap.Any("tickers", tickers))

	// make sure metadata cache is set properly
	for _, ticker := range tickers {
		_, found := pf.metaDataPerTicker.getMetadataPerTicker(ticker)
		if !found {
			_, err := pf.metaDataPerTicker.updateMetaDataCache(ticker)
			if err != nil {
				pf.logger.Debug("failed to update metadata cache", zap.Error(err))
			}
		}
	}

	for _, ticker := range tickers {
		g.Go(func() error {
			ticker := ticker
			var err error

			metadata, found := pf.metaDataPerTicker.getMetadataPerTicker(ticker)
			if !found {
				unresolvedTickerCallback(ticker, providertypes.NewErrorWithCode(
					NoOsmosisMetadataForTickerError(ticker.String()),
					providertypes.ErrorTickerMetadataNotFound,
				))

				return nil
			}

			callCtx, cancel := context.WithTimeout(ctx, pf.api.Timeout)
			defer cancel()

			resp, err := pf.client.SpotPrice(callCtx,
				metadata.PoolID,
				metadata.BaseTokenDenom,
				metadata.QuoteTokenDenom,
			)
			if err != nil {
				unresolvedTickerCallback(ticker, providertypes.NewErrorWithCode(
					err,
					providertypes.ErrorAPIGeneral,
				))

				return nil
			}

			price, err := calculatePrice(resp)
			if err != nil {
				pf.logger.Error("failed parse spot price response", zap.Error(err))

				unresolvedTickerCallback(ticker, providertypes.NewErrorWithCode(
					err,
					providertypes.ErrorFailedToParsePrice,
				))

				return nil
			}

			resolvedTickerCallback(ticker, price)

			return nil
		})
	}

	_ = g.Wait()

	return oracletypes.NewPriceResponse(resolved, unresolved)
}

func calculatePrice(resp WrappedSpotPriceResponse) (*big.Float, error) {
	return math.Float64StringToBigFloat(resp.SpotPrice)
}
