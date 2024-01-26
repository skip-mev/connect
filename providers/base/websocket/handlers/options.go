package handlers

// Option is a function that is used to configure a WebSocketConnHandler.
type Option func(*WebSocketConnHandlerImpl)

// WithPreDialHook is an option that is used to set a pre-dial hook for a websocket connection.
func WithPreDialHook(hook PreDialHook) Option {
	return func(r *WebSocketConnHandlerImpl) {
		if hook == nil {
			panic("pre-dial hook cannot be nil")
		}

		r.preDialHook = hook
	}
}
