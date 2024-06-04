package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
)

// Keeper defines the interface that must be fulfilled by the oracle keeper. This
// interface is utilized by the PreBlock handler to write oracle data to state for the
// supported assets.
//
//go:generate mockery --name OracleKeeper --filename mock_oracle_keeper.go
type OracleKeeper interface { //golint:ignore
	GetAllCurrencyPairs(ctx sdk.Context) []slinkytypes.CurrencyPair
	SetPriceForCurrencyPair(ctx sdk.Context, cp slinkytypes.CurrencyPair, qp oracletypes.QuotePrice) error
}
