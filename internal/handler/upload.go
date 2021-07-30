package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"path"

	"github.com/SatorNetwork/sator-api/internal/storage"
	"github.com/SatorNetwork/sator-api/svc/mediaservice/repository"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thedevsaddam/govalidator"
)

// Upload new item and resize according to passed parameters
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) error {
	rules := govalidator.MapData{
		"file:image": []string{"required", "ext:png", "size:2097152", "mime:image/png"},
	}
	if err := validate(r, rules, nil); err != nil {
		return err
	}

	file, header, err := r.FormFile("image")
	defer file.Close()

	// Store original image
	tx, err := h.db.Begin()
	if err != nil {
		return errors.Wrap(err, "begin db transaction")
	}
	query := h.query.WithTx(tx)

	oid := uuid.New()
	fname := fmt.Sprintf("%s%s", oid.String(), path.Ext(header.Filename))
	ct := header.Header.Get("Content-Type")
	item, err := query.CreateItem(r.Context(), repository.CreateItemParams{
		ID:         oid,
		Filename:   header.Filename,
		Filepath:   h.storage.FileURL(h.storage.FilePath(fname)),
		RelationID: oid, // TODO: add relation data.
		RelationType: sql.NullString{ // TODO: add relation data.
			String: "",
			Valid:  true,
		},
	})
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "store original item to db")
	}
	if err := h.storage.Upload(file, item.Filepath, storage.Public, ct); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "upload original image")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit original item")
	}

	return jsonResponse(w, http.StatusOK, data{
		"item": item,
		//"resized":  resizedItem,
	})
}
