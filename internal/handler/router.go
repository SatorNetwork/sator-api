package handler

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Router returns http.Handler interface
func Router(h *Handler) http.Handler {
	r := chi.NewRouter()

	r.Method(http.MethodPost, "/original", Wrap(h.Upload))
	r.Method(http.MethodGet, "/original", Wrap(h.GetOriginalItemsList))
	r.Method(http.MethodGet, "/original/{oid}", Wrap(h.GetOriginalItem))

	return r
}
