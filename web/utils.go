package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	responseHeaderContentTypeKey      = "Content-Type"
	responseHeaderContentTypeJSONUTF8 = "application/json; charset=UTF-8"
	responseHeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"
	responseHeaderAccessControlAllowMethods = "Access-Control-Allow-Methods"
	responseHeaderAccessControlAllowHeaders = "Access-Control-Allow-Headers"

	resourceNotFound = "Resource not found"
	errorMsg         = "Error"
)

// SendJSONWithHTTPCode outputs JSON with an HTTP code
func SendJSONWithHTTPCode(w http.ResponseWriter, d interface{}, code int) {
	w.Header().Set(responseHeaderContentTypeKey, responseHeaderContentTypeJSONUTF8)
	w.Header().Set(responseHeaderAccessControlAllowOrigin, "*")
	w.Header().Set(responseHeaderAccessControlAllowMethods, "*")
	w.Header().Set(responseHeaderAccessControlAllowHeaders, "*")
	w.WriteHeader(code)
	if d != nil {
		err := json.NewEncoder(w).Encode(d)
		if err != nil {
			panic(err)
		}
	}
}

// SendJSONOk outputs a JSON with http.StatusOK code
func SendJSONOk(w http.ResponseWriter, d interface{}) {
	SendJSONWithHTTPCode(w, d, http.StatusOK)
}

// SendJSONError sends error with a custom message and error code
func SendJSONError(w http.ResponseWriter, error string, code int) {
	SendJSONWithHTTPCode(w, map[string]string{errorMsg: error}, code)
}

// SendJSONNotFound produces a http.StatusNotFound response with the following JSON, '{"Error":"Resource not found"}'
func SendJSONNotFound(w http.ResponseWriter) {
	SendJSONError(w, resourceNotFound, http.StatusNotFound)
}

// NotFoundHandler return a JSON implementation of the not found handler
func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SendJSONNotFound(w)
	}
}

// ParamAsString returns an URL parameter /{name} as a string
func ParamAsString(name string, r *http.Request) string {
	vars := mux.Vars(r)
	value := vars[name]
	return value
}

// GetJSONContent returns the JSON content of a request
func GetJSONContent(v interface{}, r *http.Request) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
