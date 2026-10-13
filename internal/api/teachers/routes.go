package teachers

import (
	"net/http"

	"github.com/aaksant/school-api/internal/auth"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, tm *auth.TokenManager) {
	mux.Handle("GET /teachers", tm.Protect(
		http.HandlerFunc(h.List),
		auth.RoleExecutive,
	))
	mux.Handle("GET /teachers/{id}", tm.Protect(
		http.HandlerFunc(h.GetById),
		auth.RoleExecutive,
	))
	mux.Handle("POST /teachers", tm.Protect(
		http.HandlerFunc(h.Create),
		auth.RoleExecutive,
	))
	mux.Handle("DELETE /teachers/{id}", tm.Protect(
		http.HandlerFunc(h.Delete),
		auth.RoleExecutive,
	))
	mux.Handle("PUT /teachers/{id}", tm.Protect(
		http.HandlerFunc(h.Update),
		auth.RoleExecutive,
	))
}
