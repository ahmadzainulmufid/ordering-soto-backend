package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret     []byte
	issuer     string
	expiration time.Duration
}

type UserClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTManager(
	secret string,
	issuer string,
	expiration time.Duration,
) *JWTManager {
	return &JWTManager{
		secret:     []byte(secret),
		issuer:     issuer,
		expiration: expiration,
	}
}

func (m *JWTManager) GenerateAccessToken(
	userID int64,
	role string,
) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(m.expiration)

	claims := UserClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", 0, fmt.Errorf("gagal membuat access token: %w", err)
	}

	return signedToken, int64(m.expiration.Seconds()), nil
}

func (m *JWTManager) ValidateAccessToken(
	tokenString string,
) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("algoritma token tidak valid")
			}

			return m.secret, nil
		},
		jwt.WithIssuer(m.issuer),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, fmt.Errorf("token tidak valid: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("claims token tidak valid")
	}

	return claims, nil
}