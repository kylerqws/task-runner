package response

import "net/http"

// RespondNoContent sends an HTTP response with the provided status code and no response body.
// Despite the name, this function allows custom status codes (e.g., 204, 202) for empty responses.
func RespondNoContent(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}
