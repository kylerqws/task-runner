package response

import (
	"encoding/json"
	"net/http"
)

// RespondJSON sends a JSON response with the given status code and data.
// It sets the "Content-Type" header to "application/json".
func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}
