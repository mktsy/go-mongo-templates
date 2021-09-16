package reqreader

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ReadPathParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}

func ReadQueryParam(r *http.Request, param string) string {
	return r.URL.Query().Get(param)
}

func ReadBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func ReadHeader(r *http.Request, param string) string {
	return r.Header.Get(param)
}

func SetHeader(r *http.Request, param, value string) {
	r.Header.Set(param, value)
}
