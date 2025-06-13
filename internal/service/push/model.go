package push

import (
	"errors"
	"slices"

	"notifications/internal/lib/language"
	"notifications/pkg/util/strset"
)

const _slashDelim = "/"

// default fake message in case of user is inactive
const (
	_inactiveUserMessageID = "inactive_user#fake_message_id"
	_disabledPushMessageID = "disabled_push#fake_message_id"
	_fcmPushMessageID      = "firebase_error#fake_message_id"
	_defaultAPIClient      = "my.app"
)

// Push types
const (
	_pushType = "pushType"
	_silent   = "silent"
	_otp      = "otp"
	_push     = "push"
)

const (
	_title       = "title"
	_comment     = "comment"
	_message     = "message"
	_badge       = "badge"
	_sectionName = "sectionName"
	_category    = "category"
	_trID        = "transactionID"
)

const (
	_badge0             = "0"
	_badge1             = "1"
	_defaultSectionName = "pushHistory"
	_defaultCategory    = "OpenUISection"
)

const (
	_approved = "approved"
	_failed   = "failed"
)

const _active = "active"

type Message struct {
	UserID     int
	Token      string
	Data       map[string]string
	ShowInFeed bool
}

type Request struct {
	InternalRequest InternalRequest
	ExternalRequest ExternalRequest
	ShowInFeed      bool
	IsInternal      bool
	Sync            bool
}

type InternalRequest struct {
	UserID int
	Token  string
	Data   map[string]string
}

type ExternalRequest struct {
	ID                string
	APIClient         string
	Phone             string
	PersonExternalRef string
	PushType          string
	Title             language.Language
	Body              language.Language
}

func (r *Request) validate() error {
	if !slices.Contains([]string{_otp, _push}, r.ExternalRequest.PushType) {
		return errors.New("invalid push type")
	}

	if strset.IsEmpty(r.ExternalRequest.PersonExternalRef) && strset.IsEmpty(r.ExternalRequest.Phone) {
		return errors.New("user phone and person external ref cannot be empty")
	}

	if !strset.IsEmpty(r.ExternalRequest.Phone) {
		if strset.IsEmpty(strset.GetDigits(r.ExternalRequest.Phone)) {
			return errors.New("invalid phone number")
		}

		const _plusSymbol = '+'
		if r.ExternalRequest.Phone[0] != _plusSymbol {
			return errors.New("phone number must start with +")
		}
	}

	return nil
}
