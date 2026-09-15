package teachers

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /teachers", h.List)
	mux.HandleFunc("GET /teachers/{id}", h.GetById)
	mux.HandleFunc("POST /teachers", h.Create)
	mux.HandleFunc("DELETE /teachers/{id}", h.Delete)
	mux.HandleFunc("PUT /teachers/{id}", h.Update)
}
