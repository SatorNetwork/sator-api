package announcement

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/lib/httpencoder"
)

type (
	logger interface {
		Log(keyvals ...interface{}) error
	}
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
)

func MakeHTTPHandler(e Endpoints, log logger) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(httpencoder.EncodeError(log, codeAndMessageFrom)),
		httptransport.ServerBefore(jwtkit.HTTPToContext()),
	}

	r.Post("/", httptransport.NewServer(
		e.CreateAnnouncement,
		decodeCreateAnnouncementRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{id}", httptransport.NewServer(
		e.GetAnnouncementByID,
		decodeGetAnnouncementByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/{id}", httptransport.NewServer(
		e.UpdateAnnouncement,
		decodeUpdateAnnouncementRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/{id}", httptransport.NewServer(
		e.DeleteAnnouncement,
		decodeDeleteAnnouncementRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/", httptransport.NewServer(
		e.ListAnnouncements,
		decodeListAnnouncementsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/unread", httptransport.NewServer(
		e.ListUnreadAnnouncements,
		decodeListUnreadAnnouncementsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{id}/read", httptransport.NewServer(
		e.MarkAsRead,
		decodeMarkAsReadRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/read_all", httptransport.NewServer(
		e.MarkAllAsRead,
		decodeMarkAllAsReadRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeCreateAnnouncementRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req CreateAnnouncementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "could not decode request body")
	}

	return &req, nil
}

func decodeGetAnnouncementByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed announcement id parameter", ErrInvalidParameter)
	}
	return &GetAnnouncementByIDRequest{
		ID: id,
	}, nil
}

func decodeUpdateAnnouncementRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req UpdateAnnouncementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "could not decode request body")
	}

	return &req, nil
}

func decodeDeleteAnnouncementRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req DeleteAnnouncementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "could not decode request body")
	}

	return &req, nil
}

func decodeListAnnouncementsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return &Empty{}, nil
}

func decodeListUnreadAnnouncementsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return &Empty{}, nil
}

func decodeMarkAsReadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	announcementsID := chi.URLParam(r, "id")
	if announcementsID == "" {
		return nil, fmt.Errorf("%w: missed announcement id parameter", ErrInvalidParameter)
	}
	return &MarkAsReadRequest{
		AnnouncementID: announcementsID,
	}, nil
}

func decodeMarkAllAsReadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return &Empty{}, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}
