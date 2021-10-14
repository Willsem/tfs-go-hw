package response

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func Respond(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Status", strconv.Itoa(statusCode))
	w.Header().Set("Content-Type", "application/json")

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
