package okx

const (
	// OKX provides a few different URLs for its Websocket API. The URLs can be found
	// in the documentation here: https://www.okx.com/docs-v5/en/?shell#overview-production-trading-services
	// The two production URLs are defined in ProductionURL and ProductionAWSURL. The
	// DemoURL is used for testing purposes.

	// ProductionURL is the public OKX Websocket URL.
	ProductionURL = "wss://ws.okx.com:8443/ws/v5/public"

	// ProductionAWSURL is the public OKX Websocket URL hosted on AWS.
	ProductionAWSURL = "wss://wsaws.okx.com:8443/ws/v5/public"

	// DemoURL is the public OKX Websocket URL for test usage.
	DemoURL = "wss://wspap.okx.com:8443/ws/v5/public?brokerId=9999"
)
