package students

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /students", h.List)
	mux.HandleFunc("GET /students/{id}", h.GetById)
	mux.HandleFunc("POST /students", h.Create)
	mux.HandleFunc("DELETE /students/{id}", h.Delete)
	mux.HandleFunc("PUT /students/{id}", h.Update)
	mux.HandleFunc("PATCH /students/{id}", h.PartialUpdate)
}
