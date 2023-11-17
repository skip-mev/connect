package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	incentivetypes "github.com/skip-mev/slinky/x/incentives/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// BankKeeper defines the expected interface that the bank-keeper dependency must implement.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amounts sdk.Coins) error
}

// OracleKeeper defines the expected interface that the oracle-keeper dependency must implement.
type OracleKeeper interface {
	GetNonceForCurrencyPair(ctx sdk.Context, cp oracletypes.CurrencyPair) (uint64, error)
}

// StakingKeeper defines the expected interface that the staking-keeper dependency must implement.
type StakingKeeper interface {
	Slash(ctx context.Context, consAddr sdk.ConsAddress, infractionHeight, power int64, slashFactor math.LegacyDec) (math.Int, error)
	GetValidatorByConsAddr(ctx context.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, err error)
	BondDenom(ctx context.Context) (string, error)
}

// IncentiveKeeper defines the expected interface that the incentive-keeper dependency must implement.
type IncentiveKeeper interface {
	AddIncentives(ctx sdk.Context, incentives []incentivetypes.Incentive) error
}
