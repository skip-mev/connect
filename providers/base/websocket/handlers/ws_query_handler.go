package handlers

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/websocket/errors"
	"github.com/skip-mev/slinky/providers/base/websocket/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// WebSocketQueryHandler is an interface that encapsulates querying a websocket
// data provider for info. The handler must respect the context timeout and close
// the connection if the context is cancelled. All responses must be sent to the
// response channel. These are processed asynchronously by the provider.
//
//go:generate mockery --name WebSocketQueryHandler --output ./mocks/ --case underscore
type WebSocketQueryHandler[K providertypes.ResponseKey, V providertypes.ResponseValue] interface {
	// Start should initialize the websocket connection and start listening for
	// the data (i.e. ids). All websocket responses should be sent to the response
	// channel.
	Start(ctx context.Context, ids []K, responseCh chan<- providertypes.GetResponse[K, V]) error

	// Copy is used to create a copy of the query handler. This is useful for creating
	// multiple connections to the same data provider.
	Copy() WebSocketQueryHandler[K, V]
}

// WebSocketQueryHandlerImpl is the default websocket implementation of the
// WebSocketQueryHandler interface. This is used to establish a connection to the data
// provider and subscribe to events for a given set of IDs. It runs in a separate go
// routine and will send all responses to the response channel as they are received.
type WebSocketQueryHandlerImpl[K providertypes.ResponseKey, V providertypes.ResponseValue] struct {
	logger  *zap.Logger
	metrics metrics.WebSocketMetrics
	config  config.WebSocketConfig

	// The connection handler is used to manage the connection to the data provider. This
	// establishes the connection and sends/receives messages to/from the data provider.
	connHandler WebSocketConnHandler

	// The data handler is used to handle messages received from the data provider. This
	// is used to parse the messages and send responses to the provider.
	dataHandler WebSocketDataHandler[K, V]

	// ids is the set of IDs that the provider will fetch data for.
	ids []K
}

// NewWebSocketQueryHandler creates a new websocket query handler.
func NewWebSocketQueryHandler[K providertypes.ResponseKey, V providertypes.ResponseValue](
	logger *zap.Logger,
	config config.WebSocketConfig,
	dataHandler WebSocketDataHandler[K, V],
	connHandler WebSocketConnHandler,
	m metrics.WebSocketMetrics,
) (WebSocketQueryHandler[K, V], error) {
	if err := config.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	if !config.Enabled {
		return nil, fmt.Errorf("websocket is not enabled")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	if dataHandler == nil {
		return nil, fmt.Errorf("websocket data handler is nil")
	}

	if connHandler == nil {
		return nil, fmt.Errorf("connection handler is nil")
	}

	if m == nil {
		return nil, fmt.Errorf("websocket metrics is nil")
	}

	return &WebSocketQueryHandlerImpl[K, V]{
		logger:      logger.With(zap.String("web_socket_data_handler", config.Name)),
		config:      config,
		dataHandler: dataHandler,
		connHandler: connHandler,
		metrics:     m,
	}, nil
}

// Start is used to start the connection to the data provider and start listening for
// the data (i.e. ids). All websocket responses should be sent to the response channel.
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
		h.metrics.AddWebSocketConnectionStatus(h.config.Name, metrics.Unhealthy)
	}()

	if responseCh == nil {
		h.logger.Debug("response channel is nil")
		return fmt.Errorf("response channel is nil")
	}

	h.ids = ids
	if len(h.ids) == 0 {
		h.logger.Debug("no ids to query; exiting")
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Initialize the connection to the data provider and subscribe to the events
	// for the corresponding IDs.
	if err := h.start(); err != nil {
		responseCh <- providertypes.NewGetResponseWithErr[K, V](
			ids,
			providertypes.NewErrorWithCode(
				err,
				providertypes.ErrorWebsocketStartFail,
			),
		)
		return fmt.Errorf("failed to start connection: %w", err)
	}

	if h.config.PingInterval > 0 {
		go h.heartBeat(ctx)
	}

	// Start receiving messages from the data provider.
	h.metrics.AddWebSocketConnectionStatus(h.config.Name, metrics.Healthy)
	return h.recv(ctx, responseCh)
}

