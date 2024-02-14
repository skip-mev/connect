package constants

import (
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	// The following are the default tickers for the oracle sidecar.
	//
	// USD denominated tickers.
	ATOM_USD     = mmtypes.NewTicker("ATOM", "USD", 8, 1)
	AVAX_USD     = mmtypes.NewTicker("AVAX", "USD", 8, 1)
	BITCOIN_USD  = mmtypes.NewTicker("BITCOIN", "USD", 8, 1)
	CELESTIA_USD = mmtypes.NewTicker("CELESTIA", "USD", 8, 1)
	DYDX_USD     = mmtypes.NewTicker("DYDX", "USD", 8, 1)
	ETHEREUM_USD = mmtypes.NewTicker("ETHEREUM", "USD", 8, 1)
	OSMOSIS_USD  = mmtypes.NewTicker("OSMOSIS", "USD", 8, 1)
	SOLANA_USD   = mmtypes.NewTicker("SOLANA", "USD", 8, 1)
	USDC_USD     = mmtypes.NewTicker("USDC", "USD", 8, 1)
	USDT_USD     = mmtypes.NewTicker("USDT", "USD", 8, 1)

	// USDC denominated tickers.
	AVAX_USDC     = mmtypes.NewTicker("AVAX", "USDC", 8, 1)
	BITCOIN_USDC  = mmtypes.NewTicker("BITCOIN", "USDC", 8, 1)
	ETHEREUM_USDC = mmtypes.NewTicker("ETHEREUM", "USDC", 8, 1)
	SOLANA_USDC   = mmtypes.NewTicker("SOLANA", "USDC", 8, 1)

	// USDT denominated tickers.
	ATOM_USDT     = mmtypes.NewTicker("ATOM", "USDT", 8, 1)
	AVAX_USDT     = mmtypes.NewTicker("AVAX", "USDT", 8, 1)
	BITCOIN_USDT  = mmtypes.NewTicker("BITCOIN", "USDT", 8, 1)
	DYDX_USDT     = mmtypes.NewTicker("DYDX", "USDT", 8, 1)
	ETHEREUM_USDT = mmtypes.NewTicker("ETHEREUM", "USDT", 8, 1)
	SOLANA_USDT   = mmtypes.NewTicker("SOLANA", "USDT", 8, 1)
	USDC_USDT     = mmtypes.NewTicker("USDC", "USDT", 8, 1)

	// BTC denominated tickers.
	ETHEREUM_BITCOIN = mmtypes.NewTicker("ETHEREUM", "BITCOIN", 8, 1)
)
