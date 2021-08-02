package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// GetItem http handler
func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, "wrong id format")
	}
	item, err := h.query.GetItemByID(r.Context(), id)
	if err != nil {
		return err
	}

	return jsonResponse(w, http.StatusOK, item)
}
