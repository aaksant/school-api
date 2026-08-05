package middlewares

import (
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture the status code
// and inject a response-time header right before headers are sent.

type responseWriter struct {
	http.ResponseWriter
	statusCode  int
	start       time.Time
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.wroteHeader {
		rw.Header().Set("X-Response-Time", time.Since(rw.start).String())
		rw.statusCode = statusCode
		rw.wroteHeader = true
	}
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write ensures WriteHeader runs even if the handler never calls WriteHeader explicitly.
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func ResponseTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Wrap the ResponseWriter so we can capture the status code
		// and set the X-Response-Time header at the last possible moment.
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			start:          start,
		}

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)
		log.Printf(
			"POST %v. Request took %v with status %d\n",
			r.URL, duration, wrappedWriter.statusCode,
		)
	})
}
