package dino

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
	err     error  `json:"-"`
	log     bool   `json:"-"`
}

type ErrorOption func(*Error)

func (e *Error) Error() string {
	return e.Message
}

func NewError(code int, message string, opts ...ErrorOption) *Error {
	err := &Error{
		Code:    code,
		Message: message,
		log:     true,
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

func WithDetails(details any) ErrorOption {
	return func(err *Error) {
		err.Details = details
	}
}

func WithInternalError(internalErr error) ErrorOption {
	return func(err *Error) {
		err.err = internalErr
	}
}

func WithoutLog() ErrorOption {
	return func(err *Error) {
		err.log = false
	}
}
