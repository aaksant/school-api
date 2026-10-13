package classes

import (
	"net/http"

	"github.com/aaksant/school-api/internal/auth"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, tm *auth.TokenManager) {
	mux.HandleFunc("GET /classes", h.List)
	mux.HandleFunc("GET /classes/{id}", h.GetById)
	mux.Handle("POST /classes", tm.Protect(
		http.HandlerFunc(h.Create),
		auth.RoleExecutive,
	))
	mux.Handle("DELETE /classes/{id}", tm.Protect(
		http.HandlerFunc(h.Delete),
		auth.RoleExecutive,
	))
	mux.Handle("PUT /classes/{id}", tm.Protect(
		http.HandlerFunc(h.Update),
		auth.RoleExecutive,
	))
}
