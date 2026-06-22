package domain

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInvalidInput  = errors.New("invalid input data")
	ErrUnauthorized  = errors.New("unauthorized access")
	ErrForbidden     = errors.New("access denied")
	ErrRateLimited   = errors.New("rate limit exceeded")
)
