package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SatorNetwork/sator-api/internal/mediaservice/storage"
	"github.com/SatorNetwork/sator-api/svc/mediaservice/repository"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetItem(t *testing.T) {
	createdAt, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2006")
	item := repository.Item{
		ID:       uuid.New(),
		Filename: "image.png",
		Filepath: "http://test/uploads/image.png",
		RelationType: sql.NullString{
			String: "RelationType",
			Valid:  true,
		},
		RelationID: uuid.New(),
		CreatedAt:  createdAt,
	}
	db, mock, err := repository.NewSQLMock()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repository.GetItemByIDMock(mock, item, nil)

	repo := repository.New(db)

	s3mock := &storage.S3Mock{}
	opt := storage.Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
	}
	stor := storage.New(s3mock, opt)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/item/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", item.ID.String())
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	t.Run("success", func(t *testing.T) {
		h := &Handler{
			db:      db,
			query:   repo,
			storage: stor,
		}
		if err := h.GetItem(w, r); err != nil {
			t.Errorf("Handler.GetItem() error = %v", err)
		}

		result := w.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.Equal(t, http.StatusOK, result.StatusCode)

		expected, _ := json.Marshal(item)
		assert.JSONEqf(t, string(expected), string(body), "response does not match to expected jsonn string")
	})
}

func TestHandler_GetItem_WrongID(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/item/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	t.Run("wrong id", func(t *testing.T) {
		Wrap((&Handler{}).GetItem).ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}

func TestHandler_GetItem_NotFound(t *testing.T) {
	createdAt, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2006")
	item := repository.Item{
		ID:       uuid.New(),
		Filename: "image.png",
		Filepath: "http://test/uploads/image.png",
		RelationType: sql.NullString{
			String: "RelationType",
			Valid:  true,
		},
		RelationID: uuid.New(),
		CreatedAt:  createdAt,
	}
	db, mock, err := repository.NewSQLMock()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repository.GetItemByIDMock(mock, item, sql.ErrNoRows)

	repo := repository.New(db)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/item/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", item.ID.String())
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	t.Run("not found", func(t *testing.T) {
		h := &Handler{
			db:      db,
			query:   repo,
			storage: nil,
		}
		Wrap(h.GetItem).ServeHTTP(w, r)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}
