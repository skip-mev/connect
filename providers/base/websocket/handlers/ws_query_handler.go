package handlers

import (
	"context"
	"fmt"

	"github.com/skip-mev/slinky/providers/base/websocket/errors"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"go.uber.org/zap"
)

// WebSocketQueryHandler is an interface that encapsulates querying a web socket
// data provider for info. The handler must respect the context timeout and close
// the connection if the context is cancelled. All responses must be sent to the
// response channel. These are processed asynchronously by the provider.
//
//go:generate mockery --name WebSocketQueryHandler --output ./mocks/ --case underscore
type WebSocketQueryHandler[K comparable, V any] interface {
	// Start should initialize the web socket connection and start listening for
	// the data (i.e. ids). All web socket responses should be sent to the response
	// channel.
	Start(ctx context.Context, ids []K, responseCh chan<- providertypes.GetResponse[K, V]) error
}

// WebSocketQueryHandlerImpl is the default web socket implementation of the
// WebSocketQueryHandler interface. This is used to establish a connection to the data
// provider and subscribe to events for a given set of IDs. It runs in a separate go
// routine and will send all responses to the response channel as they are received.
type WebSocketQueryHandlerImpl[K comparable, V any] struct {
	logger *zap.Logger

	// The connection handler is used to manage the connection to the data provider. This
	// establishes the connection and sends/receives messages to/from the data provider.
	connHandler WebSocketConnHandler

	// The data handler is used to handle messages received from the data provider. This
	// is used to parse the messages and send responses to the provider.
	dataHandler WebSocketDataHandler[K, V]

	// ids is the set of IDs that the provider will fetch data for.
	ids []K
}

// NewWebSocketQueryHandler creates a new web socket query handler.
func NewWebSocketQueryHandler[K comparable, V any](
	logger *zap.Logger,
	dataHandler WebSocketDataHandler[K, V],
	connHandler WebSocketConnHandler,
) (WebSocketQueryHandler[K, V], error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	if dataHandler == nil {
		return nil, fmt.Errorf("data handler is nil")
	}

	if connHandler == nil {
		return nil, fmt.Errorf("connection is nil")
	}

	return &WebSocketQueryHandlerImpl[K, V]{
		logger:      logger.With(zap.String("web_socket_data_handler", dataHandler.Name())),
		dataHandler: dataHandler,
		connHandler: connHandler,
	}, nil
}

// Start is used to start the connection to the data provider and start listening for
// the data (i.e. ids). All web socket responses should be sent to the response channel.
// Start will first:
//  1. Create the initial set of events that the channel will subscribe to using the data
//     handler.
//  2. Start the connection to the data provider using the connection handler and url from
//     the data handler.
//  3. Send the initial payload to the data provider.
//  4. Start receiving messages from the data provider.
func (h *WebSocketQueryHandlerImpl[K, V]) Start(
	ctx context.Context,
	ids []K,
	responseCh chan<- providertypes.GetResponse[K, V],
) error {
	defer func() {
		if err := recover(); err != nil {
			h.logger.Error("panic occurred", zap.Any("err", err))
		}
	}()

	if responseCh == nil {
		h.logger.Error("response channel is nil")
		return fmt.Errorf("response channel is nil")
	}

	h.ids = ids
	if len(h.ids) == 0 {
		h.logger.Debug("no ids to query")
		return nil
	}

	// Initialize the connection to the data provider and subscribe to the events
	// for the corresponding IDs.
	if err := h.start(); err != nil {
		responseCh <- providertypes.NewGetResponseWithErr[K, V](ids, err)
		return fmt.Errorf("failed to start connection: %w", err)
	}

	// Start receiving messages from the data provider.
	return h.recv(ctx, responseCh)
}

// start is used to start the connection to the data provider.
func (h *WebSocketQueryHandlerImpl[K, V]) start() error {
	url := h.dataHandler.URL()
	h.logger.Debug(
		"creating connection to data provider",
		zap.String("url", url),
	)

	// Start the connection.
	if err := h.connHandler.Dial(url); err != nil {
		h.logger.Error(
			"failed to create connection with data provider",
			zap.Error(err),
		)

		return errors.ErrDialWithErr(err)
	}

	// Create the initial set of events that the channel will subscribe to.
	message, err := h.dataHandler.CreateMessage(h.ids)
	if err != nil {
		h.logger.Error(
			"failed to create subscription messages",
			zap.Error(err),
		)

		return errors.ErrCreateMessageWithErr(err)
	}

	h.logger.Debug(
		"connection created; sending initial payload",
		zap.Binary("payload", message),
	)

	// Send the initial payload to the data provider.
	if err := h.connHandler.Write(message); err != nil {
		h.logger.Error(
			"failed to write message to web socket connection handler",
			zap.Error(err),
		)

		return errors.ErrWriteWithErr(err)
	}

	h.logger.Debug("initial payload sent; web socket connection successfully started")
	return nil
}

// recv is used to manage the connection to the data provider.
func (h *WebSocketQueryHandlerImpl[K, V]) recv(ctx context.Context, responseCh chan<- providertypes.GetResponse[K, V]) error {
	defer func() {
		if err := recover(); err != nil {
			h.logger.Error("panic occurred", zap.Any("err", err))
		}
	}()

	h.logger.Debug(
		"starting recv",
		zap.Int("buffer_size", cap(responseCh)),
	)

	for {
		select {
		case <-ctx.Done():
			h.logger.Debug("context finished; closing connection to web socket handler")
			if err := h.connHandler.Close(); err != nil {
				h.logger.Error(
					"failed to close connection",
					zap.Error(err),
				)

				return errors.ErrCloseWithErr(err)
			}

			h.logger.Debug("connection closed")
			return ctx.Err()
		default:
			// Wait for a message from the data provider.
			message, err := h.connHandler.Read()
			if err != nil {
				h.logger.Error(
					"failed to read message from web socket handler",
					zap.Error(err),
				)

				continue
			}

			h.logger.Debug(
				"message received; attempting to handle message",
				zap.Binary("message", message),
			)

			// Handle the message.
			response, updateMessage, err := h.dataHandler.HandleMessage(message)
			if err != nil {
				h.logger.Error(
					"failed to handle web socket message",
					zap.Error(err),
				)

				continue
			}

			// Immediately send the response to the response channel. Even if this is
			// empty, it will be handled by the provider.
			responseCh <- response
			h.logger.Debug(
				"handled message successfully; sent response to response channel",
				zap.String("response", response.String()),
			)

			// If the update message is not nil, send it to the data provider.
			if len(updateMessage) != 0 {
				h.logger.Debug(
					"sending update message to data provider",
					zap.Binary("update_message", updateMessage),
				)

				if err := h.connHandler.Write(updateMessage); err != nil {
					h.logger.Error(
						"failed to write update message",
						zap.Error(err),
					)
				} else {
					h.logger.Debug("update message sent")
				}

			}
		}
	}
}
