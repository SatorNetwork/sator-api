package handler

import (
	"database/sql"
	"fmt"
	"github.com/SatorNetwork/sator-api/svc/mediaservice/repository"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/SatorNetwork/sator-api/internal/storage"

	"github.com/pkg/errors"
)

type (
	// Wrap function makes a handler compatible with the default go mux interface
	Wrap func(w http.ResponseWriter, r *http.Request) error

	// ValidationError ...
	ValidationError struct {
		errBag url.Values
	}

	// HTTPError struct
	HTTPError struct {
		Code    int         `json:"code"`
		Message interface{} `json:"message"`
	}

	resizerFunc func(f io.ReadCloser, w, h int) (io.ReadSeeker, error)

	// Handler struct
	Handler struct {
		db      *sql.DB
		query   *repository.Queries
		storage *storage.Interactor
		resize  resizerFunc
	}
)

// NewHTTPError factory
func NewHTTPError(code int, msg interface{}) error {
	return HTTPError{Code: code, Message: msg}
}

// NewValidationError factory
func NewValidationError(errs url.Values) *ValidationError {
	return &ValidationError{errBag: errs}
}

// Add validation errors to errors bag
func (e *ValidationError) Add(key string, vals []string) {
	if e.errBag == nil {
		e.errBag = url.Values{}
	}
	e.errBag[key] = vals
}

// Set validation errors bag
func (e *ValidationError) Set(errBag url.Values) {
	if e.errBag == nil {
		e.errBag = url.Values{}
	}
	e.errBag = errBag
}

// GetAll validation errors
func (e *ValidationError) GetAll() url.Values {
	return e.errBag
}

// Implementation of the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validationError: %v", e.errBag)
}

// Implementation of the error interface
func (e HTTPError) Error() string {
	return fmt.Sprintf("code: %d; message: %v", e.Code, e.Message)
}

// ServeHTTP func is http.Handler interface implementation
func (h Wrap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		switch err.(type) {
		case *ValidationError:
			er := err.(*ValidationError)
			jsonResponse(w, http.StatusUnprocessableEntity, data{"validationError": er.GetAll()})
			return
		case HTTPError:
			er := err.(HTTPError)
			jsonResponse(w, er.Code, data{"error": er.Message})
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			jsonResponse(w, http.StatusNotFound, data{"error": http.StatusText(http.StatusNotFound)})
			return
		}
		log.Println(errors.Wrap(err, "undefined http error"))
		jsonResponse(w, http.StatusInternalServerError, data{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
}

// New is a factory function,
// returns a new instance of the HTTP handler interactor
func New(db *sql.DB, q *repository.Queries, s *storage.Interactor, rf resizerFunc) *Handler {
	return &Handler{
		db:      db,
		query:   q,
		storage: s,
		resize:  rf,
	}
}
