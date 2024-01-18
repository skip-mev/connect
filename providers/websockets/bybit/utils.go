package bybit

const (
	// ByBit provides a few different URLs for its Websocket API. The URLs can be found
	// in the documentation here: https://bybit-exchange.github.io/docs/v5/ws/connect
	// The two production URLs are defined in ProductionURL and TestnetURL. The

	// ProductionURL is the public ByBit Websocket URL.
	ProductionURL = "wss://stream.bybit.com/v5/public/spot"

	// TestnetURL is the public ByBit Websocket URL hosted on AWS.
	TestnetURL = "wss://stream-testnet.bybit.com/v5/public/spot"
)
