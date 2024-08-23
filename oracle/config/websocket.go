package config

import (
	"fmt"
	"time"
)

const (
	// DefaultMaxBufferSize is the default maximum number of messages that the provider
	// will buffer at any given time.
	DefaultMaxBufferSize = 1024

	// DefaultReconnectionTimeout is the default timeout for the provider to attempt
	// to reconnect to the websocket endpoint.
	DefaultReconnectionTimeout = 10 * time.Second

	// DefaultPostConnectionTimeout is the default timeout for the provider to wait
	// after a connection is established before sending messages.
	DefaultPostConnectionTimeout = 1 * time.Second

	// DefaultReadBufferSize is the default I/O read buffer size. If a buffer size of
	// 0 is specified, then a default buffer size is used i.e. the buffers allocated
	// by the HTTP server.
	DefaultReadBufferSize = 0

	// DefaultWriteBufferSize is the default I/O write buffer size. If a buffer size of
	// 0 is specified, then a default buffer size is used i.e. the buffers allocated
	// by the HTTP server.
	DefaultWriteBufferSize = 0

	// DefaultHandshakeTimeout is the default duration for the handshake to complete.
	DefaultHandshakeTimeout = 10 * time.Second

	// DefaultEnableCompression is the default value for whether the client should
	// attempt to negotiate per message compression (RFC 7692).
	DefaultEnableCompression = false

	// DefaultReadTimeout is the default read deadline on the underlying network
	// connection.
	DefaultReadTimeout = 10 * time.Second

	// DefaultWriteTimeout is the default write deadline on the underlying network
	// connection.
	DefaultWriteTimeout = 5 * time.Second

	// DefaultPingInterval is the default interval at which the provider should send
	// ping messages to the server.
	DefaultPingInterval = 0 * time.Second

	// DefaultWriteInterval is the default interval at which the provider should send
	// write messages to the server.
	DefaultWriteInterval = 100 * time.Millisecond

	// DefaultMaxReadErrorCount is the default maximum number of read errors that
	// the provider will tolerate before closing the connection and attempting to
	// reconnect. This default value utilized by the gorilla/websocket package is
	// 1000, but we set it to a lower value to allow the provider to reconnect
	// faster.
	DefaultMaxReadErrorCount = 100

	// DefaultMaxSubscriptionsPerConnection is the default maximum subscriptions
	// a provider can handle per-connection.  When this value is 0, one connection
	// will handle all subscriptions.
	DefaultMaxSubscriptionsPerConnection = 0

	// DefaultMaxSubscriptionsPerBatch is the default maximum number of subscriptions
	// that can be assigned to a single batch/write/message.
	DefaultMaxSubscriptionsPerBatch = 1
)

