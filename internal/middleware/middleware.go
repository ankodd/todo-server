// Package middleware provides middleware functions'
package middleware

import (
	"log"
	"net/http"
)

// Logger is a middleware that logs the request and response
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Addr:%s. Method: %s. URL:%s.\n", r.RemoteAddr, r.Method, r.URL)
		next(w, r)
	}
}

// AcceptCORS is a middleware that sets the CORS headers
func AcceptCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		next(w, r)
	}
}

// SetJSONContentType is a middleware that sets the content type to JSON
func SetJSONContentType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// All is a middleware that sets the CORS headers and sets the content type to JSON
func All(next http.HandlerFunc) http.HandlerFunc {
	return AcceptCORS(SetJSONContentType(Logger(next)))
}
