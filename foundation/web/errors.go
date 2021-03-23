package web

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

// Error is used to add web info to a request error
type Error struct {
	Err    error
	Status int
	Fields []FieldError
}

// NewRequestError wraps a provided error with an HTTP status code. This function should be used when handlers encounter expected errors.
func NewRequestError(err error, status int) error {
	return &Error{Err: err, Status: status}
}

// Error implements the error interface. It uses the default message of the wrapped error. This is what will be shown in the services' logs.
func (e *Error) Error() string {
	return e.Err.Error()
}

// shutdown is a type used to help with the graceful termination of the service.
type shutdown struct {
	Message string
}

// Error is the implementation of the error interface.
func (s *shutdown) Error() string {
	return s.Message
}

// NewShutdownError returns an error that causes the framework to signal
// a graceful shutdown.
func NewShutdownError(message string) error {
	return &shutdown{message}
}
