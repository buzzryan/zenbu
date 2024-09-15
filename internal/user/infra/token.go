package infra

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/buzzryan/zenbu/internal/user/usecase"
)

type jwsTokenManager struct {
	signingKey string
}

type jwsClaims struct {
	jwt.RegisteredClaims
	UserID string
}

func NewJWSTokenManager(signingKey string) usecase.TokenManager {
	return &jwsTokenManager{signingKey: signingKey}
}

func (j *jwsTokenManager) Parse(token string) (*usecase.Claims, error) {
	t, err := jwt.ParseWithClaims(token, &jwsClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.signingKey), nil
	})
	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, usecase.ErrTokenExpired
	}
	if err != nil {
		return nil, errors.Join(usecase.ErrInvalidToken, err)
	}

	claims, ok := t.Claims.(*jwsClaims)
	if !ok {
		return nil, usecase.ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user id is not a valid UUID: %w", err)
	}
	return &usecase.Claims{
		UserID:    userID,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

func (j *jwsTokenManager) Generate(claims *usecase.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwsClaims{
		UserID: claims.UserID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(claims.ExpiresAt),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		}})
	signed, err := token.SignedString([]byte(j.signingKey))
	if err != nil {
		return "", fmt.Errorf("unexpected error when generating JWS token: %w", err)
	}
	return signed, nil
}
