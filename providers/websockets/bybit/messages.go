package bybit

import (
	"encoding/json"
	"fmt"
	"math"

	connectmath "github.com/skip-mev/connect/v2/pkg/math"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
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
	OperationPing      Operation = "ping"
	OperationPong      Operation = "pong"

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

// NewSubscriptionRequestMessage creates subscription messages corresponding to the provided tickers.
// If the number of tickers is greater than 10, the requests will be broken into 10-ticker messages.
func (h *WebSocketHandler) NewSubscriptionRequestMessage(tickers []string) ([]handlers.WebsocketEncodedMessage, error) {
	numTickers := len(tickers)
	if numTickers == 0 {
		return nil, fmt.Errorf("tickers cannot be empty")
	}

	numBatches := int(math.Ceil(float64(numTickers) / float64(h.ws.MaxSubscriptionsPerBatch)))
	msgs := make([]handlers.WebsocketEncodedMessage, numBatches)
	for i := 0; i < numBatches; i++ {
		// Get the tickers for the batch.
		start := i * h.ws.MaxSubscriptionsPerBatch
		end := connectmath.Min((i+1)*h.ws.MaxSubscriptionsPerBatch, numTickers)

		// Create the message for the tickers.
		bz, err := json.Marshal(
			SubscriptionRequest{
				BaseRequest: BaseRequest{
					Op: string(OperationSubscribe),
				},
				Args: tickers[start:end],
			},
		)
		if err != nil {
			return msgs, fmt.Errorf("unable to marshal message: %w", err)
		}

		msgs[i] = bz
	}

	return msgs, nil
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

// NewHeartbeatPingMessage returns the encoded message for sending a heartbeat message to a peer.
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

// BaseResponse is the base structure for responses sent from a peer.
type BaseResponse struct {
	Success bool   `json:"success"`
	RetMsg  string `json:"ret_msg"`
	ConnID  string `json:"conn_id"`
	Op      string `json:"op"`
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

// TickerUpdateMessage is the update sent for a subscribed ticker on the ByBit websocket API.
//
// Example:
//
//	{
//	   "topic": "tickers.BTCUSDT",
//	   "ts": 1673853746003,
//	   "type": "snapshot",
//	   "cs": 2588407389,
//	   "data": {
//	       "symbol": "BTCUSDT",
//	       "lastPrice": "21109.77",
//	       "highPrice24h": "21426.99",
//	       "lowPrice24h": "20575",
//	       "prevPrice24h": "20704.93",
//	       "volume24h": "6780.866843",
//	       "turnover24h": "141946527.22907118",
//	       "price24hPcnt": "0.0196",
//	       "usdIndexPrice": "21120.2400136"
//	   }
type TickerUpdateMessage struct {
	Topic string           `json:"topic"`
	Data  TickerUpdateData `json:"data"`
}

// TickerUpdateData is the data stored inside a ticker update message.
type TickerUpdateData struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"lastPrice"`
}
