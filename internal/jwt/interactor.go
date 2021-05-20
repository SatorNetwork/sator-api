package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var defaultSigningMethod jwt.SigningMethod = jwt.SigningMethodHS256

type (
	// JWT struct
	JWT struct {
		signingKey []byte
		expIn      time.Duration
	}
)

// NewInteractor is a factory function,
// returns a new instance of the JWT struct
func NewInteractor(signingKey string, expiresIn time.Duration) *JWT {
	return &JWT{
		signingKey: []byte(signingKey),
		expIn:      expiresIn,
	}
}

// NewWithUserData returns signed JWT string with user id and username in claims
func (i *JWT) NewWithUserData(userID uuid.UUID, username string) (uuid.UUID, string, error) {
	tokenID := uuid.New()
	claims := &Claims{
		userID.String(),
		username,
		jwt.StandardClaims{
			Id:        tokenID.String(),
			ExpiresAt: time.Now().Add(i.expIn).Unix(),
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(defaultSigningMethod, claims)
	ss, err := token.SignedString(i.signingKey)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("could not sign token: %w", err)
	}
	return tokenID, ss, nil
}
