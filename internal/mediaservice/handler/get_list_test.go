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

	"github.com/SatorNetwork/sator-api/internal/mediaservice/storage"
	"github.com/SatorNetwork/sator-api/svc/mediaservice/repository"

	"github.com/google/uuid"
)

func TestHandler_GetItemsList(t *testing.T) {
	createdAt, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2006")
	arg := repository.GetItemsListParams{
		Limit:  10,
		Offset: 0,
	}
	items := []repository.Item{
		{
			ID:       uuid.New(),
			Filename: "image.png",
			Filepath: "http://test/uploads/image.png",
			RelationType: sql.NullString{
				String: "RelationType",
				Valid:  true,
			},
			RelationID: uuid.New(),
			CreatedAt:  createdAt,
		},
		{
			ID:       uuid.New(),
			Filename: "image.png",
			Filepath: "http://test/uploads/image.png",
			RelationType: sql.NullString{
				String: "RelationType",
				Valid:  true,
			},
			RelationID: uuid.New(),
			CreatedAt:  createdAt,
		},
	}
	db, mock, err := repository.NewSQLMock()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repository.GetItemsListMock(mock, arg, items, nil)

	repo := repository.New(db)

	s3mock := &storage.S3Mock{}
	opt := storage.Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
	}
	stor := storage.New(s3mock, opt)

	r, err := http.NewRequest("GET", "/item?limit=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		db      *sql.DB
		query   *repository.Queries
		storage *storage.Interactor
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
	}{"success", fields{db, repo, stor}, args{httptest.NewRecorder(), r}, false}
	t.Run(tt.name, func(t *testing.T) {
		h := &Handler{
			db:      tt.fields.db,
			query:   tt.fields.query,
			storage: tt.fields.storage,
		}
		if err := h.GetItemsList(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
			t.Errorf("Handler.GetItemsList() error = %v, wantErr %v", err, tt.wantErr)
		}

		resp, ok := tt.args.w.(*httptest.ResponseRecorder)
		if !ok {
			t.Error("Wrong response recorder")
		}
		result := resp.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.Equal(t, result.StatusCode, http.StatusOK)

		expected, _ := json.Marshal(items)
		assert.JSONEqf(t, string(expected), string(body), "response does not match to expected json string")
	})
}

func TestHandler_GetItemsList_NotFound(t *testing.T) {
	arg := repository.GetItemsListParams{
		Limit:  10,
		Offset: 0,
	}
	db, mock, err := repository.NewSQLMock()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repository.GetItemsListMock(mock, arg, nil, sql.ErrNoRows)

	repo := repository.New(db)

	s3mock := &storage.S3Mock{}
	opt := storage.Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
	}
	stor := storage.New(s3mock, opt)

	r, err := http.NewRequest("GET", "/item?limit=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		db      *sql.DB
		query   *repository.Queries
		storage *storage.Interactor
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
	}{"not found", fields{db, repo, stor}, args{httptest.NewRecorder(), r}, true}
	t.Run(tt.name, func(t *testing.T) {
		h := &Handler{
			db:      tt.fields.db,
			query:   tt.fields.query,
			storage: tt.fields.storage,
		}
		if err := h.GetItemsList(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
			t.Errorf("Handler.GetItemsList() error = %v, wantErr %v", err, tt.wantErr)
		}
	})
}

func TestHandler_GetItemsList_LimitNotSet(t *testing.T) {
	createdAt, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2006")
	arg := repository.GetItemsListParams{
		Limit:  10,
		Offset: 0,
	}
	items := []repository.Item{
		{
			ID:       uuid.New(),
			Filename: "image.png",
			Filepath: "http://test/uploads/image.png",
			RelationType: sql.NullString{
				String: "RelationType",
				Valid:  true,
			},
			RelationID: uuid.New(),
			CreatedAt:  createdAt,
		},
		{
			ID:       uuid.New(),
			Filename: "image.png",
			Filepath: "http://test/uploads/image.png",
			RelationType: sql.NullString{
				String: "RelationType",
				Valid:  true,
			},
			RelationID: uuid.New(),
			CreatedAt:  createdAt,
		},
	}
	db, mock, err := repository.NewSQLMock()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repository.GetItemsListMock(mock, arg, items, nil)

	repo := repository.New(db)

	s3mock := &storage.S3Mock{}
	opt := storage.Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
	}
	stor := storage.New(s3mock, opt)

	r, err := http.NewRequest("GET", "/item", nil)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		db      *sql.DB
		query   *repository.Queries
		storage *storage.Interactor
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
	}{"limit not set", fields{db, repo, stor}, args{httptest.NewRecorder(), r}, false}
	t.Run(tt.name, func(t *testing.T) {
		h := &Handler{
			db:      tt.fields.db,
			query:   tt.fields.query,
			storage: tt.fields.storage,
		}
		if err := h.GetItemsList(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
			t.Errorf("Handler.GetItemsList() error = %v, wantErr %v", err, tt.wantErr)
		}

		resp, ok := tt.args.w.(*httptest.ResponseRecorder)
		if !ok {
			t.Error("Wrong response recorder")
		}
		result := resp.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.Equal(t, result.StatusCode, http.StatusOK)

		expected, _ := json.Marshal(items)
		assert.JSONEqf(t, string(expected), string(body), "response does not match to expected json string")
	})
}
