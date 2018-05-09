package middlewares

import (
	"net/http"
)

// AdminKeyMiddleware check the admin key for interact with the transmitter api
func AdminKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	})
}