// start is used to start the connection to the data provider.
func (h *WebSocketQueryHandlerImpl[K, V]) start() error {
	// Start the connection.
	h.logger.Debug("creating connection to data provider")
	if err := h.connHandler.Dial(); err != nil {
		h.logger.Debug("failed to create connection with data provider", zap.Error(err))
		h.metrics.AddWebSocketConnectionStatus(h.config.Name, metrics.DialErr)
		return errors.ErrDialWithErr(err)
	}

	// Create the initial set of events that the channel will subscribe to.
	h.metrics.AddWebSocketConnectionStatus(h.config.Name, metrics.DialSuccess)
	messages, err := h.dataHandler.CreateMessages(h.ids)
	if err != nil {
		h.logger.Debug("failed to create subscription messages", zap.Error(err))
		h.metrics.AddWebSocketDataHandlerStatus(h.config.Name, metrics.CreateMessageErr)
		return errors.ErrCreateMessageWithErr(err)
	}

	h.metrics.AddWebSocketDataHandlerStatus(h.config.Name, metrics.CreateMessageSuccess)
	for _, message := range messages {
		h.logger.Debug("connection created; sending initial payload", zap.String("payload", string(message)))

		// Send the initial payload to the data provider.
		if err := h.connHandler.Write(message); err != nil {
			h.logger.Debug("failed to write message to websocket connection handler", zap.Error(err))
			h.metrics.AddWebSocketConnectionStatus(h.config.Name, metrics.WriteErr)
			return errors.ErrWriteWithErr(err)
		}
		h.metrics.AddWebSocketConnectionStatus(h.config.Name, metrics.WriteSuccess)
	}

	h.logger.Debug("initial payload sent; websocket connection successfully started")
	return nil
}

// heartBeat is used to send heartbeats to the data provider. This will
// send a heartbeat message to the data provider every ping interval.
func (h *WebSocketQueryHandlerImpl[K, V]) heartBeat(ctx context.Context) {
	ticker := time.NewTicker(h.config.PingInterval)
	defer ticker.Stop()

	h.logger.Debug("starting heartbeat loop", zap.Duration("ping_interval", h.config.PingInterval))

	for {
		select {
		case <-ctx.Done():
			h.logger.Debug("context finished; stopping heartbeat")
			return
		case <-ticker.C:
			h.logger.Debug("creating heartbeat messages")
			msgs, err := h.dataHandler.HeartBeatMessages()
			if err != nil {
				h.metrics.AddWebSocketDataHandlerStatus(h.config.Name, metrics.HeartBeatErr)
				h.logger.Debug("failed to create heartbeat messages", zap.Error(err))
				continue
			}

			h.metrics.AddWebSocketDataHandlerStatus(h.config.Name, metrics.HeartBeatSuccess)
			h.logger.Debug("sending heartbeat messages to data provider", zap.Int("num_msgs", len(msgs)))

			for _, msg := range msgs {
				if err := h.connHandler.Write(msg); err != nil {
					h.metrics.AddWebSocketConnectionStatus(h.config.Name, metrics.WriteErr)
					h.logger.Debug("failed to write heartbeat message", zap.String("message", string(msg)), zap.Error(err))
				} else {
					h.metrics.AddWebSocketConnectionStatus(h.config.Name, metrics.WriteSuccess)
					h.logger.Debug("heartbeat message sent", zap.String("message", string(msg)))
				}
			}
		}
	}
}

