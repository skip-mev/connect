package strategies

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// OracleKeeper is an interface for interacting with the x/oracle state.
//
//go:generate mockery --name OracleKeeper --filename mock_oracle_keeper.go
type OracleKeeper interface {
	GetAllCurrencyPairs(ctx sdk.Context) []oracletypes.CurrencyPair
	GetCurrencyPairFromID(ctx sdk.Context, id uint64) (cp oracletypes.CurrencyPair, found bool)
	GetIDForCurrencyPair(ctx sdk.Context, cp oracletypes.CurrencyPair) (uint64, bool)
}
