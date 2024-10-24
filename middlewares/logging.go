package middlewares

import (
	"log"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time
		startTime := time.Now()

		// Log the incoming request details
		log.Printf("Started %s %s -- %s %s", r.Method, r.URL.Path, r.Host, w.Header())

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the request processing time
		elapsedTime := time.Since(startTime)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, elapsedTime)
	})
}
