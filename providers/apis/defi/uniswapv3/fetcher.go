package uniswapv3

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/apis/defi/ethmulticlient"
	uniswappool "github.com/skip-mev/connect/v2/providers/apis/defi/uniswapv3/pool"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

var _ types.PriceAPIFetcher = (*PriceFetcher)(nil)

// PriceFetcher is the Uniswap V3 price fetcher. This fetcher is responsible for
// querying Uniswap V3 pool contracts and returning the price of a given ticker. The price is
// derived from the slot 0 data of the pool contract.
//
// To read more about how the price is calculated, see the Uniswap V3 documentation
// https://blog.uniswap.org/uniswap-v3-math-primer.
//
// We utilize the eth client's BatchCallContext to batch the calls to the ethereum network as
// this is more performant than making individual calls or the multi call contract:
// https://docs.chainstack.com/docs/http-batch-request-vs-multicall-contract#performance-comparison.
type PriceFetcher struct {
	logger *zap.Logger
	api    config.APIConfig

	// client is the EVM client implementation. This is used to interact with the ethereum network.
	client ethmulticlient.EVMClient
	// abi is the uniswap v3 pool abi. This is used to pack the slot0 call to the pool contract
	// and parse the result.
	abi *abi.ABI
	// payload is the packed slot0 call to the pool contract. Since the slot0 payload is the same
	// for all pools, we can reuse this payload for all pools.
	payload []byte
	// poolCache is a cache of the tickers to pool configs. This is used to avoid unmarshalling
	// the metadata for each ticker.
	poolCache map[types.ProviderTicker]PoolConfig
}

// NewPriceFetcher returns a new Uniswap V3 price fetcher.
func NewPriceFetcher(
	ctx context.Context,
	logger *zap.Logger,
	apiMetrics metrics.APIMetrics,
	api config.APIConfig,
) (*PriceFetcher, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if apiMetrics == nil {
		return nil, fmt.Errorf("api metrics is nil")
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config: %w", err)
	}

	if !IsValidProviderName(api.Name) {
		return nil, fmt.Errorf("invalid api config name %s", api.Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api config for %s is not enabled", api.Name)
	}

	var (
		client ethmulticlient.EVMClient
		err    error
	)
	switch {
	case len(api.Endpoints) > 1:
		client, err = ethmulticlient.NewMultiRPCClientFromEndpoints(
			ctx,
			logger,
			api,
			apiMetrics,
		)
	case len(api.Endpoints) == 1:
		client, err = ethmulticlient.NewGoEthereumClientImpl(
			ctx,
			apiMetrics,
			api,
			0,
		)
	default:
		err = fmt.Errorf("no endpoints were provided")
	}
	if err != nil {
		return nil, err
	}

	return NewPriceFetcherWithClient(
		logger,
		api,
		client,
	)
}

// NewPriceFetcherWithClient returns a new PriceFetcher.
// It requires a pre-validated config, and initialized client.
func NewPriceFetcherWithClient(
	logger *zap.Logger,
	api config.APIConfig,
	client ethmulticlient.EVMClient,
) (*PriceFetcher, error) {
	abi, err := uniswappool.UniswapMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get uniswap abi: %w", err)
	}

	payload, err := abi.Pack(ContractMethod)
	if err != nil {
		return nil, fmt.Errorf("failed to pack slot0: %w", err)
	}

	return &PriceFetcher{
		logger:    logger.With(zap.String("fetcher", api.Name)),
		api:       api,
		client:    client,
		abi:       abi,
		payload:   payload,
		poolCache: make(map[types.ProviderTicker]PoolConfig),
	}, nil
}

