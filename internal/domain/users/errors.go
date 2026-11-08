package users

import "errors"

var (
	ErrUsernameUnavailable     = errors.New("username unavailable")
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidLoginCredentials = errors.New("invalid login credentials")
	ErrUserNotActivated        = errors.New("user not activated")
)
