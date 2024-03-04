package constants

import (
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	// The following are the default tickers for the oracle sidecar.
	//
	// USD denominated tickers.
	APE_USD      = mmtypes.NewTicker("APE", "USD", 8, 1)
	APTOS_USD    = mmtypes.NewTicker("APT", "USD", 8, 1)
	ATOM_USD     = mmtypes.NewTicker("ATOM", "USD", 8, 1)
	AVAX_USD     = mmtypes.NewTicker("AVAX", "USD", 8, 1)
	BITCOIN_USD  = mmtypes.NewTicker("BTC", "USD", 8, 1)
	CARDANO_USD  = mmtypes.NewTicker("ADA", "USD", 8, 1)
	CELESTIA_USD = mmtypes.NewTicker("TIA", "USD", 8, 1)
	DYDX_USD     = mmtypes.NewTicker("DYDX", "USD", 8, 1)
	ETHEREUM_USD = mmtypes.NewTicker("ETH", "USD", 8, 1)
	OSMOSIS_USD  = mmtypes.NewTicker("OSMO", "USD", 8, 1)
	SOLANA_USD   = mmtypes.NewTicker("SOL", "USD", 8, 1)
	USDC_USD     = mmtypes.NewTicker("USDC", "USD", 8, 1)
	USDT_USD     = mmtypes.NewTicker("USDT", "USD", 8, 1)

	// USDC denominated tickers.
	APE_USDC      = mmtypes.NewTicker("APE", "USDC", 8, 1)
	APTOS_USDC    = mmtypes.NewTicker("APT", "USDC", 8, 1)
	ATOM_USDC     = mmtypes.NewTicker("ATOM", "USDC", 8, 1)
	AVAX_USDC     = mmtypes.NewTicker("AVAX", "USDC", 8, 1)
	BITCOIN_USDC  = mmtypes.NewTicker("BTC", "USDC", 8, 1)
	CARDANO_USDC  = mmtypes.NewTicker("ADA", "USDC", 8, 1)
	CELESTIA_USDC = mmtypes.NewTicker("TIA", "USDC", 8, 1)
	DYDX_USDC     = mmtypes.NewTicker("DYDX", "USDC", 8, 1)
	ETHEREUM_USDC = mmtypes.NewTicker("ETH", "USDC", 8, 1)
	OSMOSIS_USDC  = mmtypes.NewTicker("OSMO", "USDC", 8, 1)
	SOLANA_USDC   = mmtypes.NewTicker("SOL", "USDC", 8, 1)

	// USDT denominated tickers.
	APE_USDT      = mmtypes.NewTicker("APE", "USDT", 8, 1)
	APTOS_USDT    = mmtypes.NewTicker("APT", "USDT", 8, 1)
	ATOM_USDT     = mmtypes.NewTicker("ATOM", "USDT", 8, 1)
	AVAX_USDT     = mmtypes.NewTicker("AVAX", "USDT", 8, 1)
	BITCOIN_USDT  = mmtypes.NewTicker("BTC", "USDT", 8, 1)
	CARDANO_USDT  = mmtypes.NewTicker("ADA", "USDT", 8, 1)
	CELESTIA_USDT = mmtypes.NewTicker("TIA", "USDT", 8, 1)
	DYDX_USDT     = mmtypes.NewTicker("DYDX", "USDT", 8, 1)
	ETHEREUM_USDT = mmtypes.NewTicker("ETH", "USDT", 8, 1)
	OSMOSIS_USDT  = mmtypes.NewTicker("OSMO", "USDT", 8, 1)
	SOLANA_USDT   = mmtypes.NewTicker("SOL", "USDT", 8, 1)
	USDC_USDT     = mmtypes.NewTicker("USDC", "USDT", 8, 1)

	// BTC denominated tickers.
	ETHEREUM_BITCOIN = mmtypes.NewTicker("ETH", "BTC", 8, 1)
)
