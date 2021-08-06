package jwt

import (
	"context"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt"
)

// ClaimsFactory is a factory for jwt.Claims.
// Useful in NewParser middleware.
type claimsFactory func() jwt.Claims

// NewParser creates a new JWT token parsing middleware, specifying a
// jwt.Keyfunc interface, the signing method and the claims type to be used. NewParser
// adds the resulting claims to endpoint context or returns error on invalid token.
// Particularly useful for servers.
func newParser(keyFunc jwt.Keyfunc, method jwt.SigningMethod, newClaims claimsFactory) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// tokenString is stored in the context from the transport handlers.
			tokenString, ok := ctx.Value(kitjwt.JWTTokenContextKey).(string)
			if !ok {
				return nil, kitjwt.ErrTokenContextMissing
			}

			// Parse takes the token string and a function for looking up the
			// key. The latter is especially useful if you use multiple keys
			// for your application.  The standard is to use 'kid' in the head
			// of the token to identify which key to use, but the parsed token
			// (head and claims) is provided to the callback, providing
			// flexibility.
			token, err := jwt.ParseWithClaims(tokenString, newClaims(), func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if token.Method != method {
					return nil, kitjwt.ErrUnexpectedSigningMethod
				}

				return keyFunc(token)
			})
			if err != nil {
				if e, ok := err.(*jwt.ValidationError); ok {
					switch {
					case e.Errors&jwt.ValidationErrorMalformed != 0:
						// Token is malformed
						return nil, kitjwt.ErrTokenMalformed
					case e.Errors&jwt.ValidationErrorExpired != 0:
						// Token is expired
						return nil, kitjwt.ErrTokenExpired
					case e.Errors&jwt.ValidationErrorNotValidYet != 0:
						// Token is not active yet
						return nil, kitjwt.ErrTokenNotActive
					case e.Inner != nil:
						// report e.Inner
						return nil, e.Inner
					}
					// We have a ValidationError but have no specific Go kit error for it.
					// Fall through to return original error.
				}
				return nil, err
			}

			if !token.Valid {
				return nil, kitjwt.ErrTokenInvalid
			}

			ctx = context.WithValue(ctx, kitjwt.JWTClaimsContextKey, token.Claims)

			return next(ctx, request)
		}
	}
}

// NewParser returns go-kit parser middleware
func NewParser(signingKey string) endpoint.Middleware {
	return newParser(func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	}, defaultSigningMethod, ClaimsFactory)
}
