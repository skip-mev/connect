package bybit

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"

	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
)

// WebSocketConnHandler is a struct that implements the WebSocketConnHandler interface.
type WebSocketConnHandler struct {
	sync.Mutex

	logger *zap.Logger

	// conn is the connection to the data provider.
	conn *websocket.Conn
}

// NewWebSocketHandler returns a new WebSocketConnHandler.
func NewWebSocketHandler(logger *zap.Logger) handlers.WebSocketConnHandler {
	return &WebSocketConnHandler{logger: logger.With(zap.String("web_socket_conn_handler", Name))}
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

// Heartbeat sends a heartbeat ping to the server every 20 seconds until the context is cancelled.
func (h *WebSocketConnHandler) Heartbeat(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("shutting down heartbeat routine")
			return nil
		default:
			msg, err := json.Marshal(HeartbeatPing{BaseRequest{
				ReqID: time.Now().String(),
				Op:    string(OperationPing),
			}})
			if err != nil {
				h.logger.Debug("unable to marshal heartbeat ping")
				return fmt.Errorf("unable to marshal heartbeat ping")
			}

			err = h.Write(msg)
			if err != nil {
				h.logger.Debug("unable to write heartbeat ping")
				return fmt.Errorf("unable to write heartbeat ping")
			}

			h.logger.Debug("sent heartbeat message")
			time.Sleep(20 * time.Second)
		}
	}
}
