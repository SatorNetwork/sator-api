package handler

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dmitrymomot/media-srv/repository"
	"github.com/dmitrymomot/media-srv/resizer"
	"github.com/dmitrymomot/media-srv/storage"
	"github.com/google/uuid"
)

func TestHandler_GetOriginalItemsList(t *testing.T) {
	createdAt, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2006")
	arg := repository.GetOriginalItemsListParams{
		Limit:  10,
		Offset: 0,
	}
	items := []repository.OriginalItem{
		{
			ID:        uuid.New(),
			Name:      "image.png",
			Path:      "uploads/image.png",
			URL:       "http://test/uploads/image.png",
			CreatedAt: createdAt,
		},
		{
			ID:        uuid.New(),
			Name:      "image.png",
			Path:      "uploads/image.png",
			URL:       "http://test/uploads/image.png",
			CreatedAt: createdAt,
		},
	}
	db, mock, err := repository.NewSQLMock()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repository.GetOriginalItemsListMock(mock, arg, items, nil)

	repo := repository.New(db)

	s3mock := &storage.S3Mock{}
	opt := storage.Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
	}
	stor := storage.New(s3mock, opt)

	r, err := http.NewRequest("GET", "/origin?limit=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		db      *sql.DB
		query   *repository.Queries
		storage *storage.Interactor
		resize  resizerFunc
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tt := struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{"success", fields{db, repo, stor, resizer.Resize}, args{httptest.NewRecorder(), r}, false}
	t.Run(tt.name, func(t *testing.T) {
		h := &Handler{
			db:      tt.fields.db,
			query:   tt.fields.query,
			storage: tt.fields.storage,
			resize:  tt.fields.resize,
		}
		if err := h.GetOriginalItemsList(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
			t.Errorf("Handler.GetOriginalItemsList() error = %v, wantErr %v", err, tt.wantErr)
		}

		resp, ok := tt.args.w.(*httptest.ResponseRecorder)
		if !ok {
			t.Error("Wrong response recorder")
		}
		result := resp.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.Equal(t, result.StatusCode, http.StatusOK)

		expected, _ := json.Marshal(items)
		assert.JSONEqf(t, string(expected), string(body), "response does not match to expected jsonn string")
	})
}

func TestHandler_GetOriginalItemsList_NotFound(t *testing.T) {
	arg := repository.GetOriginalItemsListParams{
		Limit:  10,
		Offset: 0,
	}
	db, mock, err := repository.NewSQLMock()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repository.GetOriginalItemsListMock(mock, arg, nil, sql.ErrNoRows)

	repo := repository.New(db)

	s3mock := &storage.S3Mock{}
	opt := storage.Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
	}
	stor := storage.New(s3mock, opt)

	r, err := http.NewRequest("GET", "/origin?limit=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		db      *sql.DB
		query   *repository.Queries
		storage *storage.Interactor
		resize  resizerFunc
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tt := struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{"not found", fields{db, repo, stor, resizer.Resize}, args{httptest.NewRecorder(), r}, true}
	t.Run(tt.name, func(t *testing.T) {
		h := &Handler{
			db:      tt.fields.db,
			query:   tt.fields.query,
			storage: tt.fields.storage,
			resize:  tt.fields.resize,
		}
		if err := h.GetOriginalItemsList(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
			t.Errorf("Handler.GetOriginalItemsList() error = %v, wantErr %v", err, tt.wantErr)
		}
	})
}

func TestHandler_GetOriginalItemsList_LimitNotSet(t *testing.T) {
	createdAt, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2006")
	arg := repository.GetOriginalItemsListParams{
		Limit:  10,
		Offset: 0,
	}
	items := []repository.OriginalItem{
		{
			ID:        uuid.New(),
			Name:      "image.png",
			Path:      "uploads/image.png",
			URL:       "http://test/uploads/image.png",
			CreatedAt: createdAt,
		},
		{
			ID:        uuid.New(),
			Name:      "image.png",
			Path:      "uploads/image.png",
			URL:       "http://test/uploads/image.png",
			CreatedAt: createdAt,
		},
	}
	db, mock, err := repository.NewSQLMock()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repository.GetOriginalItemsListMock(mock, arg, items, nil)

	repo := repository.New(db)

	s3mock := &storage.S3Mock{}
	opt := storage.Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
	}
	stor := storage.New(s3mock, opt)

	r, err := http.NewRequest("GET", "/origin", nil)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		db      *sql.DB
		query   *repository.Queries
		storage *storage.Interactor
		resize  resizerFunc
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tt := struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{"limit not set", fields{db, repo, stor, resizer.Resize}, args{httptest.NewRecorder(), r}, false}
	t.Run(tt.name, func(t *testing.T) {
		h := &Handler{
			db:      tt.fields.db,
			query:   tt.fields.query,
			storage: tt.fields.storage,
			resize:  tt.fields.resize,
		}
		if err := h.GetOriginalItemsList(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
			t.Errorf("Handler.GetOriginalItemsList() error = %v, wantErr %v", err, tt.wantErr)
		}

		resp, ok := tt.args.w.(*httptest.ResponseRecorder)
		if !ok {
			t.Error("Wrong response recorder")
		}
		result := resp.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.Equal(t, result.StatusCode, http.StatusOK)

		expected, _ := json.Marshal(items)
		assert.JSONEqf(t, string(expected), string(body), "response does not match to expected jsonn string")
	})
}
