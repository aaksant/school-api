package students

import (
	"net/http"

	"github.com/aaksant/school-api/internal/auth"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, tm *auth.TokenManager) {
	read := []auth.Role{auth.RoleExecutive, auth.RoleTeacher}

	mux.Handle("GET /students", tm.Protect(
		http.HandlerFunc(h.List),
		read...,
	))
	mux.Handle("GET /students/{id}", tm.Protect(
		http.HandlerFunc(h.GetById),
		read...,
	))
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
