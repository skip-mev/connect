package bitstamp

import (
	"encoding/json"
	"fmt"

	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
)

type (
	// ChannelType defines the channel types that can be subscribed to.
	ChannelType string

	// EventTypes defines the event types that can be sent to the server.
	EventTypes string
)

const (
	// TickerChannel is the ticker channel. A subscription to this channel
	// will provide real-time trade data, updated every time a trade occurs.
	//
	// ref: https://www.bitstamp.net/websocket/v2/
	TickerChannel ChannelType = "live_trades_%s"

	// HeartbeatEvent is the heartbeat event. This event is sent to the server
	// to keep the connection alive.
	//
	// ref: https://www.bitstamp.net/websocket/v2/
	HeartbeatEvent EventTypes = "bts:heartbeat"

	// SubscriptionEvent is the subscription event. This event is sent to the
	// server to subscribe to a channel.
	//
	// ref: https://www.bitstamp.net/websocket/v2/
	SubscriptionEvent EventTypes = "bts:subscribe"

	// ReconnectEvent is the reconnect event. After you receive this request,
	// you will have a few seconds to reconnect. Without doing so, you will
	// automatically be disconnected. If you send reconnection request, you
	// will be placed to a new server. Consequentially, you can continue without
	// any message loss.
	//
	// ref: https://www.bitstamp.net/websocket/v2/
	ReconnectEvent EventTypes = "bts:request_reconnect"
)

// BaseMessage is utilized to define the base message structure for the Bitstamp
// websocket feed.
type BaseMessage struct {
	// Event is the event type.
	Event EventTypes `json:"event"`
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

// NewSubscriptionRequestMessage returns a new subscription request message.
func NewSubscriptionRequestMessage(channels string) ([]handlers.WebsocketEncodedMessage, error) {
	if len(channels) == 0 {
		return nil, fmt.Errorf("no instruments provided")
	}

	bz, err := json.Marshal(SubscriptionRequestMessage{
		BaseMessage: BaseMessage{
			Event: SubscriptionEvent,
		},
		Data: SubscriptionRequestMessageData{
			Channel: channels,
	)
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
//		"id": "live_trades_btcusd",
//		"amount": 0.001,
//		"amount_str": "0.00100000",
//		"price": 10000,
//		"price_str": "10000.00",
//		"type": 1,
//		"timestamp": "1597679171",
//		"microtimestamp": "1597679171000000",
//		"buy_order_id": 0,
//		"sell_order_id": 0
//	}
//
// ref: https://www.bitstamp.net/websocket/v2/
type TickerResponseMessage struct {
	// ID is the trade unique ID.
	ID string `json:"id"`

	// PriceStr is the trade price represented in string format.
	PriceStr string `json:"price_str"`
}
