package osmosis

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
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

	wg := sync.WaitGroup{}
	unresolvedMtx := sync.Mutex{}
	resolveMtx := sync.Mutex{}
	wg.Add(len(tickers))

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
		go func(ticker oracletypes.ProviderTicker) {
			var err error
			defer wg.Done()

			metadata, found := pf.metaDataPerTicker.getMetadataPerTicker(ticker)
			if !found {
				unresolvedMtx.Lock()
				defer unresolvedMtx.Unlock()
				unresolved[ticker] = providertypes.UnresolvedResult{
					ErrorWithCode: providertypes.NewErrorWithCode(
						NoOsmosisMetadataForTickerError(ticker.String()),
						providertypes.ErrorTickerMetadataNotFound,
					),
				}
				return
			}

			callCtx, cancel := context.WithTimeout(ctx, pf.api.Timeout)
			defer cancel()

			resp, err := pf.client.SpotPrice(callCtx,
				metadata.PoolID,
				metadata.BaseTokenDenom,
				metadata.QuoteTokenDenom,
			)
			if err != nil {
				unresolvedMtx.Lock()
				defer unresolvedMtx.Unlock()
				pf.logger.Error("failed to fetch spot price", zap.Error(err))
				unresolved[ticker] = providertypes.UnresolvedResult{
					ErrorWithCode: providertypes.NewErrorWithCode(
						err,
						providertypes.ErrorAPIGeneral,
					),
				}
				return
			}

			price, err := calculatePrice(resp)
			if err != nil {
				pf.logger.Error("failed parse spot price response", zap.Error(err))
				unresolved[ticker] = providertypes.UnresolvedResult{
					ErrorWithCode: providertypes.NewErrorWithCode(
						err,
						providertypes.ErrorFailedToParsePrice,
					),
				}
				return
			}

			resolveMtx.Lock()
			defer resolveMtx.Unlock()
			resolved[ticker] = oracletypes.NewPriceResult(price, time.Now().UTC())
		}(ticker)
	}

	wg.Wait()

	return oracletypes.NewPriceResponse(resolved, unresolved)
}

func calculatePrice(resp SpotPriceResponse) (*big.Float, error) {
	return math.Float64StringToBigFloat(resp.SpotPrice)
}
