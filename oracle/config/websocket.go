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
}

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

	return nil
}
