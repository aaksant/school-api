package middlewares

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (gzw *gzipResponseWriter) Write(b []byte) (int, error) {
	return gzw.writer.Write(b)
}

func Compression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the client supports compression (via the Accept-Encoding request header).
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// Client doesn't support, next middleware
			next.ServeHTTP(w, r)
			return
		}

		// Set the Content-Encoding: gzip response header so the client knows to decompress it.
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		// If yes, wrap ResponseWriter.Write() so anything written
		// gets piped through a gzip.Writer instead of going straight to the connection.
		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzw := &gzipResponseWriter{
			ResponseWriter: w,
			writer:         gz,
		}

		next.ServeHTTP(gzw, r)
	})
}
