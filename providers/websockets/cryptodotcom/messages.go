package cryptodotcom

import (
	"encoding/json"
	"fmt"
	"math"

	connectmath "github.com/skip-mev/connect/v2/pkg/math"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

type (
	// Method is the method of the message received from the Crypto.com websocket API.
	Method string

	// StatusCode is the status code of the message received from the Crypto.com websocket API.
	StatusCode int64
)

const (
	// InstrumentMethod is the method used to subscribe to an instrument. This should
	// be called when the initial connection is established.
	InstrumentMethod Method = "subscribe"

	// HeartBeatRequestMethod is the method used to send a heartbeat message to the
	// Crypto.com websocket API. The heartbeat message is sent to the Crypto.com websocket
	// API every 30 seconds.
	HeartBeatRequestMethod Method = "public/heartbeat"

	// HeartBeatResponseMethod is the method used to respond to a heartbeat message
	// from the Crypto.com websocket API. Any heartbeat message received from the Crypto.com
	// websocket API must be responded to with a heartbeat response message with the same ID.
	HeartBeatResponseMethod Method = "public/respond-heartbeat"
)

const (
	// SuccessStatusCode is the status code of a successful message received from the
	// Crypto.com websocket API.
	SuccessStatusCode StatusCode = iota
)

const (
	// TickerChannel is the channel used to subscribe to the ticker of an instrument.
	TickerChannel string = "ticker.%s"
)

// HeartBeatResponseMessage is the response format that must be sent to the Crypto.com websocket API
// when a heartbeat message is received.
//
//	{
//		"id": 1587523073344,
//		"method": "public/respond-heartbeat"
//	}
type HeartBeatResponseMessage struct {
	// ID is the ID of the heartbeat message. This must be the same ID that is received
	// from the Crypto.com websocket API.
	ID int64 `json:"id"`

	// Method is the method of the heartbeat message.
	Method string `json:"method"`
}

// NewHeartBeatResponseMessage returns a new HeartBeatResponse message that can be sent to the
// Crypto.com websocket API.
func NewHeartBeatResponseMessage(id int64) ([]byte, error) {
	return json.Marshal(HeartBeatResponseMessage{
		ID:     id,
		Method: string(HeartBeatResponseMethod),
	})
}

// InstrumentRequestMessage is the request format that must be sent to the Crypto.com
// websocket API when subscribing to an instrument.
//
//	{
//			"id": 1,
//			"method": "subscribe",
//			"params": {
//		  		"channels": ["ticker.BTCUSD-PERP", "ticker.ETHUSD-PERP"]
//			},
//			"nonce": 1587523073344
//	 }
type InstrumentRequestMessage struct {
	// Method is the method of the subscribe message.
	Method string `json:"method"`

	// Params is the params of the subscribe message.
	Params InstrumentParams `json:"params"`
}

// InstrumentParams is the params of the subscribe message.
type InstrumentParams struct {
	// Channels is the channels that we want to subscribe to.
	Channels []string `json:"channels"`
}

// NewInstrumentMessage returns a new InstrumentRequestMessage that can be sent to
// the Crypto.com websocket API.
func (h *WebSocketHandler) NewInstrumentMessage(instruments []string) ([]handlers.WebsocketEncodedMessage, error) {
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

		bz, err := json.Marshal(InstrumentRequestMessage{
			Method: string(InstrumentMethod),
			Params: InstrumentParams{
				Channels: instruments[start:end],
			},
		})
		if err != nil {
			return msgs, err
		}
		msgs[i] = bz
	}

	return msgs, nil
}

// InstrumentResponseMessage is the response received from the Crypto.com websocket API when subscribing
// to an instrument i.e. price feed.
//
//	{
//		"id": -1,
//		"method": "subscribe",
//		"code": 0,
//		"result": {
//		  "instrument_name": "BTCUSD-PERP",
//		  "subscription": "ticker.BTCUSD-PERP",
//		  "channel": "ticker",
//		  "data": [{
//			"h": "51790.00",        // Price of the 24h highest trade
//			"l": "47895.50",        // Price of the 24h lowest trade, null if there weren't any trades
//			"a": "51174.500000",    // The price of the latest trade, null if there weren't any trades
//			"c": "0.03955106",      // 24-hour price change, null if there weren't any trades
//			"b": "51170.000000",    // The current best bid price, null if there aren't any bids
//			"bs": "0.1000",         // The current best bid size, null if there aren't any bids
//			"k": "51180.000000",    // The current best ask price, null if there aren't any asks
//			"ks": "0.2000",         // The current best ask size, null if there aren't any bids
//			"i": "BTCUSD-PERP",     // Instrument name
//			"v": "879.5024",        // The total 24h traded volume
//			"vv": "26370000.12",    // The total 24h traded volume value (in USD)
//			"oi": "12345.12",       // Open interest
//			"t": 1613580710768
//		  }]
//		}
//	}
type InstrumentResponseMessage struct {
	// ID is the ID of the subscribe message.
	ID int64 `json:"id"`

	// Method is the method of the subscribe message.
	Method string `json:"method"`

	// Result is the result of the subscribe message.
	Result InstrumentResult `json:"result"`

	// Code is the status of the response. 0 means success.
	Code int64 `json:"code"`
}

// InstrumentResult is the result of the subscribe message.
type InstrumentResult struct {
	// Data is the data of the subscribe message.
	Data []InstrumentData `json:"data"`
}

// InstrumentData is the data of the subscribe message.
type InstrumentData struct {
	// LatestTradePrice is the price of the latest trade, null if there weren't any trades.
	LatestTradePrice string `json:"a"`

	// Name is the instrument name.
	Name string `json:"i"`
}
