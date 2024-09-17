package aggregator

import (
	"fmt"
	"math/big"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/abci/strategies/codec"
	"github.com/skip-mev/connect/v2/abci/strategies/currencypair"
	connectabci "github.com/skip-mev/connect/v2/abci/types"
	vetypes "github.com/skip-mev/connect/v2/abci/ve/types"
	"github.com/skip-mev/connect/v2/aggregator"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
)

// Vote encapsulates the validator and oracle data contained within a vote extension.
type Vote struct {
	// ConsAddress is the validator that submitted the vote extension.
	ConsAddress sdk.ConsAddress
	// OracleVoteExtension
	OracleVoteExtension vetypes.OracleVoteExtension
}

// GetOracleVotes returns all oracle vote extensions that were injected into
// the block. Note that all vote extensions included are necessarily valid at this point
// because the vote extensions were validated by the vote extension and proposal handlers.
func GetOracleVotes(
	proposal [][]byte,
	veCodec codec.VoteExtensionCodec,
	extCommitCodec codec.ExtendedCommitCodec,
) ([]Vote, error) {
	if len(proposal) < connectabci.NumInjectedTxs {
		return nil, connectabci.MissingCommitInfoError{}
	}

	extendedCommitInfo, err := extCommitCodec.Decode(proposal[connectabci.OracleInfoIndex])
	if err != nil {
		return nil, connectabci.CodecError{
			Err: fmt.Errorf("error decoding extended-commit-info: %w", err),
		}
	}

	votes := make([]Vote, len(extendedCommitInfo.Votes))
	for i, voteInfo := range extendedCommitInfo.Votes {
		voteExtension, err := veCodec.Decode(voteInfo.VoteExtension)
		if err != nil {
			return nil, connectabci.CodecError{
				Err: fmt.Errorf("error decoding vote-extension: %w", err),
			}
		}

		votes[i] = Vote{
			ConsAddress:         voteInfo.Validator.Address,
			OracleVoteExtension: voteExtension,
		}
	}

	return votes, nil
}

// VoteAggregator is an interface that defines the methods for aggregating oracle votes into a set of prices.
// This object holds both the aggregated price resulting from a given set of votes, and the prices
// reported by each validator.
//
//go:generate mockery --name VoteAggregator --filename mock_vote_aggregator.go
type VoteAggregator interface {
	// AggregateOracleVotes ingresses vote information which contains all
	// vote extensions each validator extended in the previous block. it is important
	// to note that
	//  1. The vote extension may be nil, in which case the validator is not providing
	//     any oracle data for the current block. This could have occurred because the
	//     validator was offline, or its local oracle service was down.
	//  2. The vote extension may contain prices updates for only a subset of currency pairs.
	//     This could have occurred because the price providers for the validator were
	//     offline, or the price providers did not provide a price update for a given
	//     currency pair.
	//
	// In order for a currency pair to be included in the final oracle price, the currency
	// pair must be provided by a super-majority (2/3+) of validators. This is enforced by the
	// price aggregator but can be replaced by the application.
	//
	// Notice: This method overwrites the VoteAggregator's local view of prices.
	AggregateOracleVotes(ctx sdk.Context, votes []Vote) (map[connecttypes.CurrencyPair]*big.Int, error)

	// GetPriceForValidator gets the prices reported by a given validator. This method depends
	// on the prices from the latest set of aggregated votes.
	GetPriceForValidator(validator sdk.ConsAddress) map[connecttypes.CurrencyPair]*big.Int
}

func NewDefaultVoteAggregator(
	logger log.Logger,
	aggregateFn aggregator.AggregateFnFromContext[string, map[connecttypes.CurrencyPair]*big.Int],
	strategy currencypair.CurrencyPairStrategy,
) VoteAggregator {
	return &DefaultVoteAggregator{
		logger: logger,
		priceAggregator: aggregator.NewDataAggregator(
			aggregator.WithAggregateFnFromContext(aggregateFn),
		),
		currencyPairStrategy: strategy,
	}
}

type DefaultVoteAggregator struct {
	// validator address -> currency-pair -> price
	priceAggregator *aggregator.DataAggregator[string, map[connecttypes.CurrencyPair]*big.Int]

	// decoding prices / currency-pair ids
	currencyPairStrategy currencypair.CurrencyPairStrategy

	logger log.Logger
}

func (dva *DefaultVoteAggregator) AggregateOracleVotes(ctx sdk.Context, votes []Vote) (map[connecttypes.CurrencyPair]*big.Int, error) {
	// Reset the price aggregator and set the aggregationFn to use the latest application-state.
	dva.priceAggregator.ResetProviderData()

	// Iterate through all vote extensions and consolidate all price info before
	// aggregating.
	for _, vote := range votes {
		consAddrStr := vote.ConsAddress.String()

		if err := dva.addVoteToAggregator(ctx, consAddrStr, vote.OracleVoteExtension); err != nil {
			dva.logger.Error(
				"failed to add vote to aggregator",
				"validator_address", consAddrStr,
				"err", err,
			)

			return nil, err
		}
	}

	// Compute the final prices for each currency pair.
	dva.priceAggregator.AggregateDataFromContext(ctx)
	prices := dva.priceAggregator.GetAggregatedData()

	dva.logger.Debug(
		"aggregated oracle data",
		"num_prices", len(prices),
	)

	return prices, nil
}

// addVoteToAggregator consolidates the oracle data from a single validator
// into the price aggregator. The oracle data is provided in the form of a vote
// extension. The vote extension contains the prices for each currency pair that
// the validator is providing for the current block.
func (dva *DefaultVoteAggregator) addVoteToAggregator(ctx sdk.Context, address string, oracleData vetypes.OracleVoteExtension) error {
	if len(oracleData.Prices) == 0 {
		return nil
	}

	// Format all prices into a map of currency pair -> price.
	prices := make(map[connecttypes.CurrencyPair]*big.Int, len(oracleData.Prices))
	for cpID, priceBz := range oracleData.Prices {
		if len(priceBz) > connectabci.MaximumPriceSize {
			dva.logger.Debug(
				"failed to store price, bytes are too long",
				"currency_pair_id", cpID,
				"num_bytes", len(priceBz),
			)
			continue
		}

		// Convert the asset into a currency pair.
		cp, err := dva.currencyPairStrategy.FromID(ctx, cpID)
		if err != nil {
			dva.logger.Debug(
				"failed to convert currency pair id to currency pair",
				"currency_pair_id", cpID,
				"err", err,
			)

			// If the currency pair is not supported, continue.
			continue
		}

		price, err := dva.currencyPairStrategy.GetDecodedPrice(ctx, cp, priceBz)
		if err != nil {
			dva.logger.Debug(
				"failed to decode price",
				"currency_pair_id", cpID,
				"err", err,
			)

			// If the price cannot be decoded, continue.
			continue
		}

		prices[cp] = price
	}

	dva.logger.Debug(
		"adding oracle prices to aggregator",
		"num_prices", len(prices),
		"validator_address", address,
	)

	dva.priceAggregator.SetProviderData(address, prices)

	return nil
}

func (dva *DefaultVoteAggregator) GetPriceForValidator(validator sdk.ConsAddress) map[connecttypes.CurrencyPair]*big.Int {
	consAddrStr := validator.String()
	return dva.priceAggregator.GetDataByProvider(consAddrStr)
}
