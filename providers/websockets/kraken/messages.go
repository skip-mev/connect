package kraken

type (
	Event  string
	Status string
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
)

// SystemStatusMessage is the message that is sent to the client when the
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

// SubscibeRequestMessage is the message that is sent to the server to subscribe
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
