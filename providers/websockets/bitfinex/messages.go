package bitfinex

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
)

type (
	// Event is the event type of a message sent over the Bitfinex websocket API.
	Event string

	// Channel is the channel of a message sent over the Bitfinex websocket API.
	Channel string

	// ChannelID is the id of a channel being communicated over the BitFinex websocket API.
	ChannelID string
)

const (
	// EventSubscribe indicates a subscribe action.
	EventSubscribe Event = "subscribe"
	// EventSubscribed indicates that a subscription was successful.
	EventSubscribed Event = "subscribed"
	// EventError indicates that an error occurred.
	EventError Event = "error"
	// EventNil represents an empty event field.
	EventNil Event = ""
	// ChannelTicker is the channel name for the ticker channel.
	ChannelTicker Channel = "ticker"
	// ChannelIdHeartbeat is the channelID always used for a heartbeat.
	ChannelIdHeartbeat ChannelID = "hb"
)

type ErrorCode int64

const (
	ErrorUnknownEvent       ErrorCode = 10000
	ErrorUnknownPair        ErrorCode = 10001
	ErrorLimitOpenChannels  ErrorCode = 10305
	ErrorSubscriptionFailed ErrorCode = 10400
	ErrorNotSubscribed      ErrorCode = 104001
)

// Error returns the string representation of the ErrorCode.
func (e ErrorCode) Error() error {
	switch e {
	case ErrorUnknownEvent:
		return errors.New("unknown event")
	case ErrorUnknownPair:
		return errors.New("unknown pair")
	case ErrorLimitOpenChannels:
		return errors.New("limit of open channels reached")
	case ErrorSubscriptionFailed:
		return errors.New("subscribed failed")
	case ErrorNotSubscribed:
		return errors.New("not subscribed")
	default:
		return errors.New("unknown error type")

	}
}

// BaseMessage is the base message structure for subscription requests and responses in the
// BitFinex websocket API.
type BaseMessage struct {
	Event string `json:"event" validate:"required"`
}

// SubscribeMessage is a base message used to make a subscription request.
//
// Ex:
//
//	{
//	 event: "subscribe",
//	 channel: "ticker",
//	 symbol: SYMBOL
//	}
//
// ref: https://docs.bitfinex.com/reference/ws-public-ticker
type SubscribeMessage struct {
	BaseMessage
	Channel string `json:"channel" validate:"required"`
	Symbol  string `json:"symbol" validate:"required"`
}

// NewSubscribeMessages creates subscription messages for the given tickers.
func NewSubscribeMessages(symbols []string) ([]handlers.WebsocketEncodedMessage, error) {
	var msgs []handlers.WebsocketEncodedMessage

	if len(symbols) == 0 {
		return nil, fmt.Errorf("symbols cannot be empty")
	}

	for _, symbol := range symbols {
		msg, err := NewSubscribeMessage(symbol)
		if err != nil {
			return nil, fmt.Errorf("error marshalling subscription message: %w", err)
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

// NewSubscribeMessage creates a new subscribe message given the ticker symbol
func NewSubscribeMessage(symbol string) (handlers.WebsocketEncodedMessage, error) {
	return json.Marshal(SubscribeMessage{
		BaseMessage: BaseMessage{
			Event: string(EventSubscribe),
		},
		Channel: string(ChannelTicker),
		Symbol:  symbol,
	})
}

// SubscribedMessage is message indicating the status of a subscription request.
//
// Ex:
//
//	{
//	  event: "subscribed",
//	  channel: "ticker",
//	  chanId: CHANNEL_ID,
//	  symbol: SYMBOL,
//	  pair: PAIR
//	}
//
// ref: https://docs.bitfinex.com/reference/ws-public-ticker
type SubscribedMessage struct {
	BaseMessage
	Channel   string `json:"channel" validate:"required"`
	ChannelID string `json:"chanId" validate:"required"`
	Pair      string `json:"pair" validate:"required"`
}

// ErrorMessage represent an error returned by the peer.
//
// Ex.
//
//	{
//	  "event": "error",
//	  "msg": ERROR_MSG,
//	  "code": ERROR_CODE
//	}
//
// ref: https://docs.bitfinex.com/reference/ws-public-ticker
type ErrorMessage struct {
	BaseMessage
	Msg  string `json:"msg" validate:"required"`
	Code int64  `json:"code" validate:"required"`
}

// BaseStreamMessage is the base message for all stream messages received from a peer.
type BaseStreamMessage struct {
	ChannelID string `json:"chanId" validate:"required"`
}

// HeartbeatStream is the heartbeat message sent from the server every 15 seconds.
type HeartbeatStream BaseStreamMessage

// TickerStream is a stream message continually received after successfully
// making a subscription.
//
// Ex:
//
//	{
//	  chanId: CHANNEL_ID,
//	  lastPrice: LAST_PRICE,
//	}
//
// ref: https://docs.bitfinex.com/reference/ws-public-ticker
type TickerStream struct {
	BaseStreamMessage
	LastPrice string `json:"lastPrice" validate:"required"`
}
