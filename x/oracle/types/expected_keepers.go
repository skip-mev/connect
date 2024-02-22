package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

// MarketMapKeeper is the expected keeper interface for the market map keeper.
//
//go:generate mockery --name OracleKeeper --output ./mocks/ --case underscore
type MarketMapKeeper interface {
	GetTicker(ctx sdk.Context, tickerStr string) (types.Ticker, error)
}
