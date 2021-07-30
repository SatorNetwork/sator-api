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

	"github.com/dmitrymomot/media-srv/repository"
	"github.com/dmitrymomot/media-srv/resizer"
	"github.com/dmitrymomot/media-srv/storage"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetOriginalItem(t *testing.T) {
	createdAt, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2006")
	item := repository.OriginalItem{
		ID:        uuid.New(),
		Name:      "image.png",
		Path:      "uploads/image.png",
		URL:       "http://test/uploads/image.png",
		CreatedAt: createdAt,
	}
	db, mock, err := repository.NewSQLMock()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repository.GetOriginalItemByIDMock(mock, item, nil)

	repo := repository.New(db)

	s3mock := &storage.S3Mock{}
	opt := storage.Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
	}
	stor := storage.New(s3mock, opt)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/origin/{oid}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("oid", item.ID.String())
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	t.Run("success", func(t *testing.T) {
		h := &Handler{
			db:      db,
			query:   repo,
			storage: stor,
			resize:  resizer.Resize,
		}
		if err := h.GetOriginalItem(w, r); err != nil {
			t.Errorf("Handler.GetOriginalItem() error = %v", err)
		}

		result := w.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.Equal(t, http.StatusOK, result.StatusCode)

		expected, _ := json.Marshal(item)
		assert.JSONEqf(t, string(expected), string(body), "response does not match to expected jsonn string")
	})
}

func TestHandler_GetOriginalItem_WrongID(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/origin/{oid}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("oid", "123")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	t.Run("wrong id", func(t *testing.T) {
		Wrap((&Handler{}).GetOriginalItem).ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}

func TestHandler_GetOriginalItem_NotFound(t *testing.T) {
	createdAt, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2006")
	item := repository.OriginalItem{
		ID:        uuid.New(),
		Name:      "image.png",
		Path:      "uploads/image.png",
		URL:       "http://test/uploads/image.png",
		CreatedAt: createdAt,
	}
	db, mock, err := repository.NewSQLMock()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repository.GetOriginalItemByIDMock(mock, item, sql.ErrNoRows)

	repo := repository.New(db)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/origin/{oid}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("oid", item.ID.String())
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	t.Run("not found", func(t *testing.T) {
		h := &Handler{
			db:      db,
			query:   repo,
			storage: nil,
		}
		Wrap(h.GetOriginalItem).ServeHTTP(w, r)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}
