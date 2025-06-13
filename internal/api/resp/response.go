package resp

import (
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"

	"notifications/internal/api/resp/code"
)

const (
	_contentType     = "Content-Type"
	_applicationJSON = "application/json"
)

func GinJSONAbort(c *gin.Context, statusCode int, resp any) {
	JSON(c.Writer, statusCode, resp)
	c.Abort()
}

func JSON(w http.ResponseWriter, statusCode int, resp any) {
	body, err := sonic.Marshal(resp)
	if err != nil {
		w.WriteHeader(code.InternalErr)
		return
	}

	w.Header().Set(_contentType, _applicationJSON)
	w.WriteHeader(statusCode)
	_, _ = w.Write(body)
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Payload any    `json:"payload"`
}

func newResponse(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

var (
	Success         = newResponse(code.Success, "Success")
	Accepted        = newResponse(code.Accepted, "Accepted")
	BadRequest      = newResponse(code.BadRequest, "Bad request")
	NotFound        = newResponse(code.NotFound, "Not found")
	TokenExpired    = newResponse(code.NotFound, "App auth credits expired")
	InternalErr     = newResponse(code.InternalErr, "Internal server error")
	Unauthorized    = newResponse(code.Unauthorized, "Unauthorized")
	TooManyRequests = newResponse(code.TooManyRequests, "Too many requests")
	Forbidden       = newResponse(code.Forbidden, "Forbidden")
	DuplicateRecord = newResponse(code.BadRequest, "Record is duplicated")

	SamePassword         = newResponse(code.SamePassword, "Same password")
	RequiredFieldsAbsent = newResponse(code.RequiredFields, "Required fields absent")
	UserBlocked          = newResponse(code.UserBlocked, "User blocked")
	WrongPassword        = newResponse(code.WrongPassword, "Wrong password")
	UserNotFound         = newResponse(code.NotFound, "User not found")
	FileSizeExceeded     = newResponse(code.FileSizeExceeded, "File size is too large to upload")
)
