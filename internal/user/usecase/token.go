package usecase

import "github.com/google/uuid"

type TokenManager interface {
	Generate(*Claims) (token string, err error)
	Parse(token string) (*Claims, error)
}

type Claims struct {
	UserID uuid.UUID
}
