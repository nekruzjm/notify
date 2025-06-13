package firebase

import (
	"slices"

	"firebase.google.com/go/v4/messaging"
)

const (
	_errRegistrationTokenNotRegistered = "registration-token-not-registered"
	_errMismatchedCredential           = "mismatched-credential"
	_errInvalidToken                   = "registration token is not a valid fmc"
	_errUnregistered                   = "unregistered"
	_errSenderIDMismatched             = "senderid mismatch"
	_errEntityNotFound                 = "requested entity was not found"
	_errTooManyTopics                  = "too-many-topics"
	_errTokenNotFound                  = "Requested entity was not found."
)

func IsValidationErr(err error) bool {
	return slices.Contains([]string{
		_errRegistrationTokenNotRegistered,
		_errMismatchedCredential,
		_errInvalidToken,
		_errUnregistered,
		_errSenderIDMismatched,
		_errEntityNotFound,
		_errTooManyTopics,
		_errTokenNotFound,
	}, err.Error())
}

const (
	_apnsPriorityHeader    = "apns-priority"
	ApnsHighestPriority    = "10"
	ApnsNormalPriority     = "5"
	AndroidHighestPriority = "high"
	AndroidNormalPriority  = "normal"

	_defaultCategory = "c"
	_defaultSound    = "s"
	_defaultThreadID = "t"

	_titleKey   = "title"
	_messageKey = "message"
)

func AndroidMSG(msg *messaging.Message, data map[string]string, priority string) {
	msg.Android = &messaging.AndroidConfig{
		Priority: priority,
		Data:     data,
	}
}

func IosMSG(msg *messaging.Message, data map[string]string, priority string) {
	msg.APNS = &messaging.APNSConfig{
		Headers: map[string]string{
			_apnsPriorityHeader: priority,
		},
		Payload: &messaging.APNSPayload{
			Aps: &messaging.Aps{
				Alert: &messaging.ApsAlert{
					Title: data[_titleKey],
					Body:  data[_messageKey],
				},
				Category:         _defaultCategory,
				Sound:            _defaultSound,
				ThreadID:         _defaultThreadID,
				ContentAvailable: true,
				MutableContent:   true,
			},
			CustomData: mapConvert(data),
		},
	}
}

func mapConvert(data map[string]string) map[string]any {
	var m = make(map[string]any, len(data))
	for k, v := range data {
		m[k] = v
	}
	return m
}
