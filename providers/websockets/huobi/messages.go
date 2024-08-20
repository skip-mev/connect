package huobi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

// Status is the status of a subscription request.
type Status string

const (
	marketTickerFormatString = "market.%s.ticker"

	// StatusOk is a status indicating a request was successful.
	StatusOk Status = "ok"
)

// subFromSymbol creates the proper subscription string to a ticker for the given symbol.
func subFromSymbol(symbol string) string {
	return fmt.Sprintf(marketTickerFormatString, symbol)
}

// symbolFromSub returns the symbol from a sub topic.  If the sub message is incorrectly formatted,
// and empty string is returned.
func symbolFromSub(sub string) string {
	parts := strings.Split(sub, ".")
	if len(parts) == 3 {
		return parts[1]
	}

	return ""
}

// PingMessage is a message representing a ping.  It must have the same value as the
// corresponding pong.
type PingMessage struct {
	Ping int64 `json:"ping"`
}

// PongMessage is a message representing a pong.  It must have the same value as the
// corresponding ping.
type PongMessage struct {
	Pong int64 `json:"pong"`
}

// NewPongMessage returns a PongMessage for the corresponding PingMessage.
func NewPongMessage(
	message PingMessage,
) ([]handlers.WebsocketEncodedMessage, error) {
	bz, err := json.Marshal(PongMessage{Pong: message.Ping})
	return []handlers.WebsocketEncodedMessage{bz}, err
}

// SubscriptionRequest is the request message sent to the server to subscribe to a topic.
type SubscriptionRequest struct {
	Sub string `json:"sub"`
	ID  string `json:"id"`
}

// NewSubscriptionRequest creates a new encoded SubscriptionRequest from the given symbol.
func NewSubscriptionRequest(symbol string) (handlers.WebsocketEncodedMessage, error) {
	return json.Marshal(SubscriptionRequest{
		Sub: subFromSymbol(symbol),
		ID:  symbol,
	})
}

// SubscriptionResponse is the response message sent from the server after a subscription request.
type SubscriptionResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Subbed string `json:"subbed"`
}

// TickerStream is the stream for a given ticker sent every 100ms by the Huobi API.
type TickerStream struct {
	Channel string `json:"ch"`
	Tick    Tick   `json:"tick"`
}

// Tick is the tick payload attached to a TickerStream message.
type Tick struct {
	LastPrice float64 `json:"lastPrice"`
}
