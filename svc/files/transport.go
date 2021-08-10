package files

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/SatorNetwork/sator-api/internal/utils"
	"github.com/thedevsaddam/govalidator"

	"github.com/SatorNetwork/sator-api/internal/httpencoder"

	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
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

	return r
}

func decodeGetImagesListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return PaginationRequest{
		Page:          utils.StrToInt32(r.URL.Query().Get(pageParam)),
		ImagesPerPage: utils.StrToInt32(r.URL.Query().Get(itemsPerPageParam)),
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
		"file:image": []string{"required", "ext:png", "size:2097152", "mime:image/png"},
	}
	if err := utils.Validate(r, rules, nil); err != nil {
		return nil, err
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		return nil, fmt.Errorf("could not parse image from request: %w", err)
	}
	defer file.Close()

	var req AddImageRequest
	req.File = file
	req.FileHeader = header

	return req, nil
}

func decodeDeleteImageByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed image id", ErrInvalidParameter)
	}

	return id, nil
}
