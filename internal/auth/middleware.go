package auth

import (
	"context"
	"net/http"
	"slices"
	"strings"
)

func (tm *TokenManager) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			Unauthorize(w)
			return
		}

		tokenString, ok := strings.CutPrefix(token, "Bearer ")
		if !ok || tokenString == "" {
			Unauthorize(w)
			return
		}

		claims, err := tm.Verify(tokenString)
		if err != nil {
			Unauthorize(w)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (tm *TokenManager) Protect(h http.Handler, roles ...Role) http.Handler {
	// slices.Contains on an empty list is always false
	if len(roles) > 0 {
		h = RequireRole(roles...)(h)
	}
	return tm.Authenticate(h)
}

func RequireRole(roles ...Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				Unauthorize(w)
				return
			}

			if !slices.Contains(roles, claims.Role) {
				Forbidden(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Unauthorize(w http.ResponseWriter) {
	http.Error(w, "invalid credentials", http.StatusUnauthorized)
}

func Forbidden(w http.ResponseWriter) {
	http.Error(w, "invalid permission", http.StatusForbidden)
}
