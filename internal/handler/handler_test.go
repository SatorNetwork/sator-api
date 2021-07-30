package handler

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/dmitrymomot/media-srv/repository"
	"github.com/dmitrymomot/media-srv/resizer"
	"github.com/dmitrymomot/media-srv/storage"
	"gotest.tools/assert"
)

func TestValidationError_Add(t *testing.T) {
	type fields struct {
		errBag url.Values
	}
	type args struct {
		key  string
		vals []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"success", fields{}, args{"name", []string{"error"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ValidationError{
				errBag: tt.fields.errBag,
			}
			e.Add(tt.args.key, tt.args.vals)
			if _, ok := e.errBag["name"]; !ok {
				t.Errorf("ValidationError.Add(): key %s does not exist", tt.args.key)
			}
		})
	}
}

func TestValidationError_Set(t *testing.T) {
	errBag := url.Values{"name": []string{"error"}}
	type fields struct {
		errBag url.Values
	}
	type args struct {
		errBag url.Values
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"success", fields{}, args{errBag}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ValidationError{
				errBag: tt.fields.errBag,
			}
			e.Set(tt.args.errBag)
			if !reflect.DeepEqual(e.errBag, errBag) {
				t.Errorf("ValidationError.Set() = %v, want %v", e.errBag, errBag)
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	errBag := url.Values{"name": []string{"error"}}
	type fields struct {
		errBag url.Values
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"success", fields{errBag}, "validationError: map[name:[error]]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ValidationError{
				errBag: tt.fields.errBag,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ValidationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPError_Error(t *testing.T) {
	type fields struct {
		Code    int
		Message interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"success", fields{http.StatusBadRequest, "test"}, "code: 400; message: test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := HTTPError{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("HTTPError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrap_ServeHTTP(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		h       Wrap
		args    args
		code    int
		body    string
		wantErr bool
	}{
		{"success", Wrap(func(w http.ResponseWriter, r *http.Request) error {
			return stringResponse(w, http.StatusOK, "test")
		}), args{httptest.NewRecorder(), r}, http.StatusOK, "test", false},

		{"validation error", Wrap(func(w http.ResponseWriter, r *http.Request) error {
			err := &ValidationError{}
			err.Add("test", []string{"error"})
			return err
		}), args{httptest.NewRecorder(), r}, http.StatusUnprocessableEntity, "{\"validationError\":{\"test\":[\"error\"]}}\n", true},

		{"http error", Wrap(func(w http.ResponseWriter, r *http.Request) error {
			return NewHTTPError(http.StatusBadRequest, "error")
		}), args{httptest.NewRecorder(), r}, http.StatusBadRequest, "{\"error\":\"error\"}\n", true},

		{"undefined error", Wrap(func(w http.ResponseWriter, r *http.Request) error {
			return errors.New("error")
		}), args{httptest.NewRecorder(), r}, http.StatusInternalServerError, "{\"error\":\"Internal Server Error\"}\n", true},

		{"ont founnd error", Wrap(func(w http.ResponseWriter, r *http.Request) error {
			return sql.ErrNoRows
		}), args{httptest.NewRecorder(), r}, http.StatusNotFound, "{\"error\":\"Not Found\"}\n", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.ServeHTTP(tt.args.w, tt.args.r)

			resp, ok := tt.args.w.(*httptest.ResponseRecorder)
			if !ok {
				t.Error("Wrong response recorder")
			}
			result := resp.Result()
			body, _ := ioutil.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, tt.code)
			assert.Equal(t, string(body), tt.body)
		})

	}
}

func TestNew(t *testing.T) {
	db, _, err := repository.NewSQLMock()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := repository.New(db)

	s3mock := &storage.S3Mock{}
	opt := storage.Options{
		Bucket:         "test",
		URL:            "http://test.storage",
		ForcePathStyle: false,
	}
	stor := storage.New(s3mock, opt)

	h := &Handler{db: db, query: repo, storage: stor, resize: resizer.Resize}

	type args struct {
		db *sql.DB
		q  *repository.Queries
		s  *storage.Interactor
		rf resizerFunc
	}
	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{"success", args{db, repo, stor, resizer.Resize}, h},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.db, tt.args.q, tt.args.s, tt.args.rf); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("New() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
