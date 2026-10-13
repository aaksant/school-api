package users

import (
	"net/http"

	"github.com/aaksant/school-api/internal/auth"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, tm *auth.TokenManager) {
	mux.HandleFunc("POST /users/login", h.Login)
	mux.HandleFunc("POST /users/refresh", h.Refresh)
	mux.Handle("POST /users/register", tm.Protect(
		http.HandlerFunc(h.Register),
		auth.RoleExecutive,
	))
	mux.Handle("POST /users/signout", tm.Protect(http.HandlerFunc(h.Signout)))
}
