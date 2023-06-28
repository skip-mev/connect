package abci

import (
	"cosmossdk.io/math"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// App defines the interface that must be fulfilled by the base application. This
// interface is utilized by the proposal handler to retrieve the state context
// for writing oracle data to state.
//
//go:generate mockery --name App --filename mock_app.go
type App interface {
	GetFinalizeBlockStateCtx() sdk.Context
}

// ValidatorStore defines the interface contract required for calculating
// stake-weighted median prices + total voting power for a given currency pair.
//
//go:generate mockery --name ValidatorStore --filename mock_validator_store.go
type ValidatorStore interface {
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (stakingtypes.Validator, bool)
	TotalBondedTokens(ctx sdk.Context) math.Int
}

// OracleKeeper defines the interface that must be fulfilled by the oracle keeper. This
// interface is utilized by the proposal handler to write oracle data to state for the
// supported assets.
//
//go:generate mockery --name OracleKeeper --filename mock_oracle_keeper.go
type OracleKeeper interface {
	GetAllCurrencyPairs(ctx sdk.Context) []oracletypes.CurrencyPair
	SetPriceForCurrencyPair(ctx sdk.Context, cp oracletypes.CurrencyPair, qp oracletypes.QuotePrice) error
}

// ValidateVoteExtensionsFn defines the function for validating vote extensions. This
// function is not explicitly used to validate the oracle data but rather that
// the signed vote extensions included in the proposal are valid and provide
// a supermajority of vote extensions for the current block. This method is
// expected to be used in ProcessProposal, the expected ctx is the ProcessProposalState's ctx.
type ValidateVoteExtensionsFn func(ctx sdk.Context, currentHeight int64, extendedCommitInfo cometabci.ExtendedCommitInfo) error
