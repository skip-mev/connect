package kucoin

const (
	// URL is the Kucoin websocket URL. This URL specifically points to the public
	// spot and maring REST API.
	URL = "https://api.kucoin.com"

	// BulletPublicEndpoint is the endpoint to connect to for the public feed. This
	// requires a POST request with no body to receive a token and endpoints to
	// connect to.
	BulletPublicEndpoint = "/api/v1/bullet-public"

	// SuccessCode is the success code returned from the Kucoin API.
	SuccessCode = "200000"
)

// BulletPublicResponse represents the response from the bullet-public endpoint
// for the Kucoin exchange. This response is utilized when initially connecting
// to the websocket feed. Specifically, the response is utilized to determine the
// token and endpoints to connect to.
//
//	{
//		"code": "200000",
//		"data": {
//		  	"token": "token1234567890",
//		  	"instanceServers": [
//					{
//			  			"endpoint": "wss://ws-api-spot.kucoin.com/", // It is recommended to use a dynamic URL, which may change
//			  			"encrypt": true,
//			  			"protocol": "websocket",
//			  			"pingInterval": 18000,
//			  			"pingTimeout": 10000
//					}
//		  		]
//			}
//	}
//
// ref: https://www.kucoin.com/docs/websocket/basic-info/apply-connect-token/public-token-no-authentication-required-
type BulletPublicResponse struct {
	// Code is the response code.
	Code string `json:"code"`

	// Data is the response data.
	Data BulledPublicResponseData `json:"data"`
}

// BulledPublicResponseData is the data field of the BulletPublicResponse.
type BulledPublicResponseData struct {
	// Token is the token to use for authentication.
	Token string `json:"token"`

	// InstanceServers is the list of instance servers to connect to.
	InstanceServers []BulletPublicResponseInstanceServer `json:"instanceServers"`
}

// BulletPublicResponseInstanceServer is the instance server to connect to.
type BulletPublicResponseInstanceServer struct {
	// Endpoint is the endpoint to connect to.
	Endpoint string `json:"endpoint"`

	// Encrypt is a flag that indicates if the connection should be encrypted.
	Encrypt bool `json:"encrypt"`

	// Protocol is the protocol to use for the connection.
	Protocol string `json:"protocol"`

	// PingInterval is the interval to ping the server.
	PingInterval int64 `json:"pingInterval"`

	// PingTimeout is the timeout for the ping.
	PingTimeout int64 `json:"pingTimeout"`
}
