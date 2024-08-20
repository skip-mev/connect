package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/skip-mev/connect/v2/oracle/config"
)

type (
	// WebsocketEncodedMessage is a type alias for a websocket message encoded to bytes.
	WebsocketEncodedMessage []byte

	// PreDialHook is a function that is called before the connection is established. This
	// is useful for dynamically generating configurations for the connection.
	PreDialHook func(*WebSocketConnHandlerImpl) error
)

// WebSocketConnHandler is an interface the encapsulates the functionality of a websocket
// connection to a data provider. It provides the simple CRUD operations for a websocket
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
	Dial() error

	// Copy is used to create a copy of the connection handler. This is useful for creating
	// multiple connections to the same data provider.
	Copy() WebSocketConnHandler
}

// WebSocketConnHandlerImpl is a struct that implements the WebSocketConnHandler interface.
type WebSocketConnHandlerImpl struct {
	sync.Mutex
	cfg config.WebSocketConfig

	// conn is the connection to the data provider.
	conn *websocket.Conn

	// preDialHook is a function that is called before the connection is established.
	preDialHook PreDialHook
}

// NewWebSocketHandlerImpl returns a new WebSocketConnHandlerImpl.
func NewWebSocketHandlerImpl(cfg config.WebSocketConfig, opts ...Option) (*WebSocketConnHandlerImpl, error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, err
	}

	h := &WebSocketConnHandlerImpl{
		cfg: cfg,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h, nil
}

// CreateDialer is a function that dynamically creates a new websocket dialer.
func (h *WebSocketConnHandlerImpl) CreateDialer() *websocket.Dialer {
	return &websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  h.cfg.HandshakeTimeout,
		ReadBufferSize:    h.cfg.ReadBufferSize,
		WriteBufferSize:   h.cfg.WriteBufferSize,
		EnableCompression: h.cfg.EnableCompression,
	}
}

// Dial is used to create a new connection to the data provider with the given URL.
func (h *WebSocketConnHandlerImpl) Dial() error {
	if h.preDialHook != nil {
		if err := h.preDialHook(h); err != nil {
			return err
		}
	}

	if len(h.cfg.Endpoints) == 0 {
		return fmt.Errorf("no endpoints provided")
	}

	var err error
	h.conn, _, err = h.CreateDialer().Dial(h.cfg.Endpoints[0].URL, nil)
	return err
}

// Read is used to read data from the data provider. Each websocket data handler is responsible
// for determining how to parse the data and being aware of the data format (text, json, etc.).
func (h *WebSocketConnHandlerImpl) Read() ([]byte, error) {
	h.Lock()
	defer h.Unlock()

	if h.conn == nil {
		return nil, fmt.Errorf("connection has not been established")
	}

	// Set the read deadline to the configured read timeout.
	if err := h.conn.SetReadDeadline(time.Now().Add(h.cfg.ReadTimeout)); err != nil {
		return nil, err
	}

	_, message, err := h.conn.ReadMessage()
	return message, err
}

// Write is used to write messages to the data provider. By default, all messages are sent as
// text messages. This permits encoding json messages as text messages.
func (h *WebSocketConnHandlerImpl) Write(message []byte) error {
	h.Lock()
	defer h.Unlock()

	if h.conn == nil {
		return fmt.Errorf("connection has not been established")
	}

	// Set the write deadline to the configured write timeout.
	if err := h.conn.SetWriteDeadline(time.Now().Add(h.cfg.WriteTimeout)); err != nil {
		return err
	}

	return h.conn.WriteMessage(websocket.TextMessage, message)
}

// Close is used to close the connection to the data provider.
func (h *WebSocketConnHandlerImpl) Close() error {
	h.Lock()
	defer h.Unlock()

	if h.conn == nil {
		return fmt.Errorf("connection has not been established")
	}

	// Set the write deadline to the configured write timeout.
	if err := h.conn.SetWriteDeadline(time.Now().Add(h.cfg.WriteTimeout)); err != nil {
		return err
	}

	// Cleanly close the connection by sending a close message and then
	// waiting (with a timeout) for the server to close the connection.
	err := h.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}

	return h.conn.Close()
}

// Copy is used to create a copy of the connection handler. This is useful for creating multiple
// connections to the same data provider.
func (h *WebSocketConnHandlerImpl) Copy() WebSocketConnHandler {
	h.Lock()
	defer h.Unlock()

	return &WebSocketConnHandlerImpl{
		cfg:         h.cfg,
		preDialHook: h.preDialHook,
	}
}

// GetConfig is used to get the configuration for the connection handler.
func (h *WebSocketConnHandlerImpl) GetConfig() config.WebSocketConfig {
	h.Lock()
	defer h.Unlock()

	return h.cfg
}

// SetConfig is used to set the configuration for the connection handler.
func (h *WebSocketConnHandlerImpl) SetConfig(cfg config.WebSocketConfig) {
	h.Lock()
	defer h.Unlock()

	h.cfg = cfg
}
