package service

import "errors"

var (
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidPhone    = errors.New("invalid phone")
	ErrInvalidStatus   = errors.New("invalid status")
	ErrCourierExists   = errors.New("courier already exists")
	ErrCourierNotFound = errors.New("courier not found")
	ErrInternalError   = errors.New("internal error")
)
