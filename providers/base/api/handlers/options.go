package handlers

// Option is a function that is used to configure a RequestHandler.
type Option func(*RequestHandlerImpl)

// WithHTTPMethod is an option that is used to set the HTTP method used to make requests.
func WithHTTPMethod(method string) Option {
	return func(r *RequestHandlerImpl) {
		r.method = method
	}
}

// WithHTTPHeaders is an option that is used to set the HTTP headers used to make requests.
func WithHTTPHeaders(headers map[string]string) Option {
	return func(r *RequestHandlerImpl) {
		r.headers = headers
	}
}
