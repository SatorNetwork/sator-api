package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// GetOriginalItem http handler
func (h *Handler) GetOriginalItem(w http.ResponseWriter, r *http.Request) error {
	oid, err := uuid.Parse(chi.URLParam(r, "oid"))
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, "wrong id format")
	}
	item, err := h.query.GetItemByID(r.Context(), oid)
	if err != nil {
		return err
	}

	return jsonResponse(w, http.StatusOK, item)
}
