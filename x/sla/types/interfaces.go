package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SlashingKeeper defines the interface that must be fulfilled by the slashing keeper.
//
//go:generate mockery --name SlashingKeeper --filename mock_slashing_keeper.go
type SlashingKeeper interface {
	// Slash attempts to slash a validator. The slash is delegated to the staking
	// module to make the necessary validator changes. It specifies no intraction reason.
	Slash(
		ctx context.Context,
		consAddr sdk.ConsAddress,
		infractionHeight,
		power int64,
		slashFactor math.LegacyDec,
	) (amount math.Int, err error)
}

// StakingKeeper defines the interface that must be fulfilled by the staking keeper.
//
//go:generate mockery --name StakingKeeper --filename mock_staking_keeper.go
type StakingKeeper interface {
	// GetLastValidatorPower returns the last recorded power of a validator. Returns zero if
	// the operator was not a validator last block.
	GetLastValidatorPower(ctx context.Context, operator sdk.ValAddress) (power int64, err error)
}
