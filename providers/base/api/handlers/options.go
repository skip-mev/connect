package handlers

// Option is a function that is used to configure a RequestHandler.
type Option func(*RequestHandlerImpl)

// WithHTTPMethod is an option that is used to set the HTTP method used to make requests.
func WithHTTPMethod(method string) Option {
	return func(r *RequestHandlerImpl) {
		r.method = method
	}
}

// WithJSONHeader is an option that's used to set the HTTP headers in accordance with standard JSON-RPC
// fields.
func WithJSONHeader() Option {
	return func(r *RequestHandlerImpl) {
		r.requestHeaderPairs = map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		}
	}
}
