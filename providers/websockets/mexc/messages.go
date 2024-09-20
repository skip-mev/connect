package mexc

import (
	"encoding/json"
	"fmt"
	"math"

	connectmath "github.com/skip-mev/connect/v2/pkg/math"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

type (
	// MethodType defines the type of message that is being sent to the MEXC websocket.
	MethodType string

	// ChannelType defines the type of channel that the client is subscribing to.
	ChannelType string
)

const (
	// SubscriptionMethod is the method that is sent to the MEXC websocket to subscribe to a
	// currency pair i.e. market.
	//
	// ref: https://mexcdevelop.github.io/apidocs/spot_v3_en/#live-subscribing-unsubscribing-to-streams
	SubscriptionMethod MethodType = "SUBSCRIPTION"

	// PingMethod is the method that is sent to the MEXC websocket to ping the server. This should
	// be done every 30 seconds.
	//
	// ref: https://mexcdevelop.github.io/apidocs/spot_v3_en/#live-subscribing-unsubscribing-to-streams
	PingMethod MethodType = "PING"

	// PongMethod is the method that is sent from the server to the client to confirm that the
	// client has successfully pinged the server.
	//
	// ref: https://mexcdevelop.github.io/apidocs/spot_v3_en/#live-subscribing-unsubscribing-to-streams
	PongMethod MethodType = "PONG"

	// MiniTickerChannel is the channel that is used to subscribe to the mini ticker data for a
	// currency pair i.e. market.
	//
	// ex: spot@public.miniTicker.v3.api@BTCUSDT@UTC+8
	//
	// ref: https://mexcdevelop.github.io/apidocs/spot_v3_en/#miniticker
	MiniTickerChannel ChannelType = "spot@public.miniTicker.v3.api@"
)

// BaseMessage defines the base message that is used to determine the type of message that is being
// sent to the MEXC websocket.
type BaseMessage struct {
	// ID is the ID of the subscription request.
	ID int64 `json:"id"`

	// Code is the status code of the subscription request.
	Code int64 `json:"code"`

	// Message is the message that is sent from the MEXC websocket to confirm that
	// the client has successfully subscribed to a currency pair i.e. market.
	Message string `json:"msg"`
}

// IsEmpty returns true if no data has been set for the message.
func (m *BaseMessage) IsEmpty() bool {
	return m.ID == 0 && m.Code == 0 && len(m.Message) == 0
}

// SubscriptionRequestMessage defines the message that is sent to the MEXC websocket to subscribe to a
// currency pair i.e. market.
//
//	{
//		"method": "SUBSCRIPTION",
//		"params": [
//			"spot@public.deals.v3.api@BTCUSDT"
//		]
//	}
//
// ref: https://mexcdevelop.github.io/apidocs/spot_v3_en/#websocket-market-streams
type SubscriptionRequestMessage struct {
	// Method is the method that is being sent to the MEXC websocket.
	Method string `json:"method"`

	// Params is the list of channels that the client is subscribing to. This
	// list cannot exceed 30 channels.
	Params []string `json:"params"`
}

// NewSubscribeRequestMessage returns a new SubscriptionRequestMessage.
func (h *WebSocketHandler) NewSubscribeRequestMessage(instruments []string) ([]handlers.WebsocketEncodedMessage, error) {
	numInstruments := len(instruments)
	if numInstruments == 0 {
		return nil, fmt.Errorf("cannot subscribe to 0 instruments")
	}

	numBatches := int(math.Ceil(float64(numInstruments) / float64(h.ws.MaxSubscriptionsPerBatch)))
	msgs := make([]handlers.WebsocketEncodedMessage, numBatches)
	for i := 0; i < numBatches; i++ {
		// Get the instruments for this batch.
		start := i * h.ws.MaxSubscriptionsPerBatch
		end := connectmath.Min((i+1)*h.ws.MaxSubscriptionsPerBatch, numInstruments)

		bz, err := json.Marshal(SubscriptionRequestMessage{
			Method: string(SubscriptionMethod),
			Params: instruments[start:end],
		})
		if err != nil {
			return nil, err
		}
		msgs[i] = bz
	}

	return msgs, nil
}

// SubscriptionResponseMessage defines the message that is sent from the MEXC websocket to confirm that
// the client has successfully subscribed to a currency pair i.e. market.
//
//	{
//		"id":0,
//		"code":0,
//		"msg":"spot@public.deals.v3.api@BTCUSDT"
//	}
//
// ref: https://mexcdevelop.github.io/apidocs/spot_v3_en/#live-subscribing-unsubscribing-to-streams
type SubscriptionResponseMessage struct {
	BaseMessage
}

// PingRequestMessage defines the message that is sent to the MEXC websocket to ping the server.
//
//	{
//		"method":"PING"
//	}
//
// ref: https://mexcdevelop.github.io/apidocs/spot_v3_en/#live-subscribing-unsubscribing-to-streams
type PingRequestMessage struct {
	BaseMessage
}

// NewPingRequestMessage returns a new PingRequestMessage.
func NewPingRequestMessage() ([]handlers.WebsocketEncodedMessage, error) {
	bz, err := json.Marshal(PingRequestMessage{
		BaseMessage: BaseMessage{
			Message: string(PingMethod),
		},
	})
	if err != nil {
		return nil, err
	}

	return []handlers.WebsocketEncodedMessage{bz}, nil
}

// PongResponseMessage defines the message that is sent from the server to the client to confirm that the
// client has successfully pinged the server.
//
//	{
//		"id":0,
//		"code":0,
//		"msg":"PONG"
//	}
//
// ref: https://mexcdevelop.github.io/apidocs/spot_v3_en/#live-subscribing-unsubscribing-to-streams
type PongResponseMessage struct {
	BaseMessage
}

// TickerResponseMessage defines the message that is sent from the MEXC websocket to provide the latest
// ticker data for a currency pair i.e. market.
//
// c		string	channel name
// d		data	data
// >s		string	symbol
// >p		string	deal price
// >r		string	price Change Percent in utc8
// >tr		string	price Change Percent in time zone
// >h		string	24h high price
// >l		string	24h low price
// >v		string	24h volume
// >q		string	24h quote Volume
// >lastRT	string	etf Last Rebase Time
// >MT		string	etf Merge Times
// >NV		string	etf Net Value
//
//	{
//		"d": {
//		  "s":"BTCUSDT",
//		  "p":"36474.74",
//		  "r":"0.0354",
//		  "tr":"0.0354",
//		  "h":"36549.72",
//		  "l":"35101.68",
//		  "v":"375173478.65",
//		  "q":"10557.72895",
//		  "lastRT":"-1",
//		  "MT":"0",
//		  "NV":"--",
//		  "t":"1699502456050"
//		},
//		"c":"spot@public.miniTicker.v3.api@BTCUSDT@UTC+8",
//		"t":1699502456051,
//		"s":"BTCUSDT"
//	}
//
// ref: https://mexcdevelop.github.io/apidocs/spot_v3_en/#miniticker
type TickerResponseMessage struct {
	// Channel is the channel that the client is subscribed to.
	Channel string `json:"c"`

	// Data is the latest ticker data for a currency pair i.e. market.
	Data TickerData `json:"d"`
}

// TickerData defines the latest ticker data for a currency pair i.e. market.
type TickerData struct {
	// Symbol is the currency pair i.e. market that the ticker data is for.
	Symbol string `json:"s"`

	// Price is the latest price for the currency pair i.e. market.
	Price string `json:"p"`
}
