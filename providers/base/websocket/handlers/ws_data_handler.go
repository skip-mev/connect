package handlers

import (
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

// WebSocketDataHandler defines an interface that must be implemented by all providers that
// want to fetch data from a websocket. This interface is meant to be paired with the
// WebSocketQueryHandler. The WebSocketQueryHandler will use the WebSocketDataHandler to
// create establish a connection to the correct host, create subscription messages to be sent
// to the data provider, and handle incoming events accordingly.
//
//go:generate mockery --name WebSocketDataHandler --output ./mocks/ --case underscore
type WebSocketDataHandler[K providertypes.ResponseKey, V providertypes.ResponseValue] interface {
	// HandleMessage is used to handle a message received from the data provider. Message parsing
	// and response creation should be handled by this data handler. Given a message from the websocket
	// the handler should either return a response or a set of update messages.
	HandleMessage(message []byte) (response providertypes.GetResponse[K, V], updateMessages []WebsocketEncodedMessage, err error)

	// CreateMessages is used to update the connection to the data provider. This can be used to subscribe
	// to new events or unsubscribe from events.
	CreateMessages(ids []K) ([]WebsocketEncodedMessage, error)

	// HeartBeatMessages is used to construct heartbeat messages to be sent to the data provider. Note that
	// the handler must maintain the necessary state information to construct the heartbeat messages. This
	// can be done on the fly as messages as handled by the handler.
	HeartBeatMessages() ([]WebsocketEncodedMessage, error)

	// Copy is used to create a copy of the data handler. This is useful for creating multiple connections
	// to the same data provider. Stateful information can be managed independently for each connection.
	Copy() WebSocketDataHandler[K, V]
}
