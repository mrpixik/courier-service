package repository

import "errors"

var (
	ErrCourierExists   = errors.New("courier already exists")
	ErrCourierNotFound = errors.New("courier not found")
	ErrInternalError   = errors.New("internal error")
)
