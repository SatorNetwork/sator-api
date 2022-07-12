package settings

import (
	"encoding/json"
	"github.com/SatorNetwork/sator-api/lib/httpencoder"
	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
	"net/http"
)

type (
	logger interface {
		Log(keyvals ...interface{}) error
	}
)

func MakeHTTPHandler(e Endpoints, log logger, encodeResponse httptransport.EncodeResponseFunc) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(httpencoder.EncodeError(log, codeAndMessageFrom)),
		httptransport.ServerBefore(jwtkit.HTTPToContext()),
	}

	r.Get("/", httptransport.NewServer(
		e.GetSettings,
		decodeGetSettingsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/value-types", httptransport.NewServer(
		e.GetSettingsValueTypes,
		decodeGetSettingsValueTypesRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{key}", httptransport.NewServer(
		e.GetSettingsByKey,
		decodeGetSettingsByKeyRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/", httptransport.NewServer(
		e.AddSetting,
		decodeAddSettingRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/{key}", httptransport.NewServer(
		e.UpdateSetting,
		decodeUpdateSettingRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/{key}", httptransport.NewServer(
		e.DeleteSetting,
		decodeDeleteSettingRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetSettingsRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetSettingsValueTypesRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeAddSettingRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req AddGameSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUpdateSettingRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req UpdateGameSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.Key = chi.URLParam(r, "key")

	return req, nil
}

func decodeDeleteSettingRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return chi.URLParam(r, "key"), nil
}

func decodeGetSettingsByKeyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return chi.URLParam(r, "key"), nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}
