package handlers

import (
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// WebSocketDataHandler defines an interface that must be implemented by all providers that
// want to fetch data from a web socket. This interface is meant to be paired with the
// WebSocketQueryHandler. The WebSocketQueryHandler will use the WebSocketDataHandler to
// create establish a connection to the correct host, create subscription messages to be sent
// to the data provider, and handle incoming events accordingly.
//
//go:generate mockery --name WebSocketDataHandler --output ./mocks/ --case underscore
type WebSocketDataHandler[K comparable, V any] interface {
	// HandleMessage is used to handle a message received from the data provider. Message parsing
	// and response creation should be handled by this data handler. Given a message from the web socket
	// the handler should either return a response or an update message.
	HandleMessage(message []byte) (response providertypes.GetResponse[K, V], updateMessage []byte, err error)

	// CreateMessage is used to update the connection to the data provider. This can be used to subscribe
	// to new events or unsubscribe from events.
	CreateMessage(ids []K) ([]byte, error)
}
