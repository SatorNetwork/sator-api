package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dmitrymomot/media-srv/repository"
	"github.com/dmitrymomot/media-srv/resizer"
	"github.com/dmitrymomot/media-srv/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Upload(t *testing.T) {
	createdAt, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2006")
	oid := uuid.New()
	item := repository.OriginalItem{
		ID:        oid,
		Name:      "image.png",
		Path:      "uploads/image.png",
		URL:       "http://test/uploads/image.png",
		CreatedAt: createdAt,
	}
	rid := uuid.New()
	ritem := repository.ResizedItem{
		ID:        rid,
		OID:       oid,
		Name:      "image.png",
		Path:      "uploads/image-100x100.png",
		URL:       "http://test/uploads/image-100x100.png",
		Width:     100,
		Height:    100,
		CreatedAt: createdAt,
	}
	db, mock, err := repository.NewSQLMock()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	mock.ExpectBegin()
	repository.CreateOriginalItemMock(mock, repository.CreateOriginalItemParams{
		ID:   oid,
		Name: item.Name,
		Path: item.Path,
		URL:  item.URL,
	}, nil)
	mock.ExpectCommit()
	mock.ExpectBegin()
	repository.CreateResizedItemMock(mock, repository.CreateResizedItemParams{
		ID:     rid,
		OID:    oid,
		Name:   ritem.Name,
		Path:   ritem.Path,
		URL:    ritem.URL,
		Width:  ritem.Width,
		Height: ritem.Height,
	}, nil)
	mock.ExpectCommit()

	repo := repository.New(db)

	s3mock := &storage.S3Mock{}
	opt := storage.Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
	}
	stor := storage.New(s3mock, opt)

	fpath := "./testdata/image.png" //The path to upload the file
	file, err := os.Open(fpath)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filepath.Base(fpath))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	r, err := http.NewRequest("POST", "/origin/{oid}", body)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Content-Type", writer.FormDataContentType())
	data := url.Values{}
	data.Add("width", "100")
	data.Add("height", "100")
	r.Form = data
	w := httptest.NewRecorder()

	t.Run("success", func(t *testing.T) {
		h := &Handler{
			db:      db,
			query:   repo,
			storage: stor,
			resize:  resizer.Resize,
		}
		if err := h.Upload(w, r); err != nil {
			t.Errorf("Handler.Upload() error = %v", err)
		}

		result := w.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.Equal(t, http.StatusOK, result.StatusCode)

		expected, _ := json.Marshal(map[string]interface{}{
			"original": item,
			"resized":  ritem,
		})
		assert.JSONEqf(t, string(expected), string(body), "response does not match to expected jsonn string")
	})
}
