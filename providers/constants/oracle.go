package constants

import (
	"math/big"

	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

type (
	// PriceAPIQueryHandler is a type alias for the API query handler that is used to fetch price data
	// from the price providers.
	PriceAPIQueryHandler = apihandlers.APIQueryHandler[mmtypes.Ticker, *big.Int]

	// PriceAPIDataHandler is a type alias for the API data handler that is used to fetch resolve http
	// requests and parse the response into price data.
	PriceAPIDataHandler = apihandlers.APIDataHandler[mmtypes.Ticker, *big.Int]
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
	ATOM_USDC     = mmtypes.NewTicker("ATOM", "USDC", 8, 1)
	AVAX_USDC     = mmtypes.NewTicker("AVAX", "USDC", 8, 1)
	BITCOIN_USDC  = mmtypes.NewTicker("BITCOIN", "USDC", 8, 1)
	CELESTIA_USDC = mmtypes.NewTicker("CELESTIA", "USDC", 8, 1)
	DYDX_USDC     = mmtypes.NewTicker("DYDX", "USDC", 8, 1)
	ETHEREUM_USDC = mmtypes.NewTicker("ETHEREUM", "USDC", 8, 1)
	OSMOSIS_USDC  = mmtypes.NewTicker("OSMOSIS", "USDC", 8, 1)
	SOLANA_USDC   = mmtypes.NewTicker("SOLANA", "USDC", 8, 1)

	// USDT denominated tickers.
	ATOM_USDT     = mmtypes.NewTicker("ATOM", "USDT", 8, 1)
	AVAX_USDT     = mmtypes.NewTicker("AVAX", "USDT", 8, 1)
	BITCOIN_USDT  = mmtypes.NewTicker("BITCOIN", "USDT", 8, 1)
	CELESTIA_USDT = mmtypes.NewTicker("CELESTIA", "USDT", 8, 1)
	DYDX_USDT     = mmtypes.NewTicker("DYDX", "USDT", 8, 1)
	ETHEREUM_USDT = mmtypes.NewTicker("ETHEREUM", "USDT", 8, 1)
	OSMOSIS_USDT  = mmtypes.NewTicker("OSMOSIS", "USDT", 8, 1)
	SOLANA_USDT   = mmtypes.NewTicker("SOLANA", "USDT", 8, 1)
	USDC_USDT     = mmtypes.NewTicker("USDC", "USDT", 8, 1)

	// BTC denominated tickers.
	ETHEREUM_BITCOIN = mmtypes.NewTicker("ETHEREUM", "BITCOIN", 8, 1)
)
