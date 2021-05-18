package jwt

import (
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

// NewParser returns go-kit parser middleware
func NewParser(signingKey string) endpoint.Middleware {
	return kitjwt.NewParser(func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	}, defaultSigningMethod, ClaimsFactory)
}
