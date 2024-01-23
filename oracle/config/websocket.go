package config

import (
	"fmt"
	"time"
)

// WebSocketConfig defines a config for a websocket based data provider.
type WebSocketConfig struct {
	// Enabled is a flag that indicates whether the provider is web socket based.
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
	// that compression will be supported.
	EnableCompression bool `mapstructure:"enable_compression" toml:"enable_compression"`

	// ReadDeadline sets the read deadline on the underlying network connection.
	// After a read has timed out, the websocket connection state is corrupt and
	// all future reads will return an error. A zero value for t means reads will
	// not time out.
	ReadDeadline time.Duration `mapstructure:"read_deadline" toml:"read_deadline"`

	// WriteDeadline sets the write deadline on the underlying network
	// connection. After a write has timed out, the websocket state is corrupt and
	// all future writes will return an error. A zero value for t means writes will
	// not time out.
	WriteDeadline time.Duration `mapstructure:"write_deadline" toml:"write_deadline"`
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

	return nil
}
