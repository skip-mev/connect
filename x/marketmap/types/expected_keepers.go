package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/x/oracle/types"
)

// OracleKeeper defines the interface that must be fulfilled by the oracle keeper.
//
//go:generate mockery --name OracleKeeper --filename mock_oracle_keeper.go
type OracleKeeper interface {
	CreateCurrencyPair(ctx sdk.Context, cp types.CurrencyPair) error
	RemoveCurrencyPair(ctx sdk.Context, cp types.CurrencyPair)
}
