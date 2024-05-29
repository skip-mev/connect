package constants

import (
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
)

var (
	BITCOIN_USD = slinkytypes.NewCurrencyPair("BTC", "USD")
	ETHEREUM_USD = slinkytypes.NewCurrencyPair("ETH", "USD")
	ETHEREUM_USDT = slinkytypes.NewCurrencyPair("ETH", "USDT")
	PEPE_USD = slinkytypes.NewCurrencyPair("PEPE", "USD")
	ETHEREUM_USDC = slinkytypes.NewCurrencyPair("ETH", "USDC")
	USDT_USD = slinkytypes.NewCurrencyPair("USDT", "USD")
)
