package teachingassignments

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /teachers/{teacherId}/assignments", h.List)
	mux.HandleFunc("GET /teachers/{teacherId}/assignments/{id}", h.GetById)
	mux.HandleFunc("POST /teachers/{teacherId}/assignments", h.Create)
	mux.HandleFunc("DELETE /teachers/{teacherId}/assignments/{id}", h.Delete)
}
