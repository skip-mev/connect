package metrics

type (
	// ConnectionStatus is a type that represents the status of a connection.
	ConnectionStatus int

	// HandlerStatus is a type that represents the status of a data handler.
	HandlerStatus int
)

const (
	// DialErr indicates that the provider could not establish a connection.
	DialErr ConnectionStatus = iota
	// DialSuccess indicates that the provider successfully established a connection.
	DialSuccess
	// WriteErr indicates that the provider could not write the message to the data provider.
	WriteErr
	// WriteSuccess indicates that the provider successfully wrote the message to the data provider.
	WriteSuccess
	// ReadErr indicates that the provider could not read the message from the data provider.
	ReadErr
	// ReadSuccess indicates that the provider successfully read the message from the data provider.
	ReadSuccess
	// CloseErr indicates that the provider could not close the connection.
	CloseErr
	// CloseSuccess indicates that the provider successfully closed the connection.
	CloseSuccess
	// Healthy indicates that the provider is healthy.
	Healthy
	// Unhealthy indicates that the provider is unhealthy.
	Unhealthy
)

const (
	// CreateMessageErr indicates that the provider could not construct a valid message to send
	// to the data provider.
	CreateMessageErr HandlerStatus = iota
	// CreateMessageSuccess indicates that the provider successfully constructed a valid message
	// to send to the data provider.
	CreateMessageSuccess
	// HandleMessageErr indicates that the provider could not handle the message from the data
	// provider.
	HandleMessageErr
	// HandleMessageSuccess indicates that the provider successfully handled the message from the
	// data provider.
	HandleMessageSuccess
	// HeartBeatSuccess indicates that the provider successfully constructed a heartbeat message
	// to send to the data provider.
	HeartBeatSuccess
	// HeartBeatErr indicates that the provider could not construct a heartbeat message to send
	// to the data provider.
	HeartBeatErr
	// Unknown indicates that the provider encountered an unknown error.
	Unknown
)

// String returns a string representation of the connection status.
func (s ConnectionStatus) String() string {
	switch s {
	case DialErr:
		return "dial_err"
	case DialSuccess:
		return "dial_success"
	case WriteErr:
		return "write_err"
	case WriteSuccess:
		return "write_success"
	case ReadErr:
		return "read_err"
	case ReadSuccess:
		return "read_success"
	case CloseErr:
		return "close_err"
	case CloseSuccess:
		return "close_success"
	case Healthy:
		return "healthy"
	case Unhealthy:
		return "unhealthy"
	default:
		return "unknown_status"
	}
}

// String returns a string representation of the handler status.
func (s HandlerStatus) String() string {
	switch s {
	case CreateMessageErr:
		return "create_message_err"
	case CreateMessageSuccess:
		return "create_message_success"
	case HandleMessageErr:
		return "handle_message_err"
	case HandleMessageSuccess:
		return "handle_message_success"
	case HeartBeatSuccess:
		return "heartbeat_success"
	case HeartBeatErr:
		return "heartbeat_err"
	default:
		return "unknown_err"
	}
}
