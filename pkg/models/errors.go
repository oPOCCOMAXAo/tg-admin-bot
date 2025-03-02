package models

import "errors"

var (
	ErrAuthInvalid    = errors.New("auth invalid")
	ErrFailed         = errors.New("failed")
	ErrNotFound       = errors.New("not found")
	ErrNothingChanged = errors.New("nothing changed")
)
