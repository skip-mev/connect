package currencypair

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

// CurrencyPairStrategy is a strategy for generating a unique ID and price representation for a given currency pair.
//
//go:generate mockery --name CurrencyPairStrategy --filename mock_currency_pair_strategy.go
type CurrencyPairStrategy interface { //nolint
	ID(ctx sdk.Context, cp oracletypes.CurrencyPair) (uint64, error)
	FromID(ctx sdk.Context, id uint64) (oracletypes.CurrencyPair, error)
}
