package kucoin

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	connectmath "github.com/skip-mev/connect/v2/pkg/math"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

type (
	// MessageType represents the type of message received from the KuCoin websocket.
	MessageType string

	// TopicType represents the type of topic to subscribe to i.e. spot for this
	// implementation.
	TopicType string

	// SubjectType represents the type of subject that was subscribed to i.e. ticker.
	SubjectType string
)

const (
	// WelcomeMessage represents the welcome message received when first connecting
	// to the websocket.
	//
	// ref: https://www.kucoin.com/docs/websocket/basic-info/create-connection
	WelcomeMessage MessageType = "welcome"

	// PingMessage represents a ping / heartbeat message that must be sent to the
	// websocket server every ping interval. The Ping interval is received and configured
	// when first connecting to the websocket.
	//
	// ref: https://www.kucoin.com/docs/websocket/basic-info/ping
	PingMessage MessageType = "ping"

	// PongMessage represents a pong / heartbeat message that is sent from the server
	// to the client in response to a ping message.
	//
	// ref: https://www.kucoin.com/docs/websocket/basic-info/ping
	PongMessage MessageType = "pong"

	// SubscribeMessage represents the subscribe message that must be sent to the
	// websocket server to subscribe to a channel.
	//
	// ref: https://www.kucoin.com/docs/websocket/basic-info/subscribe/introduction
	SubscribeMessage MessageType = "subscribe"

	// AckMessage represents the response message received from the websocket server
	// after sending a subscribe message.
	//
	// ref: https://www.kucoin.com/docs/websocket/basic-info/subscribe/introduction
	AckMessage MessageType = "ack"

	// Message represents a message received from the websocket server. This is returned
	// after subscribing to a channel and contains the payload for the desired data.
	//
	// ref: https://www.kucoin.com/docs/websocket/spot-trading/public-channels/ticker
	Message MessageType = "message"

	// TickerTopic represents the ticker topic. This will subscribe to the spot market
	// ticker for the specified trading pairs.
	//
	// ref: https://www.kucoin.com/docs/websocket/spot-trading/public-channels/ticker
	TickerTopic TopicType = "/market/ticker:"

	// TickerSubject represents the ticker subject. This should be returned in the
	// response message when subscribing to the ticker topic.
	//
	// ref: https://www.kucoin.com/docs/websocket/spot-trading/public-channels/ticker
	TickerSubject SubjectType = "trade.ticker"
)

// BaseMessage is utilized to determine the type of message that was received.
type BaseMessage struct {
	// ID is the ID of the message.
	ID string `json:"id"`

	// Type is the type of message.
	Type string `json:"type"`
}

// PingRequestMessage represents the ping message that must be sent to the
// websocket server every ping interval.
//
//	{
//		"id": "1545910590801",
//		"type": "ping"
//	}
//
// ref: https://www.kucoin.com/docs/websocket/basic-info/ping
type PingRequestMessage struct {
	BaseMessage
}

// NewHeartbeatMessage returns a new heartbeat message.
func NewHeartbeatMessage() ([]handlers.WebsocketEncodedMessage, error) {
	bz, err := json.Marshal(PingRequestMessage{
		BaseMessage: BaseMessage{
			ID:   fmt.Sprintf("%d", time.Now().Unix()),
			Type: string(PingMessage),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ping request message: %w", err)
	}

	return []handlers.WebsocketEncodedMessage{bz}, nil
}

// SubscribeRequestMessage represents the subscribe message that must be sent
// to the websocket server to subscribe to a channel.
//
// Spot Demo
//
//	{
//		"id": 1545910660739, // The id should be a unique value
//		"type": "subscribe",
//		"topic": "/market/ticker:BTC-USDT,ETH-USDT", // Topic needs to be subscribed. Some topics support to divisional subscribe the information of multiple trading pairs through ",".
//		"privateChannel": false, // Adopted the private channel or not. Set as false by default.
//		"response": true // Whether the server needs to return the receipt information of this subscription or not. Set as false by default.
//	}
//
// ref: https://www.kucoin.com/docs/websocket/basic-info/subscribe/introduction
type SubscribeRequestMessage struct {
	// ID is the ID of the message.
	ID int64 `json:"id"`

	// Type is the type of message.
	Type string `json:"type"`

	// Topic is the topic to subscribe to.
	Topic string `json:"topic"`

	// PrivateChannel is a flag that indicates if the channel is private.
	PrivateChannel bool `json:"privateChannel"`

	// Response is a flag that indicates if the server should return the receipt
	// information of this subscription.
	Response bool `json:"response"`
}

// NewSubscribeRequestMessage returns a new SubscribeRequestMessage.
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
		// Get the instruments for this batch.
		start := i * h.ws.MaxSubscriptionsPerBatch
		end := connectmath.Min((i+1)*h.ws.MaxSubscriptionsPerBatch, numInstruments)

		// Create the topic for this batch.
		topic := fmt.Sprintf("%s%s", TickerTopic, strings.Join(instruments[start:end], ","))
		bz, err := json.Marshal(SubscribeRequestMessage{
			ID:             time.Now().UTC().UnixNano(),
			Type:           string(SubscribeMessage),
			Topic:          topic,
			PrivateChannel: false,
			Response:       false,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal subscribe request message: %w", err)
		}
		msgs[i] = bz
	}

	return msgs, nil
}

// TickerResponseMessage represents the ticker response message received from
// the websocket server.
//
//	{
//		"type": "message",
//		"topic": "/market/ticker:BTC-USDT",
//		"subject": "trade.ticker",
//		"data": {
//	  		"sequence": "1545896668986", // Sequence number
//	  		"price": "0.08", // Last traded price
//	  		"size": "0.011", //  Last traded amount
//	  		"bestAsk": "0.08", // Best ask price
//	  		"bestAskSize": "0.18", // Best ask size
//	  		"bestBid": "0.049", // Best bid price
//	  		"bestBidSize": "0.036", // Best bid size
//	  		"Time": 1704873323416	//The matching time of the latest transaction
//		}
//	}
//
// ref: https://www.kucoin.com/docs/websocket/spot-trading/public-channels/ticker
type TickerResponseMessage struct {
	// Type is the type of message.
	Type string `json:"type"`

	// Topic is the topic of the message.
	Topic string `json:"topic"`

	// Subject is the subject of the message.
	Subject string `json:"subject"`

	// Data is the data of the message.
	Data TickerResponseMessageData `json:"data"`
}

// TickerResponseMessageData is the data field of the TickerResponseMessage.
type TickerResponseMessageData struct {
	// Sequence is the sequence number.
	Sequence string `json:"sequence"`

	// Price is the last traded price.
	Price string `json:"price"`
}

const (
	// ExpectedTopicLength is the expected length of the topic field in the
	// TickerResponseMessage.
	ExpectedTopicLength = 2

	// TickerIndex is the index of the ticker in the topic field of the
	// TickerResponseMessage.
	TickerIndex = 1
)
