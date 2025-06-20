package response

// Predefined error messages returned by the HTTP API.
const (
	// ErrInternalServer is returned when an unexpected server-side error occurs.
	ErrInternalServer = "internal server error"

	// ErrMethodNotAllowed is returned when an HTTP method is not supported for the requested endpoint.
	ErrMethodNotAllowed = "method not allowed"
)