// recv is used to manage the connection to the data provider.
func (h *WebSocketQueryHandlerImpl[K, V]) recv(ctx context.Context, responseCh chan<- providertypes.GetResponse[K, V]) error {
	defer func() {
		if err := recover(); err != nil {
			h.logger.Error("panic occurred", zap.Any("err", err))
		}
	}()

	h.logger.Debug("starting recv", zap.Int("buffer_size", cap(responseCh)))

	// Track the number of read errors. If this exceeds the max read error count,
	// we will close the connection and return. Read errors can occur if the data
	// provider closes the connection or if the connection is interrupted.
	readErrCount := 0

	for {
		// Track the time it takes to receive a message from the data provider.
		now := time.Now().UTC()

		select {
		case <-ctx.Done():
			// Case 1: The context is cancelled. Close the connection and return.
			h.logger.Debug("context finished")
			if err := h.close(); err != nil {
				return err
			}

			return ctx.Err()
		default:
			// Case 2: The context is not cancelled. Wait for a message from the data provider.
			message, err := h.connHandler.Read()
			if err != nil {
				h.logger.Error(
					"failed to read message from websocket handler",
					zap.String("message", string(message)),
					zap.Error(err),
				)
				h.metrics.AddWebSocketConnectionStatus(h.config.Name, metrics.ReadErr)

				// If the read error count is greater than the max read error count, close the
				// connection and return.
				readErrCount++
				if readErrCount < h.config.MaxReadErrorCount {
					continue
				}

				h.logger.Error("max read errors reached", zap.Error(err))
				if err := h.close(); err != nil {
					return err
				}

				return errors.ErrReadWithErr(err)
			}

			h.logger.Debug("message received; attempting to handle message", zap.String("message", string(message)))
			h.metrics.AddWebSocketConnectionStatus(h.config.Name, metrics.ReadSuccess)

			// Handle the message.
			response, updateMessage, err := h.dataHandler.HandleMessage(message)
			if err != nil {
				h.logger.Debug("failed to handle websocket message", zap.Error(err))
				h.metrics.AddWebSocketDataHandlerStatus(h.config.Name, metrics.HandleMessageErr)
				continue
			}

			// Immediately send the response to the response channel. Even if this is
			// empty, it will be handled by the provider. Note that if the context has been
			// cancelled, we should not send the response to the channel. Otherwise, we risk
			// sending a response to a closed channel.
			select {
			case <-ctx.Done():
				h.logger.Debug("context finished")
				if err := h.close(); err != nil {
					return errors.ErrCloseWithErr(err)
				}

				return ctx.Err()
			case responseCh <- response:
				h.logger.Debug("handled message successfully; sent response to response channel", zap.String("response", response.String()))
				h.metrics.AddWebSocketDataHandlerStatus(h.config.Name, metrics.HandleMessageSuccess)
			}

			// If the update messages are not nil, send it to the data provider.
			if len(updateMessage) != 0 {
				for _, msg := range updateMessage {
					h.logger.Debug("sending update message to data provider", zap.String("update_message", string(msg)))
					h.metrics.AddWebSocketDataHandlerStatus(h.config.Name, metrics.CreateMessageSuccess)

					if err := h.connHandler.Write(msg); err != nil {
						h.logger.Error("failed to write update message", zap.Error(err))
						h.metrics.AddWebSocketConnectionStatus(h.config.Name, metrics.WriteErr)
					} else {
						h.logger.Debug("update message sent")
						h.metrics.AddWebSocketConnectionStatus(h.config.Name, metrics.WriteSuccess)
					}
				}
			}
		}

		// Record the time it took to receive the message.
		h.metrics.ObserveWebSocketLatency(h.config.Name, time.Since(now))
	}
}

// close is used to close the connection to the data provider.
func (h *WebSocketQueryHandlerImpl[K, V]) close() error {
	h.logger.Debug("closing connection to websocket handler")
	if err := h.connHandler.Close(); err != nil {
		h.logger.Debug("failed to close connection", zap.Error(err))
		h.metrics.AddWebSocketConnectionStatus(h.config.Name, metrics.CloseErr)
		return errors.ErrCloseWithErr(err)
	}

	h.logger.Debug("connection closed")
	h.metrics.AddWebSocketConnectionStatus(h.config.Name, metrics.CloseSuccess)
	return nil
}

// Copy is used to create a copy of the query handler. This is useful for creating
// multiple connections to the same data provider.
func (h *WebSocketQueryHandlerImpl[K, V]) Copy() WebSocketQueryHandler[K, V] {
	return &WebSocketQueryHandlerImpl[K, V]{
		logger:      h.logger,
		config:      h.config,
		dataHandler: h.dataHandler.Copy(),
		connHandler: h.connHandler.Copy(),
		metrics:     h.metrics,
	}
}
