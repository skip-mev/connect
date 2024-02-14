package constants

import (
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	ATOM_USD         = mmtypes.NewTicker("ATOM", "USD", 8, 1)
	BITCOIN_USD      = mmtypes.NewTicker("BITCOIN", "USD", 8, 1)
	CELESTIA_USD     = mmtypes.NewTicker("CELESTIA", "USD", 8, 1)
	DYDX_USD         = mmtypes.NewTicker("DYDX", "USD", 8, 1)
	ETHEREUM_BITCOIN = mmtypes.NewTicker("ETHEREUM", "BITCOIN", 8, 1)
	ETHEREUM_USD     = mmtypes.NewTicker("ETHEREUM", "USD", 8, 1)
	OSMOSIS_USD      = mmtypes.NewTicker("OSMOSIS", "USD", 8, 1)
	SOLANA_USD       = mmtypes.NewTicker("SOLANA", "USD", 8, 1)
)
