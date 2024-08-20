package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	incentivetypes "github.com/skip-mev/connect/v2/x/incentives/types"
)

// BankKeeper defines the expected interface that the bank-keeper dependency must implement.
//
//go:generate mockery --name BankKeeper --output ./mocks/ --case underscore
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amounts sdk.Coins) error
}

// OracleKeeper defines the expected interface that the oracle-keeper dependency must implement.
//
//go:generate mockery --name OracleKeeper --output ./mocks/ --case underscore
type OracleKeeper interface {
	HasCurrencyPair(ctx sdk.Context, cp slinkytypes.CurrencyPair) bool
}

// StakingKeeper defines the expected interface that the staking-keeper dependency must implement.
//
//go:generate mockery --name StakingKeeper --output ./mocks/ --case underscore
type StakingKeeper interface {
	Slash(ctx context.Context, consAddr sdk.ConsAddress, infractionHeight, power int64, slashFactor math.LegacyDec) (math.Int, error)
	GetValidatorByConsAddr(ctx context.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, err error)
	BondDenom(ctx context.Context) (string, error)
}

// IncentiveKeeper defines the expected interface that the incentive-keeper dependency must implement.
//
//go:generate mockery --name IncentiveKeeper --output ./mocks/ --case underscore
type IncentiveKeeper interface {
	AddIncentives(ctx sdk.Context, incentives []incentivetypes.Incentive) error
}
