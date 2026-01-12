package middleware

import (
	"net/http"
	"time"

	"github.com/razedwell/go-hand/internal/platform/logger"
)

// LoggingMiddleware logs the details of each HTTP request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}
