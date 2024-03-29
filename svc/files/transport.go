package files

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/thedevsaddam/govalidator"

	"github.com/SatorNetwork/sator-api/lib/httpencoder"
	"github.com/SatorNetwork/sator-api/lib/utils"
)

// Predefined request query keys
const (
	pageParam         = "page"
	itemsPerPageParam = "items_per_page"
)

type (
	logger interface {
		Log(keyvals ...interface{}) error
	}
)

// MakeHTTPHandler ...
func MakeHTTPHandler(e Endpoints, log logger) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(httpencoder.EncodeError(log, codeAndMessageFrom)),
		httptransport.ServerBefore(jwtkit.HTTPToContext()),
	}

	r.Route("/images", func(r chi.Router) {
		r.Post("/", httptransport.NewServer(
			e.AddImage,
			decodeAddImageRequest,
			httpencoder.EncodeResponse,
			options...,
		).ServeHTTP)

		r.Get("/", httptransport.NewServer(
			e.GetImagesList,
			decodeGetImagesListRequest,
			httpencoder.EncodeResponse,
			options...,
		).ServeHTTP)

		r.Get("/{id}", httptransport.NewServer(
			e.GetImageByID,
			decodeGetImageByIDRequest,
			httpencoder.EncodeResponse,
			options...,
		).ServeHTTP)

		r.Delete("/{id}", httptransport.NewServer(
			e.DeleteImageByID,
			decodeDeleteImageByIDRequest,
			httpencoder.EncodeResponse,
			options...,
		).ServeHTTP)
	})

	r.Post("/", httptransport.NewServer(
		e.AddFile,
		decodeAddFileRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetImagesListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return utils.PaginationRequest{
		Page:         utils.StrToInt32(r.URL.Query().Get(pageParam)),
		ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(itemsPerPageParam)),
	}, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrInvalidParameter) {
		return http.StatusBadRequest, err.Error()
	}

	return httpencoder.CodeAndMessageFrom(err)
}

func decodeGetImageByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed image id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeAddImageRequest(_ context.Context, r *http.Request) (interface{}, error) {
	rules := govalidator.MapData{
		"file:image": []string{
			"required",
			"size:2097152",
			"mime:image/png,image/jpeg",
		},
	}
	if err := utils.Validate(r, rules, nil); err != nil {
		return nil, err
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		return nil, fmt.Errorf("could not parse image from request: %w", err)
	}
	defer file.Close()

	req := AddImageResizeRequest{
		File:       file,
		FileHeader: header,
	}

	if height := r.FormValue("height"); height != "" {
		req.MaxHeight = utils.StrToUint(height)
	}

	if width := r.FormValue("width"); width != "" {
		req.MaxWidth = utils.StrToUint(width)
	}

	return req, nil
}

func decodeAddFileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	rules := govalidator.MapData{
		"file:file": []string{"required"},
	}
	if err := utils.Validate(r, rules, nil); err != nil {
		return nil, err
	}

	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		return nil, fmt.Errorf("could not parse multipart form: %w", err)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("could not parse image from request: %w", err)
	}
	defer file.Close()

	readFile, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	req := AddFileRequest{
		File:       readFile,
		FileHeader: header,
	}

	return req, nil
}

func decodeDeleteImageByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed image id", ErrInvalidParameter)
	}

	return id, nil
}
