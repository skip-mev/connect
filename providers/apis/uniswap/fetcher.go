package uniswap

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	uniswappool "github.com/skip-mev/slinky/providers/apis/uniswap/pool"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"go.uber.org/zap"
)

var _ types.PriceAPIFetcher = (*UniswapPriceFetcher)(nil)

// UniswapPriceFetcher is the Uniswap V3 price fetcher. This fetcher is responsible for
// querying the Uniswap V3 pool contract and returning the price of the pool. The price is
// derived from the slot 0 data of the pool contract. Specifically the sqrtPriceX96 value
// which is the square root of the price of the pool.
//
// To read more about how the price is calculated, see the Uniswap V3 documentation
// https://blog.uniswap.org/uniswap-v3-math-primer.
type UniswapPriceFetcher struct {
	logger  *zap.Logger
	metrics metrics.APIMetrics
	api     config.APIConfig

	// client is the go ethereum client. This is used to interact with the ethereum network.
	client *ethclient.Client
}

// NewUniswapPriceFetcher returns a new Uniswap price fetcher.
func NewUniswapPriceFetcher(
	logger *zap.Logger,
	metrics metrics.APIMetrics,
	api config.APIConfig,
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

	// Dial the ethereum client.
	client, err := ethclient.Dial(api.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ethereum client: %w", err)
	}

	return &UniswapPriceFetcher{
		logger:  logger,
		metrics: metrics,
		api:     api,
		client:  client,
	}, nil
}

// Fetch returns the price of a given ticker. This fetcher expects only 1 ticker to be passed
// in the tickers slice. If more than 1 ticker is passed, an error is returned. The fetcher
// will then query the Uniswap V3 pool contract for the price of the pool. The price is derived
// from the slot 0 data of the pool contract. Specifically the sqrtPriceX96 value which is the

func (u *UniswapPriceFetcher) Fetch(
	ctx context.Context,
	tickers []mmtypes.Ticker,
) types.PriceResponse {
	if len(tickers) != 1 {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(fmt.Errorf("expected 1 ticker"), providertypes.ErrorUnknown),
		)
	}

	var (
		resolved   = make(types.ResolvedPrices)
		unResolved = make(types.UnResolvedPrices)
	)

	// Unmarshal and validate the pool config.
	ticker := tickers[0]
	pool, cfg, err := u.getPool(ticker)
	if err != nil {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(
				fmt.Errorf("failed to get pool: %w", err),
				providertypes.ErrorUnknown,
			),
		)
	}

	price, err := u.getPriceFromPool(ctx, cfg, ticker, pool)
	if err != nil {
		u.logger.Error(
			"failed to get price from pool",
			zap.Error(err),
		)

		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorUnknown),
		)
	}

	// Add the price to the resolved prices.
	resolved[ticker] = types.NewPriceResult(price, time.Now())
	return types.NewPriceResponse(resolved, unResolved)
}

// getPool returns the uniswap pool for the given ticker. This will unmarshal the metadata
// and validate the pool config which contains all required information to query the pool.
// The pool is then returned after querying the ethereum network.
func (u *UniswapPriceFetcher) getPool(
	ticker mmtypes.Ticker,
) (*uniswappool.Uniswap, PoolConfig, error) {
	var cfg PoolConfig
	if err := json.Unmarshal([]byte(ticker.Metadata_JSON), &cfg); err != nil {
		return nil, cfg, fmt.Errorf("failed to unmarshal pool config: %w", err)
	}
	if err := cfg.ValidateBasic(); err != nil {
		return nil, cfg, fmt.Errorf("invalid pool config: %w", err)
	}

	pool, err := uniswappool.NewUniswap(common.HexToAddress(cfg.Address), u.client)
	if err != nil {
		return nil, cfg, fmt.Errorf("failed to query uniswap pool %s: %w", cfg.Address, err)
	}

	return pool, cfg, nil
}

// getPriceFromPool returns the price of the pool based on the slot 0 data.
func (u *UniswapPriceFetcher) getPriceFromPool(
	ctx context.Context,
	cfg PoolConfig,
	ticker mmtypes.Ticker,
	pool *uniswappool.Uniswap,
) (*big.Int, error) {
	// Query the contract.
	slotZero, err := pool.Slot0(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("failed to query slot 0: %w", err)
	}

	// Convert the sqrtPriceX96 to a price. This is the raw, unscaled price.
	price := ConvertSquareRootX96Price(slotZero.SqrtPriceX96)

	// Scale the price to the respective token decimals.
	scaledPrice := ScalePrice(ticker, cfg, price)
	intPrice, _ := scaledPrice.Int(nil)
	return intPrice, nil
}
