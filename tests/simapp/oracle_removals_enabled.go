//go:build oracle_removals_enabled
package simapp

import (
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

func ProvideMarketMapKeeperOptions(k oracletypes.MarketMapKeeper) func(k oracletypes.MarketMapKeeper) oracletypes.MarketMapKeeper {
	return func(k oracletypes.MarketMapKeeper) oracletypes.MarketMapKeeper {
		return nil // nullify the market-map keeper
	}
}