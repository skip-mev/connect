package handlers

import (
	"sync"

	"github.com/gorilla/websocket"
)

// WebsocketEncodedMessage is a type alias for a websocket message encoded to bytes.
type WebsocketEncodedMessage []byte

// WebSocketConnHandler is an interface the encapsulates the functionality of a web socket
// connection to a data provider. It provides the simple CRUD operations for a web socket
// connection. The connection handler is responsible for managing the connection to the
// data provider. This includes creating the connection, reading messages, writing messages,
// and closing the connection.
//
//go:generate mockery --name WebSocketConnHandler --output ./mocks/ --case underscore
type WebSocketConnHandler interface {
	// Read is used to read data from the data provider. This should block until data is
	// received from the data provider.
	Read() ([]byte, error)

	// Write is used to write messages to the data provider. Write should block until the
	// message is sent to the data provider.
	Write(message []byte) error

	// Close is used to close the connection to the data provider. Any additional cleanup
	// should be done here.
	Close() error

	// Dial is used to create the connection to the data provider.
	Dial(url string) error

	// Heartbeat is an optional routine used to keep a connection open by sending heartbeat
	// messages to the server.
	Heartbeat() error
}

// WebSocketConnHandlerImpl is a struct that implements the WebSocketConnHandler interface.
type WebSocketConnHandlerImpl struct {
	sync.Mutex

	// conn is the connection to the data provider.
	conn *websocket.Conn
}

// NewWebSocketHandlerImpl returns a new WebSocketConnHandlerImpl.
func NewWebSocketHandlerImpl() WebSocketConnHandler {
	return &WebSocketConnHandlerImpl{}
}

// Dial is used to create a new connection to the data provider with the given URL.
func (h *WebSocketConnHandlerImpl) Dial(url string) error {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	// TODO: Determine whether the default dialer is safe to use.
	var err error
	h.conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	return err
}

// Read is used to read data from the data provider. Each web socket data handler is responsible
// for determining how to parse the data and being aware of the data format (text, json, etc.).
func (h *WebSocketConnHandlerImpl) Read() ([]byte, error) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	_, message, err := h.conn.ReadMessage()
	return message, err
}

// Write is used to write messages to the data provider. By default, all messages are sent as
// text messages. This permits encoding json messages as text messages.
func (h *WebSocketConnHandlerImpl) Write(message []byte) error {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	return h.conn.WriteMessage(websocket.TextMessage, message)
}

// Close is used to close the connection to the data provider.
func (h *WebSocketConnHandlerImpl) Close() error {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	// Cleanly close the connection by sending a close message and then
	// waiting (with a timeout) for the server to close the connection.
	err := h.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}

	return h.conn.Close()
}

// Heartbeat is a no-op by default.
func (h *WebSocketConnHandlerImpl) Heartbeat() error {
	return nil
}
