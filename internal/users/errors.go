package users

import "errors"

var (
	ErrInvalidStatus     = errors.New("invalid status provided")
	ErrInvalidRole       = errors.New("invalid role provided")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidParams     = errors.New("invalid params provided")
	ErrInvalidFullname   = errors.New("invalid fullname provided")
	ErrUserAlreadyExists = errors.New("user already exists")
)
