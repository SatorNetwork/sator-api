package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// Predefined token types
const (
	AccessToken  = "access_token"
	RefreshToken = "refresh_token"
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
func (i *JWT) NewWithUserData(userID uuid.UUID, username, role string) (uuid.UUID, string, error) {
	tokenID := uuid.New()
	claims := &Claims{
		userID.String(),
		username,
		role,
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

// NewWithUserData returns signed JWT string with user id and username in claims
func (i *JWT) NewWithRefreshToken(userID uuid.UUID, username, role string) (access, refresh string, err error) {
	accessToken := jwt.NewWithClaims(defaultSigningMethod, &Claims{
		userID.String(),
		username,
		role,
		jwt.StandardClaims{
			Id:        uuid.New().String(),
			Subject:   AccessToken,
			ExpiresAt: time.Now().Add(i.expIn).Unix(),
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})
	accessTokenStr, err := accessToken.SignedString(i.signingKey)
	if err != nil {
		return "", "", fmt.Errorf("could not sign access token: %w", err)
	}

	refreshToken := jwt.NewWithClaims(defaultSigningMethod, &Claims{
		userID.String(),
		username,
		role,
		jwt.StandardClaims{
			Id:        uuid.New().String(),
			Subject:   RefreshToken,
			ExpiresAt: time.Now().Add(i.expIn + (30 * 24 * time.Hour)).Unix(), // access token exp time + 30 days
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})
	refreshTokenStr, err := refreshToken.SignedString(i.signingKey)
	if err != nil {
		return "", "", fmt.Errorf("could not sign refresh token: %w", err)
	}

	return accessTokenStr, refreshTokenStr, nil
}
