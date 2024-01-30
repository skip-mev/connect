package errors

import "errors"

var (
	// ErrHandleMessage is returned when the WebSocketDataHandler cannot handle a
	// message. This can occur if the handler is not configured to handle the given
	// message. Handlers must be able to handle heartbeat messages.
	ErrHandleMessage = errors.New("websocket data handler failed to handle message")

	// ErrCreateMessages is returned when the WebSocketDataHandler cannot create a
	// subscription messages. This can occur if the handler is not configured to
	// handle the given ids.
	ErrCreateMessages = errors.New("websocket data handler failed to create messages")

	// ErrRead is returned when the WebSocketConnHandler cannot read a message.
	ErrRead = errors.New("websocket connection handler failed to read message")

	// ErrWrite is returned when the WebSocketConnHandler cannot write a message.
	ErrWrite = errors.New("websocket connection handler failed to write message")

	// ErrClose is returned when the WebSocketConnHandler cannot close the connection.
	ErrClose = errors.New("websocket connection handler failed to close connection")

	// ErrDial is returned when the WebSocketConnHandler cannot create a connection.
	ErrDial = errors.New("websocket connection handler failed to create connection")
)

// ErrHandleMessageWithErr is used to create a new ErrHandleMessage with the given error.
// Provider's that implement the WebSocketDataHandler interface should use this function to
// create the error.
func ErrHandleMessageWithErr(err error) error {
	return errors.Join(ErrHandleMessage, err)
}

// ErrCreateMessageWithErr is used to create a new ErrCreateMessages with the given error.
// Provider's that implement the WebSocketDataHandler interface should use this function to
// create the error.
func ErrCreateMessageWithErr(err error) error {
	return errors.Join(ErrCreateMessages, err)
}

// ErrReadWithErr is used to create a new ErrRead with the given error.
// Provider's that implement the WebSocketConnHandler interface should use this function to
// create the error.
func ErrReadWithErr(err error) error {
	return errors.Join(ErrRead, err)
}

// ErrWriteWithErr is used to create a new ErrWrite with the given error.
// Provider's that implement the WebSocketConnHandler interface should use this function to
// create the error.
func ErrWriteWithErr(err error) error {
	return errors.Join(ErrWrite, err)
}

// ErrCloseWithErr is used to create a new ErrClose with the given error.
// Provider's that implement the WebSocketConnHandler interface should use this function to
// create the error.
func ErrCloseWithErr(err error) error {
	return errors.Join(ErrClose, err)
}

// ErrDialWithErr is used to create a new ErrCreate with the given error.
// Provider's that implement the WebSocketConnHandler interface should use this function to
// create the error.
func ErrDialWithErr(err error) error {
	return errors.Join(ErrDial, err)
}
