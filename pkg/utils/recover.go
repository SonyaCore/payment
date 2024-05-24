package utils

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
)

// RecoverHandler is a middleware handler function that can be used
// to recover from unexpected panics, log a stack trace and respond with a generic
// 500 Internal Server Error.
// Ensures no sensitive data is leaked during panics.
func RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				debug.PrintStack()
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
