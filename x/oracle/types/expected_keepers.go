package types

import (
	"context"

	"github.com/skip-mev/connect/v2/x/marketmap/types"
)

// MarketMapKeeper is the expected keeper interface for the market map keeper.
//
//go:generate mockery --name MarketMapKeeper --output ./mocks/ --case underscore
type MarketMapKeeper interface {
	GetMarket(ctx context.Context, tickerStr string) (types.Market, error)
}
