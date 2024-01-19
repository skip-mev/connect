package bybit

import (
	"sync"

	"github.com/gorilla/websocket"

	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
)

// WebSocketConnHandler is a struct that implements the WebSocketConnHandler interface.
type WebSocketConnHandler struct {
	sync.Mutex

	// conn is the connection to the data provider.
	conn *websocket.Conn
}

// NewWebSocketHandlerImpl returns a new WebSocketConnHandlerImpl.
func NewWebSocketHandlerImpl() handlers.WebSocketConnHandler {
	return &WebSocketConnHandler{}
}

// Dial is used to create a new connection to the data provider with the given URL.
func (h *WebSocketConnHandler) Dial(url string) error {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	// TODO: Determine whether the default dialer is safe to use.
	var err error
	h.conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	return err
}

// Read is used to read data from the data provider. Each web socket data handler is responsible
// for determining how to parse the data and being aware of the data format (text, json, etc.).
func (h *WebSocketConnHandler) Read() ([]byte, error) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	_, message, err := h.conn.ReadMessage()
	return message, err
}

// Write is used to write messages to the data provider. By default, all messages are sent as
// text messages. This permits encoding json messages as text messages.
func (h *WebSocketConnHandler) Write(message []byte) error {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	return h.conn.WriteMessage(websocket.TextMessage, message)
}

// Close is used to close the connection to the data provider.
func (h *WebSocketConnHandler) Close() error {
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
func (h *WebSocketConnHandler) Heartbeat() error {
	return nil
}
