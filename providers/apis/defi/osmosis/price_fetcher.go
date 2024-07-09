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
	"github.com/skip-mev/slinky/providers/apis/defi/osmosis/queryproto"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

const (
	// responseDecimals is the set of decimals to be used for scaling.
	// https://github.com/osmosis-labs/osmosis/blob/194ef2da5f0dcb9401e4a9bbbeaeee30aefcca67/x/gamm/keeper/grpc_query.go#L387.
	responseDecimals = 18
)

var _ oracletypes.PriceAPIFetcher = &APIPriceFetcher{}

type APIPriceFetcher struct {
	// config is the APIConfiguration for this provider
	api config.APIConfig

	// client is the osmosis gRPC client used to query the API.
	client GRPCCLient

	// metaDataPerTicker is a map of ticker.String() -> TickerMetadata
	metaDataPerTicker *metadataCache

	// logger
	logger *zap.Logger
}

// NewAPIPriceFetcher returns a new APIPriceFetcher. This method constructs the
// default Osmosis gRPC client in accordance with the config's endpoints.
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

	client, err := NewGRPCMultiClient(logger, api, apiMetrics)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client: %w", err)
	}

	return NewAPIPriceFetcherWithClient(logger, api, apiMetrics, client)
}

// NewAPIPriceFetcherWithClient returns a new APIPriceFetcher. This method constructs the
// osmosis gRPC client with the given client.
func NewAPIPriceFetcherWithClient(
	logger *zap.Logger,
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
	client GRPCCLient,
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
		api:    api,
		client: client,
		logger: logger.With(zap.String("fetcher", Name)),
	}, nil
}

// Fetch fetches prices from the osmosis gRPC API for the given currency-pairs. Specifically
// for each currency-pair,
//   - Query the spot price.
func (pf *APIPriceFetcher) Fetch(
	ctx context.Context,
	tickers []oracletypes.ProviderTicker,
) oracletypes.PriceResponse {
	resolved := make(oracletypes.ResolvedPrices)
	unresolved := make(oracletypes.UnResolvedPrices)

	wg := sync.WaitGroup{}
	wg.Add(len(tickers))

	for _, ticker := range tickers {
		go func(ticker oracletypes.ProviderTicker) {
			var err error
			defer wg.Done()

			// get or set metadata
			// TODO: how expensive is this ?
			metadata, found := pf.metaDataPerTicker.getMetadataPerTicker(ticker)
			if !found {
				metadata, err = pf.metaDataPerTicker.updateMetaDataCache(ticker)
				if err != nil {
					unresolved[ticker] = providertypes.UnresolvedResult{
						ErrorWithCode: providertypes.NewErrorWithCode(
							NoOsmosisMetadataForTickerError(ticker.String()),
							providertypes.ErrorTickerMetadataNotFound,
						),
					}
					return
				}
			}

			callCtx, cancel := context.WithTimeout(ctx, pf.api.Timeout)
			defer cancel()

			resp, err := pf.client.SpotPrice(callCtx, &queryproto.SpotPriceRequest{
				PoolId:          metadata.PoolID,
				BaseAssetDenom:  metadata.BaseTokenDenom,
				QuoteAssetDenom: metadata.QuoteTokenDenom,
			})
			if err != nil {
				pf.logger.Error("failed to fetch spot price", zap.Error(err))
				unresolved[ticker] = providertypes.UnresolvedResult{
					ErrorWithCode: providertypes.NewErrorWithCode(
						GRPCError(err),
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

			resolved[ticker] = oracletypes.NewPriceResult(price, time.Now().UTC())
		}(ticker)
	}

	wg.Wait()

	return oracletypes.NewPriceResponse(resolved, unresolved)
}

func calculatePrice(resp *queryproto.SpotPriceResponse) (*big.Float, error) {
	return math.Float64StringToBigFloat(resp.SpotPrice)
}
