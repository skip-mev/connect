package coinbase

import (
	"encoding/json"
	"fmt"
	"math"

	connectmath "github.com/skip-mev/connect/v2/pkg/math"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

type (
	// MessageType represents the type of message that is sent/received from
	// the websocket feed.
	MessageType string

	// ChannelType represents the type of channel that is sent/received from
	// the websocket feed.
	ChannelType string
)

const (
	// SubscribeMessage represents a subscribe message. This must be sent as the
	// first message when connecting to the websocket feed and must be sent within
	// 5 seconds of the initial connection. Otherwise, the connection will be
	// closed.
	//
	// ref: https://docs.cloud.coinbase.com/exchange/docs/websocket-overview#subscribe
	SubscribeMessage MessageType = "subscribe"

	// SubscriptionsMessage represents a subscriptions message. This is sent by the
	// websocket feed after a subscribe message is sent.
	//
	// ref: https://docs.cloud.coinbase.com/exchange/docs/websocket-overview#subscribe
	SubscriptionsMessage MessageType = "subscriptions"

	// TickerMessage represents a ticker message. This is sent by the websocket feed
	// when a match happens.
	//
	// ref: https://docs.cloud.coinbase.com/exchange/docs/websocket-channels#ticker-channel
	TickerMessage MessageType = "ticker"

	// HeartbeatMessage represents a heartbeat message. This is sent by the websocket feed
	// every second for the subscribed channels.
	//
	// ref: https://docs.cdp.coinbase.com/exchange/docs/websocket-channels/#heartbeat-channel
	HeartbeatMessage MessageType = "heartbeat"
)

const (
	// TickerChannel represents the ticker channel. The ticker channel providers real-time price
	// updates every time a match happens. It batches updates in case of cascading matches,
	// greatly reducing bandwidth requirements.
	//
	// ref: https://docs.cloud.coinbase.com/exchange/docs/websocket-channels#ticker-channel
	TickerChannel ChannelType = "ticker"

	// HeartbeatChannel represents the heartbeat channel. The heartbeat channel provides a
	// heartbeat every second for the subscribed channels. Heartbeats include sequence numbers
	// and last trade IDs that can be used to verify that no messages were missed.
	//
	// ref: https://docs.cdp.coinbase.com/exchange/docs/websocket-channels/#heartbeat-channel
	HeartbeatChannel ChannelType = "heartbeat"
)

// BaseMessage represents a base message. This is used to determine the type of message
// that was received.
type BaseMessage struct {
	// Type is the type of message.
	Type string `json:"type"`
}

// SubscribeRequestMessage represents a subscribe request message.
//
// Request
// Subscribe to ETH-USD and ETH-EUR with the level2, heartbeat and ticker channels,
// plus receive the ticker entries for ETH-BTC and ETH-USD
//
//	{
//	    "type": "subscribe",
//	    "product_ids": [
//	        "ETH-USD",
//	        "ETH-EUR"
//	    ],
//	    "channels": [
//	        "level2",
//	        "heartbeat",
//	        {
//	            "name": "ticker",
//	            "product_ids": [
//	                "ETH-BTC",
//	                "ETH-USD"
//	            ]
//	        }
//	    ]
//	}
//
// Request version 2 <--- This is the one we use.
//
//	{
//	    "type": "unsubscribe",
//	    "product_ids": [
//	        "ETH-USD",
//	        "ETH-EUR"
//	    ],
//	    "channels": [
//	        "ticker"
//	    ]
//	}
//
// ref: https://docs.cloud.coinbase.com/exchange/docs/websocket-overview#subscribe
type SubscribeRequestMessage struct {
	// Type is the type of message.
	Type string `json:"type"`

	// ProductIDs is a list of product IDs (markets) to subscribe to.
	ProductIDs []string `json:"product_ids"`

	// Channels is a list of channels to subscribe to.
	Channels []string `json:"channels"`
}

// NewSubscribeRequestMessage returns a new subscribe request message.
func (h *WebSocketHandler) NewSubscribeRequestMessage(instruments []string) ([]handlers.WebsocketEncodedMessage, error) {
	numInstruments := len(instruments)
	if numInstruments == 0 {
		return nil, fmt.Errorf("no instruments provided")
	}

	numBatches := int(math.Ceil(float64(numInstruments) / float64(h.ws.MaxSubscriptionsPerBatch)))
	msgs := make([]handlers.WebsocketEncodedMessage, numBatches)
	for i := 0; i < numBatches; i++ {
		// Get the instruments for the batch.
		start := i * h.ws.MaxSubscriptionsPerBatch
		end := connectmath.Min((i+1)*h.ws.MaxSubscriptionsPerBatch, numInstruments)

		bz, err := json.Marshal(SubscribeRequestMessage{
			Type:       string(SubscribeMessage),
			ProductIDs: instruments[start:end],
			Channels:   []string{string(TickerChannel), string(HeartbeatChannel)},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal subscribe request message %w", err)
		}
		msgs[i] = bz
	}

	return msgs, nil
}

// SubscribeResponseMessage represents a subscribe response message.
//
// Response
//
//	{
//	    "type": "subscriptions",
//	    "channels": [
//	        {
//	            "name": "level2",
//	            "product_ids": [
//	                "ETH-USD",
//	                "ETH-EUR"
//	            ],
//	        },
//	        {
//	            "name": "heartbeat",
//	            "product_ids": [
//	                "ETH-USD",
//	                "ETH-EUR"
//	            ],
//	        },
//	        {
//	            "name": "ticker",
//	            "product_ids": [
//	                "ETH-USD",
//	                "ETH-EUR",
//	                "ETH-BTC"
//	            ]
//	        }
//	    ]
//	}
//
// ref: https://docs.cloud.coinbase.com/exchange/docs/websocket-overview#subscriptions-message
type SubscribeResponseMessage struct {
	// Type is the type of message.
	Type string `json:"type"`

	// Channels is a list of channels that were subscribed to.
	Channels []Channel `json:"channels"`
}

// Channel represents a channel that was subscribed to.
type Channel struct {
	// Name is the name of the channel.
	Name string `json:"name"`

	// Instruments is a list of product IDs (markets) that were subscribed to.
	Instruments []string `json:"product_ids"`
}

// TickerResponseMessage represents a ticker response message.
//
// Response
//
//	{
//			"type": "ticker",
//			"sequence": 37475248783,
//			"product_id": "ETH-USD",
//			"price": "1285.22",
//			"open_24h": "1310.79",
//			"volume_24h": "245532.79269678",
//			"low_24h": "1280.52",
//			"high_24h": "1313.8",
//			"volume_30d": "9788783.60117027",
//			"best_bid": "1285.04",
//			"best_bid_size": "0.46688654",
//			"best_ask": "1285.27",
//			"best_ask_size": "1.56637040",
//			"side": "buy",
//			"time": "2022-10-19T23:28:22.061769Z",
//			"trade_id": 370843401,
//			"last_size": "11.4396987"
//	}
//
// ref: https://docs.cloud.coinbase.com/exchange/docs/websocket-channels#ticker-channel
type TickerResponseMessage struct {
	// Type is the type of message.
	Type string `json:"type"`

	// Sequence is the sequence number of the message.
	Sequence int64 `json:"sequence"`

	// Ticker is the product ID of the ticker.
	Ticker string `json:"product_id"`

	// Price is the price of the ticker.
	Price string `json:"price"`

	// TradeID is the trade ID of the ticker.
	TradeID int64 `json:"trade_id"`
}

// HeartbeatResponseMessage represents a heartbeat response message.
//
// Response
//
//	{
//				"type": "heartbeat",
//				"sequence": 90,
//				"last_trade_id": 20,
//				"product_id": "BTC-USD",
//				"time": "2014-11-07T08:19:28.464459Z"
//	}
//
// ref: https://docs.cdp.coinbase.com/exchange/docs/websocket-channels/#heartbeat-channel
type HeartbeatResponseMessage struct {
	// Type is the type of message.
	Type string `json:"type"`

	// Sequence is the sequence number of the message.
	Sequence int64 `json:"sequence"`

	// LastTradeID is the last trade ID of the message.
	LastTradeID int64 `json:"last_trade_id"`

	// Ticker is the product ID of the ticker.
	Ticker string `json:"product_id"`
}
