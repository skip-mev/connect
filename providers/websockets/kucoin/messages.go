package kucoin

// MessageType represents the type of message received from the Kucoin websocket.
type MessageType string

const (
	// WelcomeMessage represents the welcome message received when first connecting
	// to the websocket.
	//
	// ref: https://www.kucoin.com/docs/websocket/basic-info/create-connection
	WelcomeMessage MessageType = "welcome"
)

// WelcomeResponseMessage represents the welcome message received when first
// connecting to the websocket.
//
//	{
//			"id": "hQvf8jkno",
//			"type": "welcome"
//	}
//
// ref: https://www.kucoin.com/docs/websocket/basic-info/create-connection
type WelcomeResponseMessage struct {
	// ID is the ID of the message.
	ID string `json:"id"`

	// Type is the type of message.
	Type string `json:"type"`
}
