package kucoin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
	apihandlers "github.com/skip-mev/connect/v2/providers/base/api/handlers"
	wshandlers "github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

const (
	// BulletPublicEndpoint is the endpoint to connect to the public WSS feed. This
	// requires a POST request with no body to receive a token and endpoints to
	// connect to.
	BulletPublicEndpoint = "/api/v1/bullet-public"

	// SuccessCode is the success code returned from the KuCoin API.
	SuccessCode = "200000"

	// WebSocketProtocol is the expected protocol type for the KuCoin websocket feed.
	WebSocketProtocol = "websocket"
)

// BulletPublicResponse represents the response from the bullet-public endpoint
// for the KuCoin exchange. This response is utilized when initially connecting
// to the websocket feed. Specifically, the response is utilized to determine the
// token and endpoints to connect to.
//
//	{
//		"code": "200000",
//		"data": {
//		  	"token": "token1234567890", // Used to suffix the WSS URL
//		  	"instanceServers": [
//					{
//			  			"endpoint": "wss://ws-api-spot.kucoin.com/", // It is recommended to use a dynamic URL, which may change
//			  			"encrypt": true,
//			  			"protocol": "websocket",
//			  			"pingInterval": 18000, // We use this as the ping interval
//			  			"pingTimeout": 10000 // We use this as the read timeout
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

// PreDialHook is a function that is called before the connection is established.
// This function is used to fetch the token and WSS URL from the KuCoin API.
func PreDialHook(cfg config.APIConfig, requestHandler apihandlers.RequestHandler) wshandlers.PreDialHook {
	return func(handler *wshandlers.WebSocketConnHandlerImpl) error {
		resp, err := fetchCredentials(cfg, requestHandler)
		if err != nil {
			return err
		}

		// Check if the response was successful.
		if resp.Code != SuccessCode {
			return fmt.Errorf("failed to fetch token and URL from Kucoin: %s", resp.Code)
		}

		// There must be at least one instance server that we can connect to. Otherwise,
		// we cannot connect to the websocket feed.
		if len(resp.Data.InstanceServers) == 0 {
			return fmt.Errorf("no instance servers returned from %s", cfg.Name)
		}

		// Determine if there is a server that supports the desired protocol.
		var (
			server BulletPublicResponseInstanceServer
			found  bool
		)
		for _, s := range resp.Data.InstanceServers {
			if s.Protocol == WebSocketProtocol {
				server = s
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("no instance servers support the %s protocol", WebSocketProtocol)
		}

		// Create the websocket URL.
		wss := fmt.Sprintf(WSSEndpoint, server.Endpoint, resp.Data.Token)

		// Update the websocket config with the new WSS and ping interval.
		cfg := handler.GetConfig()

		cfg.Endpoints = []config.Endpoint{
			{
				URL: wss,
			},
		}
		cfg.PingInterval = time.Duration(server.PingInterval) * time.Millisecond
		cfg.ReadTimeout = time.Duration(server.PingTimeout) * time.Millisecond

		// Set the connection and updated config on the handler.
		handler.SetConfig(cfg)
		return nil
	}
}

// fetchCredentials is used to fetch the token and WSS URL from the KuCoin API.
func fetchCredentials(cfg config.APIConfig, requestHandler apihandlers.RequestHandler) (*BulletPublicResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if len(cfg.Endpoints) == 0 {
		return nil, fmt.Errorf("no endpoints provided")
	}

	// Make the request to the endpoint.
	endpoint := fmt.Sprintf("%s%s", cfg.Endpoints[0].URL, BulletPublicEndpoint)
	httpResp, err := requestHandler.Do(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request for credentials failed with status %s", httpResp.Status)
	}

	// Decode the response.
	var resp BulletPublicResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
