package bitstamp

import (
	"encoding/json"
	"fmt"

	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

type (
	// ChannelType defines the channel types that can be subscribed to.
	ChannelType string

	// EventType defines the event types that can be sent to the server.
	EventType string
)

const (
	// TickerChannel is the ticker channel. A subscription to this channel
	// will provide real-time trade data, updated every time a trade occurs.
	//
	// ref: https://www.bitstamp.net/websocket/v2/
	TickerChannel ChannelType = "live_trades_"

	// HeartbeatEvent is the heartbeat event. This event is sent to the server
	// to keep the connection alive.
	//
	// ref: https://www.bitstamp.net/websocket/v2/
	HeartbeatEvent EventType = "bts:heartbeat"

	// SubscriptionEvent is the subscription event. This event is sent to the
	// server to subscribe to a channel.
	//
	// ref: https://www.bitstamp.net/websocket/v2/
	SubscriptionEvent EventType = "bts:subscribe"

	// SubscriptionSucceededEvent is the subscription succeeded event. This
	// event is received after a subscription request is made to the ticker
	// channel.
	//
	// ref: https://www.bitstamp.net/websocket/v2/
	SubscriptionSucceededEvent EventType = "bts:subscription_succeeded"

	// ReconnectEvent is the reconnect event. After you receive this request,
	// you will have a few seconds to reconnect. Without doing so, you will
	// automatically be disconnected. If you send reconnection request, you
	// will be placed to a new server. Consequently, you can continue without
	// any message loss.
	//
	// ref: https://www.bitstamp.net/websocket/v2/
	ReconnectEvent EventType = "bts:request_reconnect"

	// TradeEvent is the trade event. This event is received after a subscription
	// request is made to the ticker channel. It contains the relevant ticker data.
	//
	// ref: https://www.bitstamp.net/websocket/v2/
	TradeEvent EventType = "trade"
)

// BaseMessage is utilized to define the base message structure for the Bitstamp
// websocket feed.
type BaseMessage struct {
	// Event is the event type.
	Event string `json:"event"`
}

// SubscriptionRequestMessage represents a subscription request message. Once an
// initial connection is established, clients can send a subscription request
// message to receive updates for a specific topic.
//
//	{
//	    "event": "bts:subscribe",
//	    "data": {
//	        "channel": "[channel_name]"
//	    }
//	}
//
// ref: https://www.bitstamp.net/websocket/v2/
type SubscriptionRequestMessage struct {
	BaseMessage

	// Data is the subscription data.
	Data SubscriptionRequestMessageData `json:"data"`
}

// SubscriptionRequestMessageData is the data field of the
// SubscriptionRequestMessage.
type SubscriptionRequestMessageData struct {
	// Channel is the channel to subscribe to.
	Channel string `json:"channel"`
}

// NewSubscriptionRequestMessages returns a new subscription request message
// for a given set of channels.
func NewSubscriptionRequestMessages(channels []string) ([]handlers.WebsocketEncodedMessage, error) {
	if len(channels) == 0 {
		return nil, fmt.Errorf("no instruments provided")
	}

	msgs := make([]handlers.WebsocketEncodedMessage, len(channels))
	for i, channel := range channels {
		msg, err := NewSubscriptionRequestMessage(channel)
		if err != nil {
			return nil, err
		}

		msgs[i] = msg
	}

	return msgs, nil
}

// NewSubscriptionRequestMessage returns a new subscription request message.
func NewSubscriptionRequestMessage(channel string) (handlers.WebsocketEncodedMessage, error) {
	bz, err := json.Marshal(SubscriptionRequestMessage{
		BaseMessage: BaseMessage{
			Event: string(SubscriptionEvent),
		},
		Data: SubscriptionRequestMessageData{
			Channel: channel,
		},
	})
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// SubscriptionResponseMessage represents a subscription response message. This
// message is received after a subscription request is made to the ticker
// channel.
//
//	{
//			"event" : "bts:subscription_succeeded",
//			"channel" : "live_trades_btcusd",
//			"data":{}
//		}
//
// ref: https://www.bitstamp.net/websocket/v2/
type SubscriptionResponseMessage struct {
	BaseMessage

	// Channel is the channel that was subscribed to.
	Channel string `json:"channel"`
}

// TickerResponseMessage represents a ticker response message. This message is
// received after a subscription request is made to the ticker channel. It contains
// the relevant ticker data.
//
// id:				Trade unique ID.
// amount:			Trade amount.
// amount_str:		Trade amount represented in string format.
// price:			Trade price.
// price_str:		Trade price represented in string format.
// type:			Trade type (0 - buy; 1 - sell).
// timestamp:		Trade timestamp.
// microtimestamp:	Trade microtimestamp.
// buy_order_id:	Trade buy order ID.
// sell_order_id:	Trade sell order ID.
//
// example:
//
//	{
//		"channel":"live_trades_ethusd",
//		"event":"trade"
//		"data": {
//				"id":317319339,
//				"timestamp":"1706302442",
//				"amount":0.062,
//				"amount_str":"0.06200000",
//				"price":2253.7,
//				"price_str":"2253.7",
//				"type":0,
//				"microtimestamp":"1706302442078000",
//				"buy_order_id":1709946746556417,
//				"sell_order_id":1709946744614913
//			}
//	}
//
// ref: https://www.bitstamp.net/websocket/v2/
type TickerResponseMessage struct {
	BaseMessage

	// Channel is the channel that was subscribed to.
	Channel string `json:"channel"`

	// Data is the ticker data.
	Data TickerData `json:"data"`
}

// TickerData is the data field of the TickerResponseMessage.
type TickerData struct {
	// PriceStr is the price represented in string format.
	PriceStr string `json:"price_str"`

	// Channel is the channel that was subscribed to.
	Channel string `json:"channel"`
}

const (
	// ExpectedTickerLength is the expected length of the ID field.
	// This field contains the channel name and the currency pair.
	ExpectedTickerLength = 2

	// TickerChannelIndex is the index of the ticker channel in the ID field.
	TickerChannelIndex = 0

	// TickerCurrencyPairIndex is the index of the ticker currency pair in the
	// ID field.
	TickerCurrencyPairIndex = 1
)

// NewHeartbeatRequestMessage returns a new heartbeat request message.
func NewHeartbeatRequestMessage() ([]handlers.WebsocketEncodedMessage, error) {
	bz, err := json.Marshal(BaseMessage{
		Event: string(HeartbeatEvent),
	})
	if err != nil {
		return nil, err
	}

	return []handlers.WebsocketEncodedMessage{bz}, nil
}

// NewReconnectRequestMessage returns a new reconnect request message.
func NewReconnectRequestMessage() ([]handlers.WebsocketEncodedMessage, error) {
	bz, err := json.Marshal(BaseMessage{
		Event: string(ReconnectEvent),
	})
	if err != nil {
		return nil, err
	}

	return []handlers.WebsocketEncodedMessage{bz}, nil
}
