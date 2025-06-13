package code

import "net/http"

const (
	Success         = http.StatusOK
	Accepted        = http.StatusAccepted
	BadRequest      = http.StatusBadRequest
	Unauthorized    = http.StatusUnauthorized
	TooManyRequests = http.StatusTooManyRequests
	InternalErr     = http.StatusInternalServerError
	NotFound        = http.StatusNotFound
	Forbidden       = http.StatusForbidden

	RequiredFields = iota + 1500
	UserBlocked
	WrongPassword
	SamePassword
	FileSizeExceeded
)
