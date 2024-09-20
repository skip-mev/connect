package binance

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	connectmath "github.com/skip-mev/connect/v2/pkg/math"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

type (
	// MethodType represents the type of message that is sent to the websocket feed.
	MethodType string
	// StreamType represents the type of stream that is sent/received from the websocket
	// feed. Streams are used to determine what instrument you are subscribing to / how to
	// handle the message.
	StreamType string
)

const (
	// SubscribeMethod represents a subscribe method. This must be sent as the first message
	// when connecting to the websocket feed.
	//
	// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#subscribe-to-a-stream
	SubscribeMethod MethodType = "SUBSCRIBE"

	// AggregateTradeStream represents the aggregate trade stream. This stream provides
	// trade information that is aggregated for a single taker order.
	//
	// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#aggregate-trade-streams
	AggregateTradeStream StreamType = "aggTrade"

	// TickerStream represents the ticker stream. This provides a 24hr rolling window ticker statistics for a single
	// symbol. These are NOT the statistics of the UTC day, but a 24hr rolling window for the previous 24hrs.
	//
	// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#individual-symbol-ticker-streams
	TickerStream StreamType = "ticker"

	// Separator is the separator used to separate the instrument and the stream type.
	Separator = "@"
)

// SubscribeMessageRequest represents a subscribe message request. This is used to subscribe
// to Binance websocket streams.
//
// Request
//
//	{
//	  "method": "SUBSCRIBE",
//	  "params": [
//	    "btcusdt@aggTrade",
//	    "btcusdt@depth"
//	  ],
//	  "id": 1
//	}
//
// The ID field is used to uniquely identify the messages going back and forth. By default, the ID
// is randomly generated with a minimum value of 1.
//
// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#live-subscribingunsubscribing-to-streams
type SubscribeMessageRequest struct {
	// Method is the method type for the message.
	Method string `json:"method"`
	// Params is the list of streams to subscribe to.
	Params []string `json:"params"`
	// ID is the unique identifier for the message.
	ID int64 `json:"id"`
}

// SubscribeMessageResponse represents a subscribe message response. This is used to determine
// whether the subscription was (un)successful.
//
// Response
//
//	{
//			"result": null,
//			"id": 1
//	}
//
// The ID field is used to uniquely identify the messages going back and forth, the same one sent in
// the initial subscription. The result is null if the subscription was successful.
//
// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#live-subscribingunsubscribing-to-streams
type SubscribeMessageResponse struct {
	// Result is the result of the subscription.
	Result interface{} `json:"result"`
	// ID is the unique identifier for the message.
	ID int64 `json:"id"`
}

// IsEmpty returns true if no data has been set for the message.
func (m *SubscribeMessageResponse) IsEmpty() bool {
	return m.ID == 0 && m.Result == nil
}

// StreamMessageResponse represents a stream message response. This is used to represent the
// data that is received from the Binance websocket. All stream data will have a stream and
// data field.
//
// # Response
//
// {
// 	"stream": "btcusdt@aggTrade",
// 	"data": {
// 	  	"e": "aggTrade",
// 	  	"E": 1716915868145,
// 	  	"s": "BTCUSDT",
// 	  	"a": 3020757327,
// 	  	"p": "67734.00000000",
// 	  	"q": "0.00230000",
// 	  	"f": 3617006743,
// 	  	"l": 3617006743,
// 	  	"T": 1716915868145,
// 	  	"m": false,
// 	  	"M": true
// 		}
//  }

// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#aggregate-trade-streams
type StreamMessageResponse struct {
	// Stream is the stream type.
	Stream string `json:"stream"`
}

// GetStreamType returns the stream type from the stream message response.
func (m *StreamMessageResponse) GetStreamType() StreamType {
	stream := strings.Split(m.Stream, Separator)
	if len(stream) != 2 {
		return ""
	}
	return StreamType(stream[1])
}

