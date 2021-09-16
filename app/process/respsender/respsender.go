package respsender

import (
	"encoding/json"
	"io"
	"net/http"
)

func ResponseString(w http.ResponseWriter, body string, status int) {
	w.WriteHeader(status)
	io.WriteString(w, body)
}

func ResponseMap(w http.ResponseWriter, body interface{}, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
