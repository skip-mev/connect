package bitfinex

type (
	// Event is the event type of a message sent over the Bitfinex websocket API.
	Event string

	// Channel is the channel of a message sent over the Bitfinex websocket API.
	Channel string
)

const (
	EVENT_SUBSCRIBE  Event = "subscribe"
	EVENT_SUBSCRIBED Event = "subscribed"

	CHANNEL_TICKER Channel = "ticker"
)

type BaseMessage struct {
	Event   string `json:"event" validate:"required"`
	Channel string `json:"channel" validate:"required"`
	Symbol  string `json:"symbol" validate:"required"`
}

type SubscribeMessage BaseMessage

type SubscribedMessage struct {
	BaseMessage
	ChannelID string `json:"chanId" validate:"required"`
	Pair      string `json:"pair" validate:"required"`
}

type TickerStream struct {
	ChannelID string  `json:"chanId" validate:"required"`
	LastPrice float64 `json:"lastPrice" validate:"required"`
}
