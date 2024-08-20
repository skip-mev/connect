package bitfinex

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

type (
	// Event is the event type of message sent over the BitFinex websocket API.
	Event string

	// Channel is the channel of a message sent over the BitFinex websocket API.
	Channel string

	// ErrorCode is a type alias for an error code sent from the BitFinex websocket API.
	ErrorCode int64
)

const (
	// EventSubscribe indicates a subscribe action.
	EventSubscribe Event = "subscribe"
	// EventSubscribed indicates that a subscription was successful.
	EventSubscribed Event = "subscribed"
	// EventError indicates that an error occurred.
	EventError Event = "error"
	// ChannelTicker is the channel name for the ticker channel.
	ChannelTicker Channel = "ticker"
	// IDHeartbeat is the id always used for a heartbeat.
	IDHeartbeat = "hb"
	// ExpectedStreamPayloadLength is the expected length of the payload of a data stream.
	ExpectedStreamPayloadLength = 10
	// ExpectedBaseStreamLength is the expected length of a stream base message.
	ExpectedBaseStreamLength = 2

	// ErrorUnknownEvent indicates an unknown event.
	ErrorUnknownEvent ErrorCode = 10000
	// ErrorUnknownPair indicates an unknown pari.
	ErrorUnknownPair ErrorCode = 10001
	// ErrorLimitOpenChannels indicates the limit of open channels has been exceeded.
	ErrorLimitOpenChannels ErrorCode = 10305
	// ErrorSubscriptionFailed indicates a subscription failed.
	ErrorSubscriptionFailed ErrorCode = 10400
	// ErrorNotSubscribed indicates you are not subscribed to the given topic.
	ErrorNotSubscribed ErrorCode = 104001
)

// Error returns the error representation of the ErrorCode.
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
	msgs := make([]handlers.WebsocketEncodedMessage, len(symbols))

	if len(symbols) == 0 {
		return nil, fmt.Errorf("symbols cannot be empty")
	}

	for i, symbol := range symbols {
		msg, err := NewSubscribeMessage(symbol)
		if err != nil {
			return nil, fmt.Errorf("error marshalling subscription message: %w", err)
		}

		msgs[i] = msg
	}

	return msgs, nil
}

// NewSubscribeMessage creates a new subscribe message given the ticker symbol.
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
	ChannelID int    `json:"chanId" validate:"required"`
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
