package handler

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Router returns http.Handler interface
func Router(h *Handler) http.Handler {
	r := chi.NewRouter()

	r.Method(http.MethodPost, "/item", Wrap(h.Upload))
	r.Method(http.MethodGet, "/item", Wrap(h.GetItemsList))
	r.Method(http.MethodGet, "/item/{id}", Wrap(h.GetItem))

	return r
}
