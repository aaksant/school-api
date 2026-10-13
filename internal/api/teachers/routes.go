package teachers

import (
	"net/http"

	"github.com/aaksant/school-api/internal/auth"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, tm *auth.TokenManager) {
	mux.HandleFunc("GET /teachers", h.List)
	mux.HandleFunc("GET /teachers/{id}", h.GetById)
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
