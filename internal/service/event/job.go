package event

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"notifications/internal/repo/repomodel"
	"notifications/internal/service/admin"
)

func (s *service) RunJob() {
	var ctx = context.Background()

	events, err := s.eventRepo.GetAllActive(ctx)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during getting active events", zap.Error(err))
		}
		return
	}

	var currentTime = time.Now()

	for _, event := range events {
		if currentTime.After(event.ScheduledAt) || currentTime.Equal(event.ScheduledAt) {
			response, err := s.RunEvent(ctx, admin.Admin{}, event.ID)
			if err != nil {
				s.logger.Error("err occurred during running event", zap.Error(err), zap.Int("id", event.ID))

				event.Status = _failed
				event.ExtraData["reason"] = err.Error()
				_, err = s.eventRepo.Update(ctx, event)
				if err != nil {
					s.sentry.CaptureException(err)
					s.logger.Error("err occurred during updating event status", zap.Error(err), zap.Int("id", event.ID))
					continue
				}
			}

			s.logger.Info("Message successfully sent to topic", zap.Any("event", event), zap.Any("response", response))
		}
	}
}
