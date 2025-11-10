package service

import "errors"

var (
	ErrInvalidName          = errors.New("invalid name")
	ErrInvalidPhone         = errors.New("invalid phone")
	ErrInvalidStatus        = errors.New("invalid status")
	ErrInternalError        = errors.New("internal error")
	ErrCourierAlreadyExists = errors.New("courier already exists")
	ErrCourierNotFound      = errors.New("courier not found")
)
