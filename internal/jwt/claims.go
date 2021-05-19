package jwt

import (
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/google/uuid"
)

// Claims struct
type Claims struct {
	UserID   string `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	jwt.StandardClaims
}

// ClaimsFactory is a ClaimsFactory that returns
// an empty Claims.
func ClaimsFactory() jwt.Claims {
	return &Claims{}
}

// UserUUID returns user uuid,
// parsed from string Claims.UserID
func (c *Claims) UserUUID() (uuid.UUID, error) {
	if c.UserID == "" {
		return uuid.Nil, ErrUserIDEmpty
	}
	id, err := uuid.Parse(c.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not parse user uuid: %w", err)
	}
	return id, nil
}

// UserIDFromContext returns user uuid from request context
func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	claims := ctx.Value(kitjwt.JWTClaimsContextKey)
	if cl, ok := claims.(*Claims); ok {
		return cl.UserUUID()
	}
	return uuid.Nil, ErrInvalidJWTClaims
}

// UsernameFromContext returns user uuid from request context
func UsernameFromContext(ctx context.Context) (string, error) {
	claims := ctx.Value(kitjwt.JWTClaimsContextKey)
	if cl, ok := claims.(*Claims); ok {
		return cl.Username, nil
	}
	return "", ErrInvalidJWTClaims
}

// TokenIDFromContext returns jwt id from request context
func TokenIDFromContext(ctx context.Context) (uuid.UUID, error) {
	claims := ctx.Value(kitjwt.JWTClaimsContextKey)
	if cl, ok := claims.(*Claims); ok {
		if cl.Id == "" {
			return uuid.Nil, ErrJWTIDEmpty
		}
		id, err := uuid.Parse(cl.Id)
		if err != nil {
			return uuid.Nil, fmt.Errorf("could not parse jwt uuid: %w", err)
		}
		return id, nil
	}
	return uuid.Nil, ErrInvalidJWTClaims
}
