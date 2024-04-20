package storage

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserExists        = errors.New("user exists")
	ErrAppNotFound       = errors.New("app not found")
	ErrExceptionNotFound = errors.New("exception not found")
	ErrExpressionExists  = errors.New("expression exists")
	ErrInvalidSign       = errors.New("invalid sign")
	ErrSignExists        = errors.New("sign exists")
)
