package binance

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
)

type (
	// MethodType represents the type of message that is sent to the websocket feed.
	// Each method type corresponds to a different type of message that is sent to the
	// websocket feed.
	MethodType string
	// StreamType represents the type of stream that is sent/received from the websocket
	// feed. Streams are used to determine what instrument you are subscribing to.
	StreamType string
)

const (
	// SubscribeMethod represents a subscribe method. This must be sent as the first message
	// when connecting to the websocket feed.
	//
	// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#subscribe-to-a-stream
	SubscribeMethod MethodType = "SUBSCRIBE"

	// AggregateTradeStream represents the aggregate trade stream. This stream provides
	// trade information that is aggregated for a single taker order.
	//
	// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#aggregate-trade-streams
	AggregateTradeStream StreamType = "aggTrade"
)

// SubscribeMessageRequest represents a subscribe message request. This is used to subscribe
// to Binance websocket streams.
//
// Request
//
//	{
//	  "method": "SUBSCRIBE",
//	  "params": [
//	    "btcusdt@aggTrade",
//	    "btcusdt@depth"
//	  ],
//	  "id": 1
//	}
//
// The ID field iused to uniquely identify the messages going back and forth. By default, the ID
// is a monotonic increasing integer based on the number of messages sent.
//
// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#live-subscribingunsubscribing-to-streams
type SubscribeMessageRequest struct {
	// Method is the method type for the message.
	Method string `json:"method"`
	// Params is the list of streams to subscribe to.
	Params []string `json:"params"`
	// ID is the unique identifier for the message.
	ID int64 `json:"id"`
}

// SubscribeMessageResponse represents a subscribe message response. This is used to determine
// whether the subscription was successful.
//
// Response
//
//	{
//			"result": null,
//			"id": 1
//		}
//
// The ID field is used to uniquely identify the messages going back and forth, the same one sent in
// the initial subscription. The result is null if the subscription was successful.
//
// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#live-subscribingunsubscribing-to-streams
type SubscribeMessageResponse struct {
	// Result is the result of the subscription.
	Result interface{} `json:"result"`
	// ID is the unique identifier for the message.
	ID int64 `json:"id"`
}

// IsEmpty returns true if no data has been set for the message.
func (m *SubscribeMessageResponse) IsEmpty() bool {
	return m.ID == 0 && m.Result == nil
}

// AggregatedTradeMessageResponse represents an aggregated trade message response. This is used to
// represent the aggregated trade data that is received from the Binance websocket.
//
// # Response
//
//	{
//	  "e": "aggTrade",    // Event type
//	  "E": 1672515782136, // Event time
//	  "s": "BNBBTC",      // Symbol
//	  "a": 12345,         // Aggregate trade ID
//	  "p": "0.001",       // Price
//	  "q": "100",         // Quantity
//	  "f": 100,           // First trade ID
//	  "l": 105,           // Last trade ID
//	  "T": 1672515782136, // Trade time
//	  "m": true,          // Is the buyer the market maker?
//	  "M": true           // Ignore
//	}
//
// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#aggregate-trade-streams
type AggregatedTradeMessageResponse struct {
	// StreamType is the stream type.
	StreamType string `json:"e"`
	// Ticker is the symbol.
	Ticker string `json:"s"`
	// Price is the price.
	Price string `json:"p"`
}

// NewSubscribeRequestMessage returns a set of messages to subscribe to the Binance websocket.
func (h *WebSocketHandler) NewSubscribeRequestMessage(instruments []string) ([]handlers.WebsocketEncodedMessage, error) {
	if len(instruments) == 0 {
		return nil, nil
	}

	params := make([]string, len(instruments))
	for i, instrument := range instruments {
		params[i] = fmt.Sprintf("%s@%s", strings.ToLower(instrument), string(AggregateTradeStream))
	}

	// Generate a random ID.
	id := rand.Int63() + 1 // Ensure the ID is not 0.
	msg, err := json.Marshal(SubscribeMessageRequest{
		Method: string(SubscribeMethod),
		Params: params,
		ID:     id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal subscribe message: %w", err)
	}

	// Set the IDs
	h.messageIDs[id] = instruments
	return []handlers.WebsocketEncodedMessage{msg}, nil
}
