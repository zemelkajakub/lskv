package profile

import "errors"

var (
	// ErrProfileNotFound is returned when a profile does not exist
	ErrProfileNotFound = errors.New("profile not found")
)
