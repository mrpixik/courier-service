package repository

import "errors"

var (
	ErrCourierNotFound      = errors.New("courier not found")
	ErrCourierAlreadyExists = errors.New("courier already exists")
	ErrUnknownError         = errors.New("internal error")
)
