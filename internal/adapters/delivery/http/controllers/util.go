package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// ParseBody parses the body of an HTTP request and unmarshals it into the provided value.
//
// It takes a pointer to an http.Request and a value of any type as parameters.
// It returns an error.
func RequestParseBody(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func ResponseJson(w http.ResponseWriter, code int, body interface{}) {
	defaultHeaders(w)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}

func ResponseError(w http.ResponseWriter, code int, message string) {
	defaultHeaders(w)
	w.WriteHeader(code)
	body := map[string]string{"message": message, "status_code": strconv.Itoa(code)}
	json.NewEncoder(w).Encode(body)
}

func ResponseMessage(w http.ResponseWriter, code int, message *string) {
	defaultHeaders(w)
	body := map[string]*string{"message": message}
	json.NewEncoder(w).Encode(body)
}

func ResponseSendStatus(w http.ResponseWriter, code int) {
	defaultHeaders(w)
	w.WriteHeader(code)
}

func defaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Server", "hyperzoop 0.0.1")
}
