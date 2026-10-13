package teaching_assignments

import (
	"net/http"

	"github.com/aaksant/school-api/internal/auth"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, tm *auth.TokenManager) {
	permittedRoles := []auth.Role{auth.RoleExecutive, auth.RoleTeacher}

	mux.Handle("GET /teachers/{teacherId}/assignments", tm.Protect(
		http.HandlerFunc(h.List),
		permittedRoles...,
	))
	mux.Handle("GET /teachers/{teacherId}/assignments/{id}", tm.Protect(
		http.HandlerFunc(h.GetById),
		permittedRoles...,
	))
	mux.Handle("POST /teachers/{teacherId}/assignments", tm.Protect(
		http.HandlerFunc(h.Create),
		auth.RoleExecutive,
	))
	mux.Handle("DELETE /teachers/{teacherId}/assignments/{id}", tm.Protect(
		http.HandlerFunc(h.Delete),
		auth.RoleExecutive,
	))
}
