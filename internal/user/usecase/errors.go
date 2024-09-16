package usecase

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("user with this username already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidPassword       = errors.New("invalid password")
)
