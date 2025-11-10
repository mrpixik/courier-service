package repository

import "errors"

var (
	ErrCourierNotFound = errors.New("courier not found")
	ErrUnknownError    = errors.New("internal error")
)
