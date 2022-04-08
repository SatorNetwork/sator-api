package httpencoder

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	jwti "github.com/SatorNetwork/sator-api/lib/jwt"
	"github.com/SatorNetwork/sator-api/lib/rbac"
	"github.com/SatorNetwork/sator-api/lib/sumsub"
	"github.com/SatorNetwork/sator-api/lib/validator"
	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
)

type (
	logger interface {
		Log(keyvals ...interface{}) error
	}
)

// EncodeError ...
func EncodeError(l logger, codeAndMessageFrom func(err error) (int, interface{})) httptransport.ErrorEncoder {
	return func(_ context.Context, err error, w http.ResponseWriter) {
		if err == nil {
			l.Log("msg", "encodeError with nil error")
			return
		}

		code, msg := codeAndMessageFrom(err)

		if code == http.StatusInternalServerError {
			// Log only unexpected errors
			l.Log("msg", fmt.Errorf("http transport error: %w", err))
		}

		w.Header().Set(ContentTypeHeader, ContentType)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": msg,
		})
	}
}

// CodeAndMessageFrom helper
func CodeAndMessageFrom(err error) (int, interface{}) {
	if verr, ok := err.(validator.ValidationError); ok {
		return http.StatusUnprocessableEntity, verr.ErrorsBag
	}

	if errors.Is(err, jwt.ErrTokenContextMissing) {
		return http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)
	}

	if errors.Is(err, jwt.ErrTokenExpired) ||
		errors.Is(err, jwt.ErrTokenInvalid) ||
		errors.Is(err, jwt.ErrTokenMalformed) ||
		errors.Is(err, jwt.ErrTokenNotActive) ||
		errors.Is(err, jwt.ErrUnexpectedSigningMethod) {
		return http.StatusUnauthorized, err.Error()
	}

	if errors.Is(err, jwti.ErrInvalidJWTSubject) || errors.Is(err, sumsub.ErrKYCInProgress) ||
		errors.Is(err, sumsub.ErrKYCUserIsDisabled) || errors.Is(err, rbac.ErrAccessDenied) {
		return http.StatusForbidden, err.Error()
	}

	if errors.Is(err, sumsub.ErrKYCNeeded) || errors.Is(err, sumsub.ErrKYCRequiredDocs) {
		return http.StatusProxyAuthRequired, err.Error()
	}

	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound, err.Error()
	}

	switch err {
	case jwt.ErrTokenContextMissing:
		return http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)
	case jwt.ErrTokenExpired,
		jwt.ErrTokenInvalid,
		jwt.ErrTokenMalformed,
		jwt.ErrTokenNotActive,
		jwt.ErrUnexpectedSigningMethod:
		return http.StatusUnauthorized, err.Error()
	default:
		return http.StatusInternalServerError, err.Error()
	}
}
