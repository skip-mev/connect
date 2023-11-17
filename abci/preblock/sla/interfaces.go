package sla

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	slakeeper "github.com/skip-mev/slinky/x/sla/keeper"
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
	GetAllCurrencyPairs(ctx sdk.Context) []oracletypes.CurrencyPair
}

// StakingKeeper defines the interface that must be fulfilled by the staking keeper.
//
//go:generate mockery --name StakingKeeper --filename mock_staking_keeper.go
type StakingKeeper interface {
	// GetBondedValidatorsByPower returns all bonded validators that have currently been stored to state.
	GetBondedValidatorsByPower(ctx context.Context) ([]stakingtypes.Validator, error)
}
