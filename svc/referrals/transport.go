package referrals

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SatorNetwork/sator-api/internal/httpencoder"
	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
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

	r.Get("/my", httptransport.NewServer(
		e.GetMyReferralCode,
		decodeGetMyReferralCodeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/confirm/{code}", httptransport.NewServer(
		e.StoreUserWithValidCode,
		decodeStoreUserWithValidCodeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetMyReferralCodeRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeStoreUserWithValidCodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	code := chi.URLParam(r, "code")
	if code == "" {
		return nil, fmt.Errorf("%w: missed referral code", ErrInvalidParameter)
	}

	return code, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}
