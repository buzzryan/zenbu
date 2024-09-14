package domain

import "github.com/google/uuid"

type Password string

type User struct {
	ID uuid.UUID

	Username string
	Password Password
}
