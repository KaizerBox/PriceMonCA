package domain

import "errors"

var (
    ErrNotFound     = errors.New("not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrUnavailable  = errors.New("unavailable")
    ErrConflict     = errors.New("conflict")
)