// AggregatedTradeMessageResponse represents an aggregated trade message response. This is used to
// represent the aggregated trade data that is received from the Binance websocket.
//
// # Response
//
//	{
//	  	"e": "aggTrade",    // Event type
//	  	"E": 1672515782136, // Event time
//	  	"s": "BNBBTC",      // Symbol
//	  	"a": 12345,         // Aggregate trade ID
//	  	"p": "0.001",       // Price
//	  	"q": "100",         // Quantity
//	  	"f": 100,           // First trade ID
//	  	"l": 105,           // Last trade ID
//	  	"T": 1672515782136, // Trade time
//	  	"m": true,          // Is the buyer the market maker?
//	 	"M": true           // Ignore
//	}
//
// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#aggregate-trade-streams
type AggregatedTradeMessageResponse struct {
	Data struct {
		// Ticker is the symbol.
		Ticker string `json:"s"`
		// Price is the price.
		Price string `json:"p"`
	} `json:"data"`
}

// TickerMessageResponse represents a ticker message response. This is used to represent the
// ticker data that is received from the Binance websocket.
//
// # Response
//
//	{
//			"e": "24hrTicker",  // Event type
//			"E": 1672515782136, // Event time
//			"s": "BNBBTC",      // Symbol
//			"p": "0.0015",      // Price change
//			"P": "250.00",      // Price change percent
//			"w": "0.0018",      // Weighted average price
//			"x": "0.0009",      // First trade(F)-1 price (first trade before the 24hr rolling window)
//			"c": "0.0025",      // Last price
//			"Q": "10",          // Last quantity
//			"b": "0.0024",      // Best bid price
//			"B": "10",          // Best bid quantity
//			"a": "0.0026",      // Best ask price
//			"A": "100",         // Best ask quantity
//			"o": "0.0010",      // Open price
//			"h": "0.0025",      // High price
//			"l": "0.0010",      // Low price
//			"v": "10000",       // Total traded base asset volume
//			"q": "18",          // Total traded quote asset volume
//			"O": 0,             // Statistics open time
//			"C": 86400000,      // Statistics close time
//			"F": 0,             // First trade ID
//			"L": 18150,         // Last trade Id
//			"n": 18151          // Total number of trades
//	}
//
// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#individual-symbol-ticker-streams
type TickerMessageResponse struct {
	Data struct {
		// Ticker is the symbol.
		Ticker string `json:"s"`
		// LastPrice is the last price.
		LastPrice string `json:"c"`
		// StatisticsCloseTime is the statistics close time.
		//
		// Note: This is unused but is included since json.Unmarshal requires all fields with same character but different casing
		// to be present.
		StatisticsCloseTime int64 `json:"C"`
	} `json:"data"`
}

// NewSubscribeRequestMessage returns a set of messages to subscribe to the Binance websocket. This will
// subscribe each instrument to the aggregate trade and ticker streams.
func (h *WebSocketHandler) NewSubscribeRequestMessage(instruments []string) ([]handlers.WebsocketEncodedMessage, error) {
	numInstruments := len(instruments)
	if numInstruments == 0 {
		return nil, fmt.Errorf("no instruments to subscribe to")
	}

	numBatches := int(math.Ceil(float64(numInstruments) / float64(h.ws.MaxSubscriptionsPerBatch)))
	msgs := make([]handlers.WebsocketEncodedMessage, numBatches)
	for i := 0; i < numBatches; i++ {
		// Get the instruments for the batch.
		start := i * h.ws.MaxSubscriptionsPerBatch
		end := connectmath.Min((i+1)*h.ws.MaxSubscriptionsPerBatch, numInstruments)
		batch := instruments[start:end]

		// Create the subscriptions for the instruments.
		params := make([]string, 0)
		for _, instrument := range batch {
			params = append(params, fmt.Sprintf("%s%s%s", strings.ToLower(instrument), Separator, string(AggregateTradeStream)))
			params = append(params, fmt.Sprintf("%s%s%s", strings.ToLower(instrument), Separator, string(TickerStream)))
		}

		// Generate a random ID.
		id := h.GenerateID()
		msg, err := json.Marshal(SubscribeMessageRequest{
			Method: string(SubscribeMethod),
			Params: params,
			ID:     id,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal subscribe message: %w", err)
		}

		// Set the IDs
		h.SetIDForInstruments(id, batch)
		msgs[i] = msg
	}

	return msgs, nil
}

// SetIDForInstruments sets the ID for the given instruments. This is used to set the ID for the
// instruments that are being subscribed to.
func (h *WebSocketHandler) SetIDForInstruments(id int64, instruments []string) {
	h.messageIDs[id] = instruments
}

// GenerateID generates a random ID for the message.
func (h *WebSocketHandler) GenerateID() int64 {
	h.nextID++
	return h.nextID
}
