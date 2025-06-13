package repomodel

import "errors"

var (
	ErrUniqueViolation = errors.New("unique violation")
	ErrNotFound        = errors.New("not found")
)
