package push

import (
	"notifications/internal/lib/language"
)

const (
	_requestID = "requestID"
	_apiClient = "apiClient"
)

type externalRequest struct {
	Phone             string            `json:"phone" example:"+992111111111" validate:"required"`
	PersonExternalRef string            `json:"personExternalRef" example:"123456"`
	PushType          string            `json:"type" example:"otp, push" validate:"required"`
	Title             language.Language `json:"title" validate:"required"`
	Body              language.Language `json:"body" validate:"required"`
	ShowInFeed        bool              `json:"showInFeed" validate:"required"`
}
