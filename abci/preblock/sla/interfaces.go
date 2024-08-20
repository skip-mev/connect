package sla

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slakeeper "github.com/skip-mev/connect/v2/x/sla/keeper"
)

// Keeper defines the interface that must be fulfilled by the SLA keeper.
//
//go:generate mockery --name Keeper --filename mock_sla_keeper.go
type Keeper interface {
	// UpdatePriceFeeds will update all price feed incentives given the latest updates.
	UpdatePriceFeeds(ctx sdk.Context, updates slakeeper.PriceFeedUpdates) error
}

// OracleKeeper defines the interface that must be fulfilled by the oracle keeper.
//
//go:generate mockery --name OracleKeeper --filename mock_oracle_keeper.go
type OracleKeeper interface {
	// GetAllCurrencyPairs returns all CurrencyPairs that have currently been stored to state.
	GetAllCurrencyPairs(ctx sdk.Context) []slinkytypes.CurrencyPair
}

// StakingKeeper defines the interface that must be fulfilled by the staking keeper.
//
//go:generate mockery --name StakingKeeper --filename mock_staking_keeper.go
type StakingKeeper interface {
	// GetBondedValidatorsByPower returns all bonded validators that have currently been stored to state.
	GetBondedValidatorsByPower(ctx context.Context) ([]stakingtypes.Validator, error)
}
