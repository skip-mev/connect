package abci

import (
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// ValidatorStore defines the interface contract required for calculating
// stake-weighted median prices + total voting power for a given currency pair.
//
//go:generate mockery --srcpkg=github.com/cosmos/cosmos-sdk/baseapp --name ValidatorStore --filename mock_validator_store.go

// ValidateVoteExtensionsFn defines the function for validating vote extensions. This
// function is not explicitly used to validate the oracle data but rather that
// the signed vote extensions included in the proposal are valid and provide
// a supermajority of vote extensions for the current block. This method is
// expected to be used in ProcessProposal, the expected ctx is the ProcessProposalState's ctx.
type ValidateVoteExtensionsFn func(
	_ sdk.Context,
	_ baseapp.ValidatorStore,
	_ int64,
	_ string,
	_ cometabci.ExtendedCommitInfo,
) error

// OracleKeeper defines the interface that must be fulfilled by the oracle keeper. This
// interface is utilized by the proposal handler to write oracle data to state for the
// supported assets.
//
//go:generate mockery --name OracleKeeper --filename mock_oracle_keeper.go
type OracleKeeper interface {
	GetAllCurrencyPairs(ctx sdk.Context) []oracletypes.CurrencyPair
	SetPriceForCurrencyPair(ctx sdk.Context, cp oracletypes.CurrencyPair, qp oracletypes.QuotePrice) error
}
