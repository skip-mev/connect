package bitfinex

type (
	// Event is the event type of a message sent over the Bitfinex websocket API.
	Event string

	// Channel is the channel of a message sent over the Bitfinex websocket API.
	Channel string
)

const (
	// EVENT_SUBSCRIBE indicates a subscribe action.
	EVENT_SUBSCRIBE Event = "subscribe"
	// EVENT_SUBSCRIBED indicates that a subscription was successful.
	EVENT_SUBSCRIBED Event = "subscribed"
	// CHANNEL_TICKER is the channel name for the ticker channel.
	CHANNEL_TICKER Channel = "ticker"
)

// BaseMessage is the base message structure for subscription requests and responses in the
// BitFinex websocket API.
type BaseMessage struct {
	Event   string `json:"event" validate:"required"`
	Channel string `json:"channel" validate:"required"`
	Symbol  string `json:"symbol" validate:"required"`
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
type SubscribeMessage BaseMessage

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
	ChannelID string `json:"chanId" validate:"required"`
	Pair      string `json:"pair" validate:"required"`
}

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
	ChannelID string  `json:"chanId" validate:"required"`
	LastPrice float64 `json:"lastPrice" validate:"required"`
}
