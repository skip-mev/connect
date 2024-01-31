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

	// DefaultReadBufferSize is the default I/O read buffer size. If a buffer size of
	// 0 is specified, then a default buffer size is used i.e. the buffers allocated
	// by the HTTP server.
	DefaultReadBufferSize = 0

	// DefaultWriteBufferSize is the default I/O write buffer size. If a buffer size of
	// 0 is specified, then a default buffer size is used i.e. the buffers allocated
	// by the HTTP server.
	DefaultWriteBufferSize = 0

	// DefaultHandshakeTimeout is the default duration for the handshake to complete.
	DefaultHandshakeTimeout = 45 * time.Second

	// DefaultEnableCompression is the default value for whether the client should
	// attempt to negotiate per message compression (RFC 7692).
	DefaultEnableCompression = false

	// DefaultReadTimeout is the default read deadline on the underlying network
	// connection.
	DefaultReadTimeout = 45 * time.Second

	// DefaultWriteTimeout is the default write deadline on the underlying network
	// connection.
	DefaultWriteTimeout = 45 * time.Second

	// DefaultPingInterval is the default interval at which the provider should send
	// ping messages to the server.
	DefaultPingInterval = 0 * time.Second

	// DefaultMaxReadErrorCount is the default maximum number of read errors that
	// the provider will tolerate before closing the connection and attempting to
	// reconnect. This default is taken from the default value for the websocket
	// library (gorilla/websocket).
	DefaultMaxReadErrorCount = 1000
)

// WebSocketConfig defines a config for a websocket based data provider.
type WebSocketConfig struct {
	// Enabled is a flag that indicates whether the provider is websocket based.
	Enabled bool `mapstructure:"enabled" toml:"enabled"`

	// MaxBufferSize is the maximum number of messages that the provider will buffer
	// at any given time. If the provider receives more messages than this, it will
	// block receiving messages until the buffer is cleared.
	MaxBufferSize int `mapstructure:"max_buffer_size" toml:"max_buffer_size"`

	// ReconnectionTimeout is the timeout for the provider to attempt to reconnect
	// to the websocket endpoint.
	ReconnectionTimeout time.Duration `mapstructure:"reconnection_timeout" toml:"reconnection_timeout"`

	// WSS is the websocket endpoint for the provider.
	WSS string `mapstructure:"wss" toml:"wss"`

	// Name is the name of the provider that corresponds to this config.
	Name string `mapstructure:"name" toml:"name"`

	// ReadBufferSize specifies the I/O read buffer size. if a buffer size of 0 is
	// specified, then a default buffer size is used.
	ReadBufferSize int `mapstructure:"read_buffer_size" toml:"read_buffer_size"`

	// WriteBufferSize specifies the I/O write buffer size. if a buffer size of 0 is
	// specified, then a default buffer size is used.
	WriteBufferSize int `mapstructure:"write_buffer_size" toml:"write_buffer_size"`

	// HandshakeTimeout specifies the duration for the handshake to complete.
	HandshakeTimeout time.Duration `mapstructure:"handshake_timeout" toml:"handshake_timeout"`

	// EnableCompression specifies if the client should attempt to negotiate per
	// message compression (RFC 7692). Setting this value to true does not guarantee
	// that compression will be supported. Note that enabling compression may
	EnableCompression bool `mapstructure:"enable_compression" toml:"enable_compression"`

	// ReadTimeout sets the read deadline on the underlying network connection.
	// After a read has timed out, the websocket connection state is corrupt and
	// all future reads will return an error. A zero value for t means reads will
	// not time out.
	ReadTimeout time.Duration `mapstructure:"read_deadline" toml:"read_deadline"`

	// WriteTimeout sets the write deadline on the underlying network
	// connection. After a write has timed out, the websocket state is corrupt and
	// all future writes will return an error. A zero value for t means writes will
	// not time out.
	WriteTimeout time.Duration `mapstructure:"write_deadline" toml:"write_deadline"`

	// PingInterval is the interval to ping the server. Note that a ping interval
	// of 0 disables pings.
	PingInterval time.Duration `mapstructure:"ping_interval" toml:"ping_interval"`

	// MaxReadErrorCount is the maximum number of read errors that the provider
	// will tolerate before closing the connection and attempting to reconnect.
	MaxReadErrorCount int `mapstructure:"max_read_error_count" toml:"max_read_error_count"`
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

	if len(c.WSS) == 0 {
		return fmt.Errorf("websocket endpoint cannot be empty")
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

	if c.MaxReadErrorCount < 0 {
		return fmt.Errorf("websocket max read error count cannot be negative")
	}

	return nil
}