// Fetch returns the price of a given set of tickers. This fetch utilizes the batch call to lower
// overhead of making individual RPC calls for each ticker. The fetcher will query the Uniswap V3
// pool contract for the price of the pool. The price is derived from the slot 0 data of the pool
// contract, specifically the sqrtPriceX96 value.
func (u *PriceFetcher) Fetch(
	ctx context.Context,
	tickers []types.ProviderTicker,
) types.PriceResponse {
	var (
		resolved   = make(types.ResolvedPrices)
		unResolved = make(types.UnResolvedPrices)
	)

	// Create a batch element for each ticker and pool.
	batchElems := make([]rpc.BatchElem, len(tickers))
	pools := make([]PoolConfig, len(tickers))
	for i, ticker := range tickers {
		pool, err := u.GetPool(ticker)
		if err != nil {
			u.logger.Debug(
				"failed to get pool for ticker",
				zap.String("ticker", ticker.String()),
				zap.Error(err),
			)

			return types.NewPriceResponseWithErr(
				tickers,
				providertypes.NewErrorWithCode(
					fmt.Errorf("failed to get pool: %w", err),
					providertypes.ErrorFailedToDecode,
				),
			)
		}

		// Create a batch element for the ticker and pool.
		var result string
		batchElems[i] = rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				map[string]interface{}{
					"to":   common.HexToAddress(pool.Address),
					"data": hexutil.Bytes(u.payload), // slot0 call to the pool contract.
				},
				"latest", // latest signifies the latest block.
			},
			Result: &result,
		}
		pools[i] = pool
	}

	// Batch call to the EVM.
	if err := u.client.BatchCallContext(ctx, batchElems); err != nil {
		u.logger.Debug(
			"failed to batch call to ethereum network for all tickers",
			zap.Error(err),
		)

		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorAPIGeneral),
		)
	}

	// Parse the result from the batch call for each ticker.
	for i, ticker := range tickers {
		result := batchElems[i]
		if result.Error != nil {
			u.logger.Debug(
				"failed to batch call to ethereum network for ticker",
				zap.String("ticker", ticker.String()),
				zap.Error(result.Error),
			)

			unResolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(
					result.Error,
					providertypes.ErrorUnknown,
				),
			}

			continue
		}

		// Parse the sqrtPriceX96 from the result.
		sqrtPriceX96, err := u.ParseSqrtPriceX96(result.Result)
		if err != nil {
			u.logger.Debug(
				"failed to parse sqrt price x96",
				zap.String("ticker", ticker.String()),
				zap.Error(err),
			)

			unResolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(
					err,
					providertypes.ErrorFailedToParsePrice,
				),
			}

			continue
		}

		// Convert the sqrtPriceX96 to a price. This is the raw, unscaled price.
		price := ConvertSquareRootX96Price(sqrtPriceX96)

		// Scale the price to the respective token decimals.
		scaledPrice := ScalePrice(pools[i], price)
		resolved[ticker] = types.NewPriceResult(scaledPrice, time.Now().UTC())
	}

	// Add the price to the resolved prices.
	return types.NewPriceResponse(resolved, unResolved)
}

// GetPool returns the uniswap pool for the given ticker. This will unmarshal the metadata
// and validate the pool config which contains all required information to query the EVM.
func (u *PriceFetcher) GetPool(
	ticker types.ProviderTicker,
) (PoolConfig, error) {
	if pool, ok := u.poolCache[ticker]; ok {
		return pool, nil
	}

	var cfg PoolConfig
	if err := json.Unmarshal([]byte(ticker.GetJSON()), &cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal pool config on ticker: %w", err)
	}
	if err := cfg.ValidateBasic(); err != nil {
		return cfg, fmt.Errorf("invalid ticker pool config: %w", err)
	}

	u.poolCache[ticker] = cfg
	return cfg, nil
}

// ParseSqrtPriceX96 parses the sqrtPriceX96 from the result of the batch call.
func (u *PriceFetcher) ParseSqrtPriceX96(
	result interface{},
) (*big.Int, error) {
	r, ok := result.(*string)
	if !ok {
		return nil, fmt.Errorf("expected result to be a string, got %T", result)
	}

	if r == nil {
		return nil, fmt.Errorf("result is nil")
	}

	bz, err := hexutil.Decode(*r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex result: %w", err)
	}

	out, err := u.abi.Methods[ContractMethod].Outputs.UnpackValues(bz)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack values: %w", err)
	}

	// Parse the sqrtPriceX96 from the result.
	sqrtPriceX96 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return sqrtPriceX96, nil
}
