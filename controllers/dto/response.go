package dto

// Response wraps service results so controllers know status, message, and body.
type Response[T any] struct {
	Success    bool
	Message    string
	StatusCode int
	Body       *T
}

func NewResponse[T any](success bool, message string, status int, body *T) *Response[T] {
	return &Response[T]{
		Success:    success,
		Message:    message,
		StatusCode: status,
		Body:       body,
	}
}

func NewSuccessResponse[T any](body T, status int) *Response[T] {
	return &Response[T]{
		Success:    true,
		Message:    "ok",
		StatusCode: status,
		Body:       &body,
	}
}
