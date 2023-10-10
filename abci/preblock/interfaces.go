package preblock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// OracleKeeper defines the interface that must be fulfilled by the oracle keeper. This
// interface is utilized by the PreBlock handler to write oracle data to state for the
// supported assets.
//
//go:generate mockery --name OracleKeeper --filename mock_oracle_keeper.go
type OracleKeeper interface {
	GetAllCurrencyPairs(ctx sdk.Context) []oracletypes.CurrencyPair
	SetPriceForCurrencyPair(ctx sdk.Context, cp oracletypes.CurrencyPair, qp oracletypes.QuotePrice) error
}
