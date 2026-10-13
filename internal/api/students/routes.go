package students

import (
	"net/http"

	"github.com/aaksant/school-api/internal/auth"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, tm *auth.TokenManager) {
	mux.HandleFunc("GET /students", h.List)
	mux.HandleFunc("GET /students/{id}", h.GetById)
	mux.Handle("POST /students", tm.Protect(
		http.HandlerFunc(h.Create),
		auth.RoleExecutive,
	))
	mux.Handle("DELETE /students/{id}", tm.Protect(
		http.HandlerFunc(h.Delete),
		auth.RoleExecutive,
	))
	mux.Handle("PUT /students/{id}", tm.Protect(
		http.HandlerFunc(h.Update),
		auth.RoleExecutive,
	))
	mux.Handle("PATCH /students/{id}", tm.Protect(
		http.HandlerFunc(h.PartialUpdate),
		auth.RoleExecutive,
	))
}
