package sentry

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

const _dev = "dev"

func (s *sentryObserver) CurrentHub() *sentry.Hub {
	return s.hub
}

func (s *sentryObserver) CaptureException(err error) {
	if s.stage == _dev {
		s.logger.Debug("dev stage. Skip sending err to Sentry", zap.Error(err))
		return
	}
	s.hub.CaptureException(err)
}
