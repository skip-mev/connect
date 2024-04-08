package raydium

import (
	"context"
	"fmt"
	"math/big"
	"time"

	binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	oracletypes "github.com/skip-mev/slinky/oracle/types"
	oraclemath "github.com/skip-mev/slinky/pkg/math/oracle"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var _ oracletypes.PriceAPIFetcher = (*APIPriceFetcher)(nil)

// SolanaJSONRPCClient is the expected interface for a solana JSON-RPC client according
// to the APIPriceFetcher.
//
//go:generate mockery --name SolanaJSONRPCClient --output ./mocks/ --case underscore
type SolanaJSONRPCClient interface {
	GetMultipleAccountsWithOpts(
		ctx context.Context,
		accounts []solana.PublicKey,
		opts *rpc.GetMultipleAccountsOpts,
	) (out *rpc.GetMultipleAccountsResult, err error)
}

// APIPriceFetcher is responsible for interacting with the solana API and querying information
// about the price of a given currency pair.
type APIPriceFetcher struct {
	// config is the APIConfiguration for this provider
	config config.APIConfig

	// client is the solana JSON-RPC client used to query the API.
	client SolanaJSONRPCClient

	// metaDataPerTicker is a map of ticker.String() -> TickerMetadata
	metaDataPerTicker map[string]TickerMetadata

	// logger
	logger *zap.Logger
}

// NewAPIPriceFetcher returns a new APIPriceFetcher. This method constructs the
// default solana JSON-RPC client in accordance with the config's URL param.
func NewAPIPriceFetcher(
	config config.APIConfig,
	logger *zap.Logger,
) (*APIPriceFetcher, error) {
	return NewAPIPriceFetcherWithClient(
		config,
		rpc.New(config.URL),
		logger,
	)
}

// NewAPIPriceFetcherWithClient returns a new APIPriceFetcher. This method requires
// that the given market + config are valid, otherwise a nil implementation + an error
// will be returned.
func NewAPIPriceFetcherWithClient(
	config config.APIConfig,
	client SolanaJSONRPCClient,
	logger *zap.Logger,
) (*APIPriceFetcher, error) {
	if err := config.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("config for raydium is invalid: %w", err)
	}

	// check fields of config
	if config.Name != Name {
		return nil, fmt.Errorf("configured name is incorrect; expected: %s, got: %s", Name, config.Name)
	}

	if !config.Enabled {
		return nil, fmt.Errorf("config is not enabled")
	}

	// generate metadata per ticker
	return &APIPriceFetcher{
		config:            config,
		client:            client,
		metaDataPerTicker: make(map[string]TickerMetadata),
		logger:            logger.With(zap.String("raydium_api_price_fetcher", Name)),
	}, nil
}

// FetchPrices fetches prices from the solana JSON-RPC API for the given currency-pairs. Specifically
// for each currency-pair,
//   - Query the raydium API base (coin) / quote (pc) token vault addresses
//   - Normalize the token balances by 1e18
//   - Calculate the price as quote / base, and scale by ticker.Decimals
func (pf *APIPriceFetcher) Fetch(
	ctx context.Context,
	tickers []types.ProviderTicker,
) types.PriceResponse {
	// get the acounts to query in order of the tickers given
	expectedNumAccounts := len(tickers) * 2
	accounts := make([]solana.PublicKey, expectedNumAccounts)

	for i, ticker := range tickers {
		metadata, err := pf.updateMetaDataCache(ticker)
		if err != nil {
			return types.NewPriceResponseWithErr(
				tickers,
				providertypes.NewErrorWithCode(
					NoRaydiumMetadataForTickerError(ticker.String()),
					providertypes.ErrorUnknownPair,
				),
			)
		}

		accounts[i*2] = solana.MustPublicKeyFromBase58(metadata.BaseTokenVault.TokenVaultAddress)
		accounts[i*2+1] = solana.MustPublicKeyFromBase58(metadata.QuoteTokenVault.TokenVaultAddress)
	}

	// query the accounts
	// We assume that the solana JSON-RPC response returns all accounts in the order
	// that they were queried, there is not a very good way to handle if this order is incorrect
	// or verify that the order is correct, as there is no way to correlate account data <> address
	accountsResp, err := pf.client.GetMultipleAccountsWithOpts(ctx, accounts, &rpc.GetMultipleAccountsOpts{
		Commitment: rpc.CommitmentFinalized,
		// TODO(nikhil): Keep track of latest height queried as well?
	})
	if err != nil {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(
				SolanaJSONRPCError(err),
				providertypes.ErrorAPIGeneral,
			),
		)
	}

	// expect a base / quote vault account for each ticker queried
	if len(accountsResp.Value) != expectedNumAccounts {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(
				SolanaJSONRPCError(fmt.Errorf("expected %d accounts, got %d", expectedNumAccounts, len(accountsResp.Value))),
				providertypes.ErrorAPIGeneral,
			),
		)
	}

	resolved := make(oracletypes.ResolvedPrices)
	unresolved := make(oracletypes.UnResolvedPrices)
	for i, ticker := range tickers {
		baseAccount := accountsResp.Value[i*2]
		quoteAccount := accountsResp.Value[i*2+1]

		metadata := pf.metaDataPerTicker[ticker.String()]

		// parse the token balances
		baseTokenBalance, err := getScaledTokenBalance(baseAccount, metadata.BaseTokenVault.TokenDecimals)
		if err != nil {
			pf.logger.Debug("error getting base token balance", zap.Error(err))
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(
					SolanaJSONRPCError(err),
					providertypes.ErrorAPIGeneral,
				),
			}
			continue
		}

		quoteTokenBalance, err := getScaledTokenBalance(quoteAccount, metadata.QuoteTokenVault.TokenDecimals)
		if err != nil {
			pf.logger.Debug("error getting quote token balance", zap.Error(err))
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(
					SolanaJSONRPCError(err),
					providertypes.ErrorAPIGeneral,
				),
			}
			continue
		}

		pf.logger.Debug("balances", zap.String("base", baseTokenBalance.String()), zap.String("quote", quoteTokenBalance.String()))

		// calculate the price
		price := calculatePrice(baseTokenBalance, quoteTokenBalance)

		// return the price
		resolved[ticker] = oracletypes.NewPriceResult(price, time.Now())
	}

	return oracletypes.NewPriceResponse(resolved, unresolved)
}

func getScaledTokenBalance(account *rpc.Account, tokenDecimals uint64) (*big.Int, error) {
	// if the account is nil, return error
	if account == nil {
		return nil, fmt.Errorf("account is nil")
	}

	// if the account is empty, return error
	if account.Data == nil {
		return nil, fmt.Errorf("account data is nil")
	}

	// unmarshal the account data into a token account
	var tokenAccount token.Account
	if err := binary.NewBinDecoder(account.Data.GetBinary()).Decode(&tokenAccount); err != nil {
		return nil, err
	}

	// get the token balance + scale by decimals
	balance := new(big.Int).SetUint64(tokenAccount.Amount)
	return oraclemath.ScaleUpCurrencyPairPrice(tokenDecimals, balance)
}

func calculatePrice(baseTokenBalance, quoteTokenBalance *big.Int) *big.Float {
	// calculate the price as quote / base
	return new(big.Float).Quo(
		new(big.Float).SetInt(quoteTokenBalance),
		new(big.Float).SetInt(baseTokenBalance),
	)
}
