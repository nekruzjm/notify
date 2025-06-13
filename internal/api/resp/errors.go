package resp

import (
	"errors"
	"fmt"

	"notifications/internal/api/resp/code"
)

var (
	ErrUnauthorized        errResponder = &Err{code.Unauthorized, "unauthorized admin user"}
	ErrNotFound            errResponder = &Err{code.NotFound, "not found"}
	ErrUserNotFound        errResponder = &Err{code.NotFound, "user not found"}
	ErrBadRequest          errResponder = &Err{code.BadRequest, "bad request"}
	ErrRequiredFieldAbsent errResponder = &Err{code.RequiredFields, "required fields absent"}
	ErrForbidden           errResponder = &Err{code.Forbidden, "forbidden"}
	ErrUniqueViolation     errResponder = &Err{code.InternalErr, "unique Violation error"}
	ErrAlreadyExpired      errResponder = &Err{code.NotFound, "cache already expired"}
	ErrDuplicateItem       errResponder = &Err{code.BadRequest, "duplicate item"}
	ErrNoRowsAffected      errResponder = &Err{code.BadRequest, "no rows were affected"}
	ErrUserBlocked         errResponder = &Err{code.UserBlocked, "user blocked"}
	ErrWrongPassword       errResponder = &Err{code.WrongPassword, "wrong password"}
	ErrFileSizeExceeded    errResponder = &Err{code.FileSizeExceeded, "file size is too large"}
	ErrInternalErr         errResponder = &Err{code.InternalErr, "internal error"}
)

type errResponder interface {
	Error() string
	ErrCode() int
}

type Err struct {
	code   int
	errMsg string
}

func newErr(message string, code int) *Err {
	return &Err{
		code:   code,
		errMsg: message,
	}
}

func (e *Err) Error() string {
	return e.errMsg
}

func (e *Err) ErrCode() int {
	return e.code
}

func Wrap(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

func RespondErr(err error) (response Response) {
	var unwrapped = errors.Unwrap(err)
	if unwrapped == nil {
		unwrapped = err
	}

	var apiErr *Err
	errors.As(unwrapped, &apiErr)

	if apiErr == nil {
		apiErr = newErr(err.Error(), code.InternalErr)
	}

	switch {
	case errors.Is(apiErr, ErrNotFound):
		response = NotFound
	case errors.Is(apiErr, ErrUnauthorized):
		response = Unauthorized
	case errors.Is(apiErr, ErrUserNotFound):
		response = UserNotFound
	case errors.Is(apiErr, ErrBadRequest), errors.Is(apiErr, ErrNoRowsAffected):
		response = BadRequest
	case errors.Is(apiErr, ErrRequiredFieldAbsent):
		response = RequiredFieldsAbsent
	case errors.Is(apiErr, ErrForbidden):
		response = Forbidden
	case errors.Is(apiErr, ErrUniqueViolation), errors.Is(apiErr, ErrDuplicateItem):
		response = DuplicateRecord
	case errors.Is(apiErr, ErrAlreadyExpired):
		response = TokenExpired
	case errors.Is(apiErr, ErrUserBlocked):
		response = UserBlocked
	case errors.Is(apiErr, ErrWrongPassword):
		response = WrongPassword
	case errors.Is(apiErr, ErrFileSizeExceeded):
		response = FileSizeExceeded
	default:
		response = InternalErr
	}

	response.Code = apiErr.ErrCode()
	response.Message = err.Error()
	return
}
