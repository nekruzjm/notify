package email

import (
	"context"
	"strings"

	"go.uber.org/zap"

	"notifications/pkg/lib/notifier/email"
)

func (s *service) Send(_ context.Context, request Email) error {
	var (
		host      = s.config.GetString("email.host")
		port      = s.config.GetString("email.port")
		from      = s.config.GetString("email.from")
		receivers = []string{request.Body[_userEmail]}
		text      = request.Body[_text]
		subject   = request.Body[_subject]
		body      strings.Builder
	)

	body.WriteString("Subject: ")
	body.WriteString(subject)
	body.WriteString("\n")
	body.WriteString(_defaultMimeTextPlain)
	body.WriteString(_defaultTmplHeader)
	body.WriteString(text)
	body.WriteString(_defaultTmplFooter)

	err := email.New().
		SetHost(host).
		SetPort(port).
		SetPlainAuth(s.plainAuth).
		SetSender(from).
		SetReceiver(receivers).
		SetBody([]byte(body.String())).
		Send()
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to send email", zap.Error(err))
		return err
	}

	return nil
}
