package bybit

import (
	"encoding/json"
	"fmt"

	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
)

type (
	// Operation is the operation to perform. This is used to construct subscription messages
	// when initially connecting to the websocket. This can later be extended to support
	// other operations.
	Operation string

	// Channel is the channel to subscribe to. The channel is used to determine the type of
	// price data that we want. This can later be extended to support other channels.
	Channel string
)

const (
	// OperationSubscribe is the operation to subscribe to a channel.
	OperationSubscribe Operation = "subscribe"

	OperationPing Operation = "ping"

	// TickerChannel is the channel for spot price updates.
	TickerChannel Channel = "tickers"
)

type BaseRequest struct {
	ReqID string `json:"req_id"`
	Op    string `json:"op"`
}

// SubscriptionRequest is a request to the server to subscribe to ticker updates for currency pairs.
//
// Example:
//
//	{
//	   "req_id": "test", // optional
//	   "op": "subscribe",
//	   "args": [
//	       "orderbook.1.BTCUSDT",
//	       "publicTrade.BTCUSDT",
//	       "orderbook.1.ETHUSDT"
//	   ]
//	}
type SubscriptionRequest struct {
	BaseRequest
	Args []string `json:"args"`
}

func NewSubscriptionRequestMessage(tickers []string) ([]handlers.WebsocketEncodedMessage, error) {
	if len(tickers) == 0 {
		return nil, fmt.Errorf("tickers cannot be empty")
	}

	bz, err := json.Marshal(
		SubscriptionRequest{
			BaseRequest: BaseRequest{
				Op: string(OperationPing),
			},
			Args: tickers,
		},
	)

	return []handlers.WebsocketEncodedMessage{bz}, err
}

// HeartbeatPing is the ping sent to the server.
//
// Example:
//
//	{
//	   "req_id": "100010",
//	   "op": "ping"
//	}
type HeartbeatPing struct {
	BaseRequest
}

func NewHeartbeatPingMessage() ([]handlers.WebsocketEncodedMessage, error) {
	bz, err := json.Marshal(
		HeartbeatPing{
			BaseRequest{
				Op: string(OperationPing),
			},
		},
	)

	return []handlers.WebsocketEncodedMessage{bz}, err
}

type BaseResponse struct {
	Success bool `json:"success"`

	RetMsg string `json:"ret_msg"`

	ConnID string `json:"conn_id"`

	Op string `json:"op"`
}

// HeartbeatPong is the pong sent back from the server after a ping.
//
// Example:
//
//	{
//	   "success": true,
//	   "ret_msg": "pong",
//	   "conn_id": "0970e817-426e-429a-a679-ff7f55e0b16a",
//	   "op": "ping"
//	}
type HeartbeatPong struct {
	BaseResponse
}

// SubscriptionResponse is the response for a subscribe event.
//
// Example:
//
//	{
//	   "success": true,
//	   "ret_msg": "subscribe",
//	   "conn_id": "2324d924-aa4d-45b0-a858-7b8be29ab52b",
//	   "req_id": "10001",
//	   "op": "subscribe"
//	}
type SubscriptionResponse struct {
	BaseResponse
	ReqID string `json:"req_id"`
}

type TickerUpdateMessage struct {
	Topic string           `json:"topic"`
	Data  TickerUpdateData `json:"data"`
}

type TickerUpdateData struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"lastPrice"`
}
