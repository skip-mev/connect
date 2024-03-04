package constants

import (
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	// The following are the default tickers for the oracle sidecar.
	//
	// USD denominated tickers.
	APE_USD       = mmtypes.NewTicker("APE", "USD", 8, 1)
	APTOS_USD     = mmtypes.NewTicker("APT", "USD", 8, 1)
	ARBITRUM_USD  = mmtypes.NewTicker("ARB", "USD", 8, 1)
	ATOM_USD      = mmtypes.NewTicker("ATOM", "USD", 8, 1)
	AVAX_USD      = mmtypes.NewTicker("AVAX", "USD", 8, 1)
	BCH_USD       = mmtypes.NewTicker("BCH", "USD", 8, 1)
	BITCOIN_USD   = mmtypes.NewTicker("BTC", "USD", 8, 1)
	BLUR_USD      = mmtypes.NewTicker("BLUR", "USD", 8, 1)
	CARDANO_USD   = mmtypes.NewTicker("ADA", "USD", 8, 1)
	CELESTIA_USD  = mmtypes.NewTicker("TIA", "USD", 8, 1)
	CHAINLINK_USD = mmtypes.NewTicker("LINK", "USD", 8, 1)
	COMPOUND_USD  = mmtypes.NewTicker("COMP", "USD", 8, 1)
	CURVE_USD     = mmtypes.NewTicker("CRV", "USD", 8, 1)
	DOGE_USD      = mmtypes.NewTicker("DOGE", "USD", 8, 1)
	DYDX_USD      = mmtypes.NewTicker("DYDX", "USD", 8, 1)
	ETC_USD       = mmtypes.NewTicker("ETC", "USD", 8, 1)
	ETHEREUM_USD  = mmtypes.NewTicker("ETH", "USD", 8, 1)
	FILECOIN_USD  = mmtypes.NewTicker("FIL", "USD", 8, 1)
	LIDO_USD      = mmtypes.NewTicker("LDO", "USD", 8, 1)
	LITECOIN_USD  = mmtypes.NewTicker("LTC", "USD", 8, 1)
	OSMOSIS_USD   = mmtypes.NewTicker("OSMO", "USD", 8, 1)
	POLKADOT_USD  = mmtypes.NewTicker("DOT", "USD", 8, 1)
	POLYGON_USD   = mmtypes.NewTicker("MATIC", "USD", 8, 1)
	MAKER_USD     = mmtypes.NewTicker("MKR", "USD", 8, 1)
	NEAR_USD      = mmtypes.NewTicker("NEAR", "USD", 8, 1)
	OPTIMISM_USD  = mmtypes.NewTicker("OP", "USD", 8, 1)
	PEPE_USD      = mmtypes.NewTicker("PEPE", "USD", 8, 1)
	RIPPLE_USD    = mmtypes.NewTicker("XRP", "USD", 8, 1)
	SEI_USD       = mmtypes.NewTicker("SEI", "USD", 8, 1)
	SOLANA_USD    = mmtypes.NewTicker("SOL", "USD", 8, 1)
	STELLAR_USD   = mmtypes.NewTicker("XLM", "USD", 8, 1)
	SUI_USD       = mmtypes.NewTicker("SUI", "USD", 8, 1)
	TRON_USD      = mmtypes.NewTicker("TRX", "USD", 8, 1)
	UNISWAP_USD   = mmtypes.NewTicker("UNI", "USD", 8, 1)
	WORLD_USD     = mmtypes.NewTicker("WLD", "USD", 8, 1)
	USDC_USD      = mmtypes.NewTicker("USDC", "USD", 8, 1)
	USDT_USD      = mmtypes.NewTicker("USDT", "USD", 8, 1)

	// USDC denominated tickers.
	APE_USDC       = mmtypes.NewTicker("APE", "USDC", 8, 1)
	APTOS_USDC     = mmtypes.NewTicker("APT", "USDC", 8, 1)
	ARBITRUM_USDC  = mmtypes.NewTicker("ARB", "USDC", 8, 1)
	ATOM_USDC      = mmtypes.NewTicker("ATOM", "USDC", 8, 1)
	AVAX_USDC      = mmtypes.NewTicker("AVAX", "USDC", 8, 1)
	BCH_USDC       = mmtypes.NewTicker("BCH", "USDC", 8, 1)
	BITCOIN_USDC   = mmtypes.NewTicker("BTC", "USDC", 8, 1)
	BLUR_USDC      = mmtypes.NewTicker("BLUR", "USDC", 8, 1)
	CARDANO_USDC   = mmtypes.NewTicker("ADA", "USDC", 8, 1)
	CELESTIA_USDC  = mmtypes.NewTicker("TIA", "USDC", 8, 1)
	CHAINLINK_USDC = mmtypes.NewTicker("LINK", "USDC", 8, 1)
	COMPOUND_USDC  = mmtypes.NewTicker("COMP", "USDC", 8, 1)
	CURVE_USDC     = mmtypes.NewTicker("CRV", "USDC", 8, 1)
	DOGE_USDC      = mmtypes.NewTicker("DOGE", "USDC", 8, 1)
	DYDX_USDC      = mmtypes.NewTicker("DYDX", "USDC", 8, 1)
	ETC_USDC       = mmtypes.NewTicker("ETC", "USDC", 8, 1)
	ETHEREUM_USDC  = mmtypes.NewTicker("ETH", "USDC", 8, 1)
	FILECOIN_USDC  = mmtypes.NewTicker("FIL", "USDC", 8, 1)
	LIDO_USDC      = mmtypes.NewTicker("LDO", "USDC", 8, 1)
	LITECOIN_USDC  = mmtypes.NewTicker("LTC", "USDC", 8, 1)
	OSMOSIS_USDC   = mmtypes.NewTicker("OSMO", "USDC", 8, 1)
	POLKADOT_USDC  = mmtypes.NewTicker("DOT", "USDC", 8, 1)
	POLYGON_USDC   = mmtypes.NewTicker("MATIC", "USDC", 8, 1)
	MAKER_USDC     = mmtypes.NewTicker("MKR", "USDC", 8, 1)
	NEAR_USDC      = mmtypes.NewTicker("NEAR", "USDC", 8, 1)
	OPTIMISM_USDC  = mmtypes.NewTicker("OP", "USDC", 8, 1)
	PEPE_USDC      = mmtypes.NewTicker("PEPE", "USDC", 8, 1)
	RIPPLE_USDC    = mmtypes.NewTicker("XRP", "USDC", 8, 1)
	SEI_USDC       = mmtypes.NewTicker("SEI", "USDC", 8, 1)
	SUI_USDC       = mmtypes.NewTicker("SUI", "USDC", 8, 1)
	SOLANA_USDC    = mmtypes.NewTicker("SOL", "USDC", 8, 1)
	STELLAR_USDC   = mmtypes.NewTicker("XLM", "USDC", 8, 1)
	TRON_USDC      = mmtypes.NewTicker("TRX", "USDC", 8, 1)
	UNISWAP_USDC   = mmtypes.NewTicker("UNI", "USDC", 8, 1)
	WORLD_USDC     = mmtypes.NewTicker("WLD", "USDC", 8, 1)

	// USDT denominated tickers.
	APE_USDT       = mmtypes.NewTicker("APE", "USDT", 8, 1)
	APTOS_USDT     = mmtypes.NewTicker("APT", "USDT", 8, 1)
	ARBITRUM_USDT  = mmtypes.NewTicker("ARB", "USDT", 8, 1)
	ATOM_USDT      = mmtypes.NewTicker("ATOM", "USDT", 8, 1)
	AVAX_USDT      = mmtypes.NewTicker("AVAX", "USDT", 8, 1)
	BCH_USDT       = mmtypes.NewTicker("BCH", "USDT", 8, 1)
	BITCOIN_USDT   = mmtypes.NewTicker("BTC", "USDT", 8, 1)
	BLUR_USDT      = mmtypes.NewTicker("BLUR", "USDT", 8, 1)
	CARDANO_USDT   = mmtypes.NewTicker("ADA", "USDT", 8, 1)
	CELESTIA_USDT  = mmtypes.NewTicker("TIA", "USDT", 8, 1)
	CHAINLINK_USDT = mmtypes.NewTicker("LINK", "USDT", 8, 1)
	COMPOUND_USDT  = mmtypes.NewTicker("COMP", "USDT", 8, 1)
	CURVE_USDT     = mmtypes.NewTicker("CRV", "USDT", 8, 1)
	DOGE_USDT      = mmtypes.NewTicker("DOGE", "USDT", 8, 1)
	DYDX_USDT      = mmtypes.NewTicker("DYDX", "USDT", 8, 1)
	ETC_USDT       = mmtypes.NewTicker("ETC", "USDT", 8, 1)
	ETHEREUM_USDT  = mmtypes.NewTicker("ETH", "USDT", 8, 1)
	FILECOIN_USDT  = mmtypes.NewTicker("FIL", "USDT", 8, 1)
	LIDO_USDT      = mmtypes.NewTicker("LDO", "USDT", 8, 1)
	LITECOIN_USDT  = mmtypes.NewTicker("LTC", "USDT", 8, 1)
	OSMOSIS_USDT   = mmtypes.NewTicker("OSMO", "USDT", 8, 1)
	POLKADOT_USDT  = mmtypes.NewTicker("DOT", "USDT", 8, 1)
	POLYGON_USDT   = mmtypes.NewTicker("MATIC", "USDT", 8, 1)
	MAKER_USDT     = mmtypes.NewTicker("MKR", "USDT", 8, 1)
	NEAR_USDT      = mmtypes.NewTicker("NEAR", "USDT", 8, 1)
	OPTIMISM_USDT  = mmtypes.NewTicker("OP", "USDT", 8, 1)
	PEPE_USDT      = mmtypes.NewTicker("PEPE", "USDT", 8, 1)
	RIPPLE_USDT    = mmtypes.NewTicker("XRP", "USDT", 8, 1)
	SEI_USDT       = mmtypes.NewTicker("SEI", "USDT", 8, 1)
	SOLANA_USDT    = mmtypes.NewTicker("SOL", "USDT", 8, 1)
	STELLAR_USDT   = mmtypes.NewTicker("XLM", "USDT", 8, 1)
	SUI_USDT       = mmtypes.NewTicker("SUI", "USDT", 8, 1)
	TRON_USDT      = mmtypes.NewTicker("TRX", "USDT", 8, 1)
	UNISWAP_USDT   = mmtypes.NewTicker("UNI", "USDT", 8, 1)
	USDC_USDT      = mmtypes.NewTicker("USDC", "USDT", 8, 1)
	WORLD_USDT     = mmtypes.NewTicker("WLD", "USDT", 8, 1)

	// BTC denominated tickers.
	ETHEREUM_BITCOIN = mmtypes.NewTicker("ETH", "BTC", 8, 1)
)
