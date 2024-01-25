package kucoin

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/skip-mev/slinky/oracle/config"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	wshandlers "github.com/skip-mev/slinky/providers/base/websocket/handlers"
)

// WebSocketConnHandler handles the process of reading/writing to a websocket
// connection for the Kucoin exchange. Since Kucoin requires dynamic fetching
// of tokens and URLs, this handler will handle the initial connection and inhe
// all other functionality from the default websocket handler.
type WebSocketConnHandler struct {
	*wshandlers.WebSocketConnHandlerImpl

	// requestHandler is the request handler to use for making requests to the
	// Kucoin API.
	requestHandler apihandlers.RequestHandler
}

// NewWebSocketHandler returns a new WebSocketConnHandler.
func NewWebSocketHandler(cfg config.WebSocketConfig, requestHandler apihandlers.RequestHandler) (wshandlers.WebSocketConnHandler, error) {
	handler, err := wshandlers.NewWebSocketHandlerImpl(cfg)
	if err != nil {
		return nil, err
	}

	if requestHandler == nil {
		return nil, fmt.Errorf("request handler cannot be nil")
	}

	return &WebSocketConnHandler{
		WebSocketConnHandlerImpl: handler,
		requestHandler:           requestHandler,
	}, nil
}

// CreateDialer is a function that dynamically creates a new websocket dialer. Per
// the Kucoin documentation, the dialer should first fetch the required token and
// WSS URL from the /api/v1/bullet-public endpoint. The token and URL should then
// be used to create the dialer.
func (h *WebSocketConnHandler) Dial(url string) error {
	// Use the handshake timeout from the config to create a context.
	cfg := h.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.HandshakeTimeout)
	defer cancel()

	// Make the request to the endpoint.
	endpoint := fmt.Sprintf("%s%s", url, BulletPublicEndpoint)
	httpResp, err := h.requestHandler.Do(ctx, endpoint)
	if err != nil {
		return err
	}

	// Decode the response.
	var resp BulletPublicResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
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

	// Create the websocket URL.
	server := resp.Data.InstanceServers[0]
	wss := fmt.Sprintf(WSSEndpoint, server.Endpoint, resp.Data.Token)

	// Update the web socket config with the new WSS and ping interval.
	cfg.WSS = wss
	cfg.PingInterval = time.Duration(server.PingInterval) * time.Millisecond

	conn, _, err := h.CreateDialer().Dial(wss, nil)
	if err != nil {
		return err
	}

	// Set the connection on the handler.
	h.SetConnection(conn)
	return nil
}
