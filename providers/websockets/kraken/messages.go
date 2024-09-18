package kraken

import (
	"encoding/json"
	"fmt"
	"math"

	connectmath "github.com/skip-mev/connect/v2/pkg/math"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

type (
	// Event correspond to the various message types that are sent to the client.
	Event string

	// Status correspond to the various status types that are sent to the client.
	Status string

	// Channel correspond to the various channels that are available for
	// subscription.
	Channel string
)

const (
	// SystemStatusEvent is the event name that is sent to the client when the
	// connection is first established. The status notifies the client of any
	// outages or planned maintenance.
	//
	// https://docs.kraken.com/websockets/#message-systemStatus
	SystemStatusEvent Event = "systemStatus"

	// HeartbeatEvent is the event name that is sent to the client periodically
	// to ensure that the connection is still alive.
	//
	// https://docs.kraken.com/websockets/#message-heartbeat
	HeartbeatEvent Event = "heartbeat"

	// SubscriptionStatusEvent is the event name that is used to send a subscribe
	// message to the server.
	//
	// https://docs.kraken.com/websockets/#message-subscribe
	SubscriptionStatusEvent Event = "subscriptionStatus"

	// SubscribeEvent is the event name that is used to send a subscribe message
	// to the server.
	//
	// https://docs.kraken.com/websockets/#message-subscribe
	SubscribeEvent Event = "subscribe"
)

const (
	// OnlineStatus is the status that is sent to the client when the connection
	// is first established. The status notifies the client that the connection
	// is online.
	OnlineStatus Status = "online"

	// MaintenanceStatus is the status that is sent to the client when the
	// connection is first established. The status notifies the client that the
	// connection is online, but that there is ongoing maintenance.
	MaintenanceStatus Status = "maintenance"

	// SubscribedStatus is the status that is sent to the client when the server
	// has received the subscription request.
	SubscribedStatus Status = "subscribed"

	// ErrorStatus is the status that is sent to the client when the server has
	// received the subscription request.
	ErrorStatus Status = "error"
)

const (
	// TickerChannel is the channel name for the ticker channel.
	//
	// https://docs.kraken.com/websockets/#message-ticker
	TickerChannel Channel = "ticker"
)

// BaseMessage is the template used to determine the type of message that is
// received from the server.
type BaseMessage struct {
	// Event is the event name that is sent to the client.
	Event string `json:"event"`
}

// SystemStatusResponseMessage is the message that is sent to the client when the
// connection is first established.
//
//	{
//			"connectionID": 8628615390848610000,
//			"event": "systemStatus",
//			"status": "online",
//			"version": "1.0.0"
//	}
//
// ref: https://docs.kraken.com/websockets/#message-systemStatus
type SystemStatusResponseMessage struct {
	// ConnectionID is the unique identifier for the connection.
	ConnectionID uint64 `json:"connectionID"`

	// Event is the event name that is sent to the client when the connection is
	// first established.
	Event string `json:"event"`

	// Status is the status that is sent to the client when the connection is
	// first established.
	Status string `json:"status"`

	// Version is the version of the API.
	Version string `json:"version"`
}

// HeartbeatResponseMessage is the message that is sent to the client
// periodically to ensure that the connection is still alive. In particular,
// the server will send a heartbeat if no subscription traffic is received
// within a 60-second period.
//
//	{
//			"event": "heartbeat",
//	}
//
// ref: https://docs.kraken.com/websockets/#message-heartbeat
type HeartbeatResponseMessage struct {
	// Event is the event name that is sent to the client periodically to ensure
	// that the connection is still alive.
	Event string `json:"event"`
}

// SubscribeRequestMessage is the message that is sent to the server to subscribe
// to a channel.
//
//	{
//			"event": "subscribe",
//			"pair": [
//				"XBT/USD",
//				"XBT/EUR"
//			],
//			"subscription": {
//				"name": "ticker"
//			}
//	}
//
// ref: https://docs.kraken.com/websockets/#message-subscribe
type SubscribeRequestMessage struct {
	// Event is the event name that is sent to the server to subscribe to a
	// channel.
	Event string `json:"event"`

	// Pair is the asset pair to subscribe to.
	Pair []string `json:"pair"`

	// Subscription is the subscription details.
	Subscription Subscription `json:"subscription"`
}

// Subscription is the subscription details.
type Subscription struct {
	// Name is the name of the subscription.
	Name string `json:"name"`
}

// NewSubscribeRequestMessage returns a new SubscribeRequestMessage with the
// given asset pairs.
func (h *WebSocketHandler) NewSubscribeRequestMessage(
	instruments []string,
) ([]handlers.WebsocketEncodedMessage, error) {
	numInstruments := len(instruments)
	if numInstruments == 0 {
		return nil, fmt.Errorf("no instruments specified")
	}

	numBatches := int(math.Ceil(float64(numInstruments) / float64(h.ws.MaxSubscriptionsPerBatch)))
	msgs := make([]handlers.WebsocketEncodedMessage, numBatches)
	for i := 0; i < numBatches; i++ {
		// Get the instruments for the batch.
		start := i * h.ws.MaxSubscriptionsPerBatch
		end := connectmath.Min((i+1)*h.ws.MaxSubscriptionsPerBatch, numInstruments)

		bz, err := json.Marshal(
			SubscribeRequestMessage{
				Event: string(SubscribeEvent),
				Pair:  instruments[start:end],
				Subscription: Subscription{
					Name: string(TickerChannel),
				},
			},
		)
		if err != nil {
			return msgs, err
		}
		msgs[i] = bz
	}
	return msgs, nil
}

// SubscribeResponseMessage is the message that is sent to the client when the
// server has received the subscription request.
//
// Good response:
//
//	{
//		"channelID": 10001,
//		"channelName": "ticker",
//		"event": "subscriptionStatus",
//		"pair": "XBT/EUR",
//		"status": "subscribed",
//		"subscription": {
//			"name": "ticker"
//		}
//	}
//
// Bad response:
//
//		{
//			"errorMessage": "Subscription depth not supported",
//			"event": "subscriptionStatus",
//			"pair": "XBT/USD",
//			"status": "error",
//			"subscription": {
//				"depth": 42,
//				"name": "book"
//			}
//	}
//
// ref: https://docs.kraken.com/websockets/#message-subscriptionStatus
type SubscribeResponseMessage struct {
	// ChannelID is the channel ID.
	ChannelID uint64 `json:"channelID"`

	// ChannelName is the channel name.
	ChannelName string `json:"channelName"`

	// Event is the event name that is sent to the client when the server has
	// received the subscription request.
	Event string `json:"event"`

	// Pair is the asset pair that was subscribed to.
	Pair string `json:"pair"`

	// Status is the status that is sent to the client when the server has
	// received the subscription request.
	Status string `json:"status"`

	// Subscription is the subscription details.
	Subscription Subscription `json:"subscription"`

	// ErrorMessage is the error message that is sent to the client when the
	// server has received the subscription request.
	ErrorMessage string `json:"errorMessage"`
}

// TickerResponseMessage is the message that is sent to the client when the
// server has a price update for the subscribed asset pair. This is specific
// to the ticker subscription.
//
//		[
//	  	0, 						// ChannelID
//	  	{
//	    	"a": [ 				// Ask array
//	      		"5525.40000", 	// Best ask price
//	     	 	1, 				// Whole lot volume
//	     	 	"1.000" 		// Lot volume
//	    	],
//	    	"b": [ 				// Bid array
//	      		"5525.10000", 	// Best bid price
//	     	 	1, 				// Whole lot volume
//	      		"1.000" 		// Lot volume
//	    	],
//	    	"c": [ 				// Close array
//	      		"5525.10000", 	// Price
//	      		"0.00398963" 	// Lot volume
//	    	],
//	    	"h": [ 				// High price array
//	      		"5783.00000", 	// Today
//	      		"5783.00000" 	// Last 24 hours
//	    	],
//	    	"l": [ 				// Low price array
//	      		"5505.00000", 	// Today
//	      		"5505.00000"	// Last 24 hours
//	   	 	],
//	    	"o": [				// Open price array
//	     	 	"5760.70000", 	// Today
//	     	 	"5763.40000"	// Last 24 hours
//	   	 	],
//	  	 	"p": [ 				// Volume weighted average price array <- This is the value we want
//	   	   		"5631.44067", 	// Value Today
//	   	   		"5653.78939" 	// Value 24h
//	  	  	],
//	   	 	"t": [ 				// Number of trades array
//	   	   		11493,			// Today
//	   	   		16267 			// Last 24 hours
//	  	 	],
//	   	 	"v": [				// Volume array
//	    	  "2634.11501494", 	// Value Today
//	    	  "3591.17907851"	// Value 24h
//	   	 	]
//	 	 },
//	  	"ticker", 				// Channel name
//	  	"XBT/USD" 				// Asset pair
//
// ]
//
// ref: https://docs.kraken.com/websockets/#message-ticker
type TickerResponseMessage struct {
	// ChannelID is the channel ID.
	ChannelID int

	// TickerData is the ticker data corresponding to the asset pair.
	TickerData TickerData

	// ChannelName is the channel name.
	ChannelName string

	// Pair is the asset pair that was subscribed to.
	Pair string
}

const (
	// ExpectedTickerResponseMessageLength is the expected length of the ticker
	// response message.
	ExpectedTickerResponseMessageLength = 4

	// ChannelIDIndex is the index of the channel ID in the ticker response
	// message.
	ChannelIDIndex = iota - 1

	// TickerDataIndex is the index of the ticker data in the ticker response
	// message.
	TickerDataIndex

	// ChannelNameIndex is the index of the channel name in the ticker response
	// message.
	ChannelNameIndex

	// PairIndex is the index of the asset pair in the ticker response message.
	PairIndex
)

// TickerData is the ticker data.
type TickerData struct {
	// VolumeWeightedAveragePrice is the volume weighted average price.
	VolumeWeightedAveragePrice []string `json:"p"`
}

const (
	// TodayPriceIndex is the index of the today's price in the ticker's
	// VolumeWeightedAveragePrice array.
	TodayPriceIndex = 0

	// ExpectedVolumeWeightedAveragePriceLength is the expected length of the ticker's
	// VolumeWeightedAveragePrice array.
	ExpectedVolumeWeightedAveragePriceLength = 2
)
