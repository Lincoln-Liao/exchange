package http

import "net/http"

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}
