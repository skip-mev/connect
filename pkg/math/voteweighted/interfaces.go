package voteweighted

import (
	"context"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ValidatorStore defines the interface contract required for calculating stake-weighted median
// prices + total voting power for a given currency pair.
//
//go:generate mockery --srcpkg=github.com/cosmos/cosmos-sdk/x/staking/types --name ValidatorI --filename mock_validator.go
//go:generate mockery --name ValidatorStore --filename mock_validator_store.go
type ValidatorStore interface {
	ValidatorByConsAddr(ctx context.Context, addr sdk.ConsAddress) (stakingtypes.ValidatorI, error)
	TotalBondedTokens(ctx context.Context) (math.Int, error)
}

// PriceStore
type PriceStore interface {
	PriceForCurrencyPair(cp slinkytypes.CurrencyPair) *big.Int
}

type OracleKeeper interface {
	GetPriceForCurrencyPair(ctx sdk.Context, cp slinkytypes.CurrencyPair) (oracletypes.QuotePrice, error)
}
