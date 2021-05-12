package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

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

// NewWithUserID returns signed JWT string with user id in claims
func (i *JWT) NewWithUserID(userID uuid.UUID) (uuid.UUID, string, error) {
	tokenID := uuid.New()
	claims := &Claims{
		userID.String(),
		jwt.StandardClaims{
			Id:        tokenID.String(),
			ExpiresAt: time.Now().Add(i.expIn).Unix(),
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(i.signingKey)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("could not sign token: %w", err)
	}
	return tokenID, ss, nil
}
