package voteweighted

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
)

// ValidatorStore defines the interface contract required for calculating stake-weighted median
// prices + total voting power for a given currency pair.
//
//go:generate mockery --srcpkg=github.com/cosmos/cosmos-sdk/x/staking/types --name ValidatorI --filename mock_validator.go
//go:generate mockery --name ValidatorStore --filename mock_validator_store.go
type ValidatorStore interface {
	ValidatorByConsAddr(ctx context.Context, addr sdk.ConsAddress) (stakingtypes.ValidatorI, error)
	TotalBondedTokens(ctx context.Context) (math.Int, error)
}

// CCValidatorStore defines the interface contract required for the cross chain validator consumer store.
//
//go:generate mockery --name CCValidatorStore --filename mock_cc_validator_store.go
type CCValidatorStore interface {
	GetAllCCValidator(ctx sdk.Context) []types.CrossChainValidator
	GetCCValidator(ctx sdk.Context, addr []byte) (types.CrossChainValidator, bool)
}
