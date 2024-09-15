package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// TokenManager is an interface for generating and parsing tokens.
type TokenManager interface {
	Generate(*Claims) (token string, err error)
	Parse(token string) (*Claims, error)
}

type Claims struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
}

var (
	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)
