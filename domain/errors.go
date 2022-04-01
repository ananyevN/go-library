package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrNotFound            = errors.New("book was not found")
	ErrConflict            = errors.New("book already exist")
)
