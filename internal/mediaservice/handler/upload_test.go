package handler

import (
	"bytes"
	"database/sql"
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

	"github.com/SatorNetwork/sator-api/internal/mediaservice/storage"
	"github.com/SatorNetwork/sator-api/svc/mediaservice/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Upload(t *testing.T) {
	createdAt, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2006")
	id := uuid.New()
	item := repository.Item{
		ID:       id,
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
	mock.ExpectBegin()
	repository.CreateItemMock(mock, repository.CreateItemParams{
		ID:           id,
		Filename:     item.Filename,
		Filepath:     item.Filepath,
		RelationType: item.RelationType,
		RelationID:   item.RelationID,
	}, nil)
	mock.ExpectCommit()

	repo := repository.New(db)

	s3mock := &storage.S3Mock{}
	opt := storage.Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
	}
	storage := storage.New(s3mock, opt)

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

	r, err := http.NewRequest("POST", "/item/{id}", body)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Content-Type", writer.FormDataContentType())
	r.Header.Set("RelationID", uuid.New().String())
	r.Header.Set("RelationType", "RelationType")
	data := url.Values{}
	//data.Add("width", "100")
	//data.Add("height", "100")
	r.Form = data
	w := httptest.NewRecorder()

	t.Run("success", func(t *testing.T) {
		h := &Handler{
			db:      db,
			query:   repo,
			storage: storage,
		}
		if err := h.Upload(w, r); err != nil {
			t.Errorf("Handler.Upload() error = %v", err)
		}

		result := w.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.Equal(t, http.StatusOK, result.StatusCode)

		expected, _ := json.Marshal(map[string]interface{}{
			"item": item,
		})
		assert.JSONEqf(t, string(expected), string(body), "response does not match to expected json string")
	})
}