// WebSocketConfig defines a config for a websocket based data provider.
type WebSocketConfig struct {
	// Enabled indicates if the provider is enabled.
	Enabled bool `json:"enabled"`

	// MaxBufferSize is the maximum number of messages that the provider will buffer
	// at any given time. If the provider receives more messages than this, it will
	// block receiving messages until the buffer is cleared.
	MaxBufferSize int `json:"maxBufferSize"`

	// ReconnectionTimeout is the timeout for the provider to attempt to reconnect
	// to the websocket endpoint.
	ReconnectionTimeout time.Duration `json:"reconnectionTimeout"`

	// PostConnectionTimeout is the timeout for the provider to wait after a connection
	// is established before sending messages.
	PostConnectionTimeout time.Duration `json:"postConnectionTimeout"`

	// Endpoints are the websocket endpoints for the provider. At least one endpoint
	// must be specified.
	Endpoints []Endpoint `json:"endpoints"`

	// Name is the name of the provider that corresponds to this config.
	Name string `json:"name"`

	// ReadBufferSize specifies the I/O read buffer size. if a buffer size of 0 is
	// specified, then a default buffer size is used.
	ReadBufferSize int `json:"readBufferSize"`

	// WriteBufferSize specifies the I/O write buffer size. if a buffer size of 0 is
	// specified, then a default buffer size is used.
	WriteBufferSize int `json:"writeBufferSize"`

	// HandshakeTimeout specifies the duration for the handshake to complete.
	HandshakeTimeout time.Duration `json:"handshakeTimeout"`

	// EnableCompression specifies if the client should attempt to negotiate per
	// message compression (RFC 7692). Setting this value to true does not guarantee
	// that compression will be supported. Note that enabling compression may
	EnableCompression bool `json:"enableCompression"`

	// ReadTimeout sets the read deadline on the underlying network connection.
	// After a read has timed out, the websocket connection state is corrupt and
	// all future reads will return an error. A zero value for t means reads will
	// not time out.
	ReadTimeout time.Duration `json:"readTimeout"`

	// WriteTimeout sets the write deadline on the underlying network
	// connection. After a write has timed out, the websocket state is corrupt and
	// all future writes will return an error. A zero value for t means writes will
	// not time out.
	WriteTimeout time.Duration `json:"writeTimeout"`

	// PingInterval is the interval to ping the server. Note that a ping interval
	// of 0 disables pings.
	PingInterval time.Duration `json:"pingInterval"`

	// WriteInterval is the interval at which the provider should wait before sending
	// consecutive write messages to the server.
	WriteInterval time.Duration `json:"writeInterval"`

	// MaxReadErrorCount is the maximum number of read errors that the provider
	// will tolerate before closing the connection and attempting to reconnect.
	MaxReadErrorCount int `json:"maxReadErrorCount"`

	// MaxSubscriptionsPerConnection is the maximum amount of subscriptions that
	// can be assigned to a single connection for this provider.  The null value (0),
	// indicates that there is no limit per connection.
	MaxSubscriptionsPerConnection int `json:"maxSubscriptionsPerConnection"`

	// MaxSubscriptionsPerBatch is the maximum number of subscription messages that the
	// provider will send in a single batch/write.
	MaxSubscriptionsPerBatch int `json:"maxSubscriptionsPerBatch"`
}

// ValidateBasic performs basic validation of the websocket config.
func (c *WebSocketConfig) ValidateBasic() error {
	if !c.Enabled {
		return nil
	}

	if c.MaxBufferSize < 1 {
		return fmt.Errorf("websocket max buffer size must be greater than 0")
	}

	if c.ReconnectionTimeout <= 0 {
		return fmt.Errorf("websocket reconnection timeout must be greater than 0")
	}

	if c.PostConnectionTimeout < 0 {
		return fmt.Errorf("websocket post connection timeout must be greater than 0")
	}

	if len(c.Endpoints) == 0 {
		return fmt.Errorf("websocket endpoints cannot be empty")
	}

	for i, e := range c.Endpoints {
		if err := e.ValidateBasic(); err != nil {
			return fmt.Errorf("endpoint %d: %w", i, err)
		}
	}

	if len(c.Name) == 0 {
		return fmt.Errorf("websocket name cannot be empty")
	}

	if c.ReadBufferSize < 0 {
		return fmt.Errorf("websocket read buffer size cannot be negative")
	}

	if c.WriteBufferSize < 0 {
		return fmt.Errorf("websocket write buffer size cannot be negative")
	}

	if c.HandshakeTimeout <= 0 {
		return fmt.Errorf("websocket handshake timeout must be greater than 0")
	}

	if c.ReadTimeout <= 0 {
		return fmt.Errorf("websocket read timeout must be greater than 0")
	}

	if c.WriteTimeout <= 0 {
		return fmt.Errorf("websocket write timeout must be greater than 0")
	}

	if c.PingInterval < 0 {
		return fmt.Errorf("websocket ping interval cannot be negative")
	}

	if c.WriteInterval < 0 {
		return fmt.Errorf("websocket write interval must be greater than 0")
	}

	if c.MaxReadErrorCount < 0 {
		return fmt.Errorf("websocket max read error count cannot be negative")
	}

	if c.MaxSubscriptionsPerConnection < 0 {
		return fmt.Errorf("websocket max subscriptions per connection cannot be negative")
	}

	if c.MaxSubscriptionsPerBatch < 1 {
		return fmt.Errorf("websocket max subscriptions per batch must be greater than 0")
	}

	return nil
}
