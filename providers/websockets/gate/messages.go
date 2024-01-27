package gate

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	"time"
)

type (
	// ErrorCode is a type alias for an int error code.
	ErrorCode int
	// Channel is a type alias for a channel identifier.
	Channel string
	// Event is a type alias for an event identifier.
	Event string
)

const (
	// ChannelTickers is the tickers channel to subscribe to.
	ChannelTickers Channel = "spot.tickers"

	// EventSubscribe is the event for subscribing to a topic.
	EventSubscribe Event = "subscribe"

	// EventUpdate is the event indicating an update.
	EventUpdate Event = "update"

	// ErrorInvalidRequestBody is returned for an invalid body in the request.
	ErrorInvalidRequestBody ErrorCode = 1
	// ErrorInvalidArgument is returned for an invalid argument in the request.
	ErrorInvalidArgument ErrorCode = 2
	// ErrorServer is returned when there is a server side error.
	ErrorServer ErrorCode = 3
)

// Error returns the error representation of the ErrorCode.
func (e ErrorCode) Error() error {
	switch e {
	case ErrorInvalidRequestBody:
		return errors.New("invalid body in request")
	case ErrorInvalidArgument:
		return errors.New("invalid argument in request")
	case ErrorServer:
		return errors.New("server side error")
	default:
		return errors.New("unknown error")
	}
}

// BaseMessage is a base message for a request/response from a peer.
type BaseMessage struct {
	// Time is the time of the message.
	Time int `json:"time"`
	// ID is the optional ID for the message.
	ID int `json:"id"`
	// Channel is the channel to subscribe to.
	Channel string `json:"channel"`
	// Event is the event the request/response is taking.
	Event string `json:"event"`
}

// ErrorMessage represents an error returned from the Gate.io websocket API.
type ErrorMessage struct {
	// Code is the integer representation of the error code.
	Code int `json:"code"`
	// Message is the accompanying error message.
	Message string `json:"message"`
}

// RequestResult is the result message returned in a response from the Gate.io websocket API.
type RequestResult struct {
	// Status is the status of the result.
	Status string `json:"status"`
}

// SubscribeRequest is a subscription request sent to the Gate.io websocket API.
type SubscribeRequest struct {
	BaseMessage
	// Payload is the argument payload sent for the corresponding request.
	Payload []string `json:"payload"`
}

// NewSubscribeRequest returns a new SubscribeRequest encoded message for the given symbols
func NewSubscribeRequest(symbols []string) ([]handlers.WebsocketEncodedMessage, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("cannot attach payload of 0 length")
	}

	bz, err := json.Marshal(SubscribeRequest{
		BaseMessage: BaseMessage{
			Time:    time.Now().UTC().Second(),
			ID:      time.Now().UTC().Second(),
			Channel: string(ChannelTickers),
			Event:   string(EventSubscribe),
		},
		Payload: symbols,
	})

	return []handlers.WebsocketEncodedMessage{bz}, err
}

// SubscribeResponse is a subscription response sent from the Gate.io websocket API.
type SubscribeResponse struct {
	BaseMessage
	// Error is the error message returned.  Will be empty if no error is returned.
	Error ErrorMessage `json:"error"`
	// Result is the result returned from the server.
	Result RequestResult `json:"result"`
}

// TickerStream is the data stream returned for a ticker subscription.
type TickerStream struct {
	// Time is the time of the message.
	Time int `json:"time"`
	// Channel is the channel to subscribe to.
	Channel string `json:"channel"`
	// Event is the event the request/response is taking.
	Event string `json:"event"`
	// Result is the result body of the data stream.
	Result TickerResult `json:"result"`
}

// TickerResult is the result returned in a TickerStream message.
type TickerResult struct {
	// CurrencyPair is the currency pair for the given data stream.
	CurrencyPair string `json:"currency_pair"`
	// Last is the last price of the pair.
	Last string `json:"last"`
}
