package classes

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /classes", h.List)
	mux.HandleFunc("GET /classes/{id}", h.GetById)
	mux.HandleFunc("POST /classes", h.Create)
	mux.HandleFunc("DELETE /classes/{id}", h.Delete)
	mux.HandleFunc("PUT /classes/{id}", h.Update)
}
