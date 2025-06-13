package sms

import (
	"context"
	"time"

	"go.uber.org/zap"

	"notifications/pkg/lib/notifier/sms"
)

func (s *service) Send(ctx context.Context, message map[string]string) error {
	err := s.sms.Send(ctx, sms.Request{
		Phone:         message[_phone],
		Text:          message[_text],
		SenderAddress: _defaultSender,
		Priority:      _defaultPriority,
		ExpiresIn:     _defaultExpiration,
		SmsType:       _defaultType,
		ScheduledAt:   time.Now().UTC().Add(-(time.Second * 5)),
	})
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("err occurred during send message", zap.Error(err))
		return err
	}

	return nil
}
