package middlewares

import (
	"net/http"
	"slices"
)

// Add frontend origins later
var allowedOrigins []string = []string{
	"https://localhost:3000",
	"https://localhost:3001", // Simulate frontend server
}

func isOriginAllowed(origin string) bool {
	return slices.Contains(allowedOrigins, origin)
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			http.Error(w, "Origin not allowed", http.StatusForbidden)
			return
		}

		if r.Method == http.MethodOptions {
			return
		}

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		next.ServeHTTP(w, r)
	})
}
