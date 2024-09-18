package okx

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
	// price data that we want. This can later be extended to support other channels. Currently,
	// only the index tickers (spot markets) channel is supported.
	Channel string
	// EventType is the event type. This is the expected event type that we want to receive
	// from the websocket. The event types pertain to subscription events.
	EventType string
)

const (
	// OperationSubscribe is the operation to subscribe to a channel.
	OperationSubscribe Operation = "subscribe"
)

const (
	// TickersChannel is the channel for tickers. This includes the spot price of the instrument.
	//
	// ref: https://www.okx.com/docs-v5/en/#order-book-trading-market-data-ws-tickers-channel
	TickersChannel Channel = "tickers"
)

const (
	// EventSubscribe is the event denoting that we have successfully subscribed to a channel.
	EventSubscribe EventType = "subscribe"
	// EventTickers is the event for tickers. By default, this field will not be populated
	// in a properly formatted message. So we set the default value to an empty string.
	EventTickers EventType = ""
	// EventError is the event for an error.
	EventError EventType = "error"
)

// BaseMessage is utilized to determine the type of message that was received.
type BaseMessage struct {
	// Event is the event that occurred.
	Event string `json:"event" validate:"required"`
}

// SubscribeRequestMessage is the request message for subscribing to a channel. The
// format of the message is:
//
//	{
//			"op": "subscribe",
//			"args": ["<SubscriptionTopic>"]
//	}
//
// Example:
//
//	{
//		"op": "subscribe",
//		"args": [
//			{
//				"channel": "index-tickers",
//				"instId": "LTC-USD-200327"
//			},
//			{
//				"channel": "candle1m",
//				"instId": "LTC-USD-200327"
//			}
//		]
//	}
//
// For more information, see https://www.okx.com/docs-v5/en/?shell#overview-websocket-subscribe
type SubscribeRequestMessage struct {
	// Operation is the operation to perform.
	Operation string `json:"op" validate:"required"`

	// Arguments is the list of arguments for the operation.
	Arguments []SubscriptionTopic `json:"args" validate:"required"`
}

// SubscriptionTopic is the topic to subscribe to.
type SubscriptionTopic struct {
	// Channel is the channel to subscribe to.
	Channel string `json:"channel" validate:"required"`

	// InstrumentID is the instrument ID to subscribe to.
	InstrumentID string `json:"instId" validate:"required"`
}

// NewSubscribeToTickersRequestMessage returns a new SubscribeRequestMessage for subscribing
// to the tickers channel.
func (h *WebSocketHandler) NewSubscribeToTickersRequestMessage(
	instruments []SubscriptionTopic,
) ([]handlers.WebsocketEncodedMessage, error) {
	numInstruments := len(instruments)
	if numInstruments == 0 {
		return nil, fmt.Errorf("instruments cannot be empty")
	}

	numBatches := int(math.Ceil(float64(numInstruments) / float64(h.ws.MaxSubscriptionsPerBatch)))
	msgs := make([]handlers.WebsocketEncodedMessage, numBatches)
	for i := 0; i < numBatches; i++ {
		// Get the instruments for this batch.
		start := i * h.ws.MaxSubscriptionsPerBatch
		end := connectmath.Min((i+1)*h.ws.MaxSubscriptionsPerBatch, numInstruments)

		bz, err := json.Marshal(
			SubscribeRequestMessage{
				Operation: string(OperationSubscribe),
				Arguments: instruments[start:end],
			},
		)
		if err != nil {
			return msgs, err
		}
		msgs[i] = bz
	}

	return msgs, nil
}

// SubscribeResponseMessage is the response message for subscribing to a channel. The
// format of the message is:
// Good Response:
//
//	{
//			"arg": {
//				"channel": "tickers",
//				"instId": "LTC-USD-200327"
//			},
//			"event": "subscribe",
//			"connId": "asdf"
//	}
//
// Bad Response:
//
//	{
//			"event": "error",
//			"code": "60012",
//			"msg": "Invalid request: {\"op\": \"subscribe\", \"argss\":[{ \"channel\" : \"index-tickers\", \"instId\" : \"BTC-USDT\"}]}",
//			"connId": "a4d3ae55"
//	}
//
// For more information, see https://www.okx.com/docs-v5/en/?shell#overview-websocket-subscribe
type SubscribeResponseMessage struct {
	// Arguments is the list of arguments for the operation.
	Arguments SubscriptionTopic `json:"arg"`

	// Event is the event that occurred.
	Event string `json:"event" validate:"required"`

	// ConnectionID is the connection ID.
	ConnectionID string `json:"connId" validate:"required"`

	// Code is the error code.
	Code string `json:"code,omitempty"`

	// Message is the error message. Note that the field will be populated with the same exact
	// initial message that was sent to the websocket.
	Message string `json:"msg,omitempty"`
}

// TickersResponseMessage is the response message for index ticker updates. This message
// type is sent when the index price changes. Price changes are pushed every 100ms if there
// is a change in price. Otherwise, the message is sent every second. The format of the message
// is:
//
//	{
//		"arg": {
//		  "channel": "tickers",
//		  "instId": "BTC-USDT"
//		},
//		"data": [
//		  {
//			"instType": "SPOT",
//			"instId": "BTC-USDT",
//			"last": "9999.99",
//			"lastSz": "0.1",
//			"askPx": "9999.99",
//			"askSz": "11",
//			"bidPx": "8888.88",
//			"bidSz": "5",
//			"open24h": "9000",
//			"high24h": "10000",
//			"low24h": "8888.88",
//			"volCcy24h": "2222",
//			"vol24h": "2222",
//			"sodUtc0": "2222",
//			"sodUtc8": "2222",
//			"ts": "1597026383085"
//		  }
//		]
//	}
//
// For more information, see https://www.okx.com/docs-v5/en/?shell#public-data-websocket-index-tickers-channel
type TickersResponseMessage struct {
	// Arguments is the list of arguments for the operation.
	Arguments SubscriptionTopic `json:"arg" validate:"required"`

	// Data is the list of index ticker data.
	Data []IndexTicker `json:"data" validate:"required"`
}

// IndexTicker is the index ticker data.
type IndexTicker struct {
	// ID is the instrument ID.
	ID string `json:"instId" validate:"required"`

	// LastPrice is the last price.
	LastPrice string `json:"last" validate:"required"`
}
