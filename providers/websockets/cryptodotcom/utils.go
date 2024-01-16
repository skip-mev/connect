package cryptodotcom

const (
	// URL is the URL used to connect to the Crypto.com websocket API. This can be found here
	// https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html?javascript#websocket-root-endpoints
	// Note that Crypto.com offers a sandbox and production environment.

	// ProductionURL is the URL used to connect to the Crypto.com production websocket API.
	ProductionURL = "wss://stream.crypto.com/exchange/v1/market"

	// SandboxURL is the URL used to connect to the Crypto.com sandbox websocket API. This will
	// return static prices.
	SandboxURL = "wss://uat-stream.3ona.co/exchange/v1/market"
)
