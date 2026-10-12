package users

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /users/register", h.Register)
	mux.HandleFunc("POST /users/login", h.Login)
	mux.HandleFunc("POST /users/refresh", h.Refresh)
	mux.Handle("POST /users/signout", h.tokenManager.Authenticate(
		http.HandlerFunc(h.Signout)),
	)
}
