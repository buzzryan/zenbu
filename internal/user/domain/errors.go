package domain

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("user with this username already exists")
)
