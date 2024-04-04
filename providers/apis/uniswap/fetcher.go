package uniswap

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
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	uniswappool "github.com/skip-mev/slinky/providers/apis/uniswap/pool"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ types.PriceAPIFetcher = (*UniswapPriceFetcher)(nil)

// UniswapPriceFetcher is the Uniswap V3 price fetcher. This fetcher is responsible for
// querying the Uniswap V3 pool contract and returning the price of the pool. The price is
// derived from the slot 0 data of the pool contract. Specifically the sqrtPriceX96 value
// which is the square root of the price of the pool.
//
// To read more about how the price is calculated, see the Uniswap V3 documentation
// https://blog.uniswap.org/uniswap-v3-math-primer.
//
// Additionally, we utilize the eth client's BatchCallContext to batch the calls to the
// ethereum network this is more performant than making individual calls or the multi call
// contract: https://docs.chainstack.com/docs/http-batch-request-vs-multicall-contract.
type UniswapPriceFetcher struct {
	logger  *zap.Logger
	metrics metrics.APIMetrics
	api     config.APIConfig

	// client is the go ethereum client. This is used to interact with the ethereum network.
	client EVMClient
	// abi is the uniswap v3 pool abi. This is used to pack the slot0 call to the pool contract
	// and parse the result.
	abi *abi.ABI
	// payload is the packed slot0 call to the pool contract. Since the slot0 payload is the same
	// for all pools, we can reuse this payload for all pools.
	payload []byte
}

// NewUniswapPriceFetcher returns a new Uniswap price fetcher.
func NewUniswapPriceFetcher(
	logger *zap.Logger,
	metrics metrics.APIMetrics,
	api config.APIConfig,
	client EVMClient,
) (*UniswapPriceFetcher, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if metrics == nil {
		return nil, fmt.Errorf("metrics cannot be nil")
	}

	if api.Name != Name {
		return nil, fmt.Errorf("expected api config name %s, got %s", Name, api.Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api config for %s is not enabled", Name)
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config: %w", err)
	}

	abi, err := uniswappool.UniswapMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get uniswap abi: %w", err)
	}

	payload, err := abi.Pack("slot0")
	if err != nil {
		return nil, fmt.Errorf("failed to pack slot0: %w", err)
	}

	return &UniswapPriceFetcher{
		logger:  logger,
		metrics: metrics,
		api:     api,
		client:  client,
		abi:     abi,
		payload: payload,
	}, nil
}

// Fetch returns the price of a given ticker. This fetcher expects only 1 ticker to be passed
// in the tickers slice. If more than 1 ticker is passed, an error is returned. The fetcher
// will then query the Uniswap V3 pool contract for the price of the pool. The price is derived
// from the slot 0 data of the pool contract. Specifically the sqrtPriceX96 value which is the
// square root of the price of the pool.
func (u *UniswapPriceFetcher) Fetch(
	ctx context.Context,
	tickers []mmtypes.Ticker,
) types.PriceResponse {
	var (
		resolved   = make(types.ResolvedPrices)
		unResolved = make(types.UnResolvedPrices)
	)

	// Create a batch element for each ticker and pool.
	batchElems := make([]rpc.BatchElem, len(tickers))
	pools := make([]PoolConfig, len(tickers))
	for i, ticker := range tickers {
		pool, err := u.getPool(ticker)
		if err != nil {
			return types.NewPriceResponseWithErr(
				tickers,
				providertypes.NewErrorWithCode(
					fmt.Errorf("failed to get pool: %w", err),
					providertypes.ErrorUnknown,
				),
			)
		}

		// Create a batch element for the ticker and pool.
		batchElems[i] = u.createBatchElement(pool)
		pools[i] = pool
	}

	// Batch call to the ethereum network.
	if err := u.client.BatchCallContext(ctx, batchElems); err != nil {
		u.logger.Error(
			"failed to batch call to ethereum network for all tickers",
			zap.Error(err),
		)

		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorUnknown),
		)
	}

	// Parse the result from the batch call for each ticker.
	for i, ticker := range tickers {
		result := batchElems[i]
		if result.Error != nil {
			u.logger.Error(
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
		sqrtPriceX96, err := u.parseSqrtPriceX96(result.Result)
		if err != nil {
			u.logger.Error(
				"failed to parse sqrt price x96",
				zap.String("ticker", ticker.String()),
				zap.Error(err),
			)

			unResolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(
					err,
					providertypes.ErrorUnknown,
				),
			}

			continue
		}

		// Convert the sqrtPriceX96 to a price. This is the raw, unscaled price.
		price := ConvertSquareRootX96Price(sqrtPriceX96)

		// Scale the price to the respective token decimals.
		scaledPrice := ScalePrice(ticker, pools[i], price)
		intPrice, _ := scaledPrice.Int(nil)
		resolved[ticker] = types.NewPriceResult(intPrice, time.Now())
	}

	// Add the price to the resolved prices.
	return types.NewPriceResponse(resolved, unResolved)
}

// getPool returns the uniswap pool for the given ticker. This will unmarshal the metadata
// and validate the pool config which contains all required information to query the pool.
// The pool is then returned after querying the ethereum network.
func (u *UniswapPriceFetcher) getPool(
	ticker mmtypes.Ticker,
) (PoolConfig, error) {
	var cfg PoolConfig
	if err := json.Unmarshal([]byte(ticker.Metadata_JSON), &cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal pool config: %w", err)
	}
	if err := cfg.ValidateBasic(); err != nil {
		return cfg, fmt.Errorf("invalid pool config: %w", err)
	}

	return cfg, nil
}

// createBatchElement creates a batch element for the given ticker and pool. This will be utilized
// to batch the calls to the ethereum network to retrieve all pricing information.
func (u *UniswapPriceFetcher) createBatchElement(
	pool PoolConfig,
) rpc.BatchElem {
	var result string
	return rpc.BatchElem{
		Method: "eth_call",
		Args: []interface{}{
			map[string]interface{}{
				"to":   common.HexToAddress(pool.Address),
				"data": hexutil.Bytes(u.payload),
			},
			"latest", // latest signifies the latest block.
		},
		Result: &result,
	}
}

// parseSqrtPriceX96 parses the sqrtPriceX96 from the result of the batch call. The sqrtPriceX96
// is the square root of the price of the pool. This is the raw, unscaled price.
func (u *UniswapPriceFetcher) parseSqrtPriceX96(
	result interface{},
) (*big.Int, error) {
	r, ok := result.(*string)
	if !ok {
		return nil, fmt.Errorf("expected result to be a string, got %T", result)
	}

	bz, err := hexutil.Decode(*r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex result: %w", err)
	}

	out, err := u.abi.Methods["slot0"].Outputs.UnpackValues(bz)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack values: %w", err)
	}

	// Parse the sqrtPriceX96 from the result.
	sqrtPriceX96 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return sqrtPriceX96, nil
}
