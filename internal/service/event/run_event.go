package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"notifications/internal/api/resp"
	"notifications/internal/api/transport/broker/stream"
	"notifications/internal/api/transport/broker/subject"
	"notifications/internal/repo/repomodel"
	"notifications/internal/repo/rom"
	"notifications/internal/service/admin"
)

func (s *service) RunEvent(ctx context.Context, a admin.Admin, id int) (any, error) {
	s.logger.Info("RunEvent start", zap.Int("eventID", id))

	tx := s.transactor.New()
	err := tx.Begin(ctx)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("err occurred during beginning transaction", zap.Error(err))
		return nil, err
	}

	defer func() {
		if err != nil {
			if errX := tx.Rollback(ctx); errX != nil {
				err = errors.Join(err, fmt.Errorf("errX: %w", errX))
			}
			return
		}
		err = tx.Commit(ctx)
		if err != nil {
			err = errors.Join(err, errors.New("tx commit failed"))
		}
	}()

	selectedEvent, err := tx.EventRepo().GetActiveByIDWithLock(ctx, id)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during getting event", zap.Error(err), zap.Int("id", id))
			return nil, err
		}
		return nil, resp.Wrap(resp.ErrNotFound, "event not found or not active")
	}

	if selectedEvent.Status == _sent {
		return nil, resp.Wrap(resp.ErrBadRequest, "cannot load users for event with sent status")
	}

	var oldEvent = *selectedEvent

	selectedEvent.Status = _sent
	err = tx.EventRepo().UpdateStatus(ctx, selectedEvent.ID, selectedEvent.Status)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("err occurred during updating event", zap.Error(err), zap.Int("eventID", id))
		return nil, err
	}

	userIDs, err := tx.UserRepo().GetUserIDsByEventID(ctx, id)
	if err != nil {
		if errors.Is(err, repomodel.ErrNotFound) {
			return nil, resp.Wrap(resp.ErrNotFound, "there are no users subscribed to this event")
		}
		s.sentry.CaptureException(err)
		s.logger.Error("err occurred during getting user relations", zap.Error(err), zap.Int("eventID", id))
		return nil, err
	}

	err = s.romRepo.BatchInsert(ctx, id, userIDs, &rom.Event{
		ID:        selectedEvent.ID,
		Title:     selectedEvent.Title,
		Body:      selectedEvent.Body,
		Image:     selectedEvent.Image,
		ExtraData: selectedEvent.ExtraData,
	})
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("err from BatchInsert", zap.Error(err), zap.Int("eventID", id))
		return nil, resp.Wrap(resp.ErrInternalErr, err.Error())
	}

	var messages = setupMessages(selectedEvent)

	s.logger.Info("firebase messaging request", zap.Any("messages", messages), zap.Int("eventID", id))

	response, err := s.fcmSender.SendEach(ctx, messages)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("err occurred during sending push", zap.Error(err), zap.Int("eventID", id))
		return nil, err
	}

	s.logger.Info("firebase messaging response", zap.Any("response", response), zap.Int("eventID", id))

	var serviceResponse = &struct {
		SuccessCount int   `json:"successCount"`
		FailedCount  int   `json:"failedCount"`
		Result       []any `json:"result"`
	}{
		SuccessCount: response.SuccessCount,
		FailedCount:  response.FailureCount,
		Result:       make([]any, 0, len(response.Responses)),
	}

	for _, res := range response.Responses {
		serviceResponse.Result = append(serviceResponse.Result, struct {
			MessageID string `json:"messageID"`
			Error     error  `json:"error"`
			Success   bool   `json:"success"`
		}{
			MessageID: res.MessageID,
			Error:     res.Error,
			Success:   res.Success,
		})
	}

	s.logger.Info("RunEvent end", zap.Int("eventID", id))

	if a != (admin.Admin{}) {
		err = s.nats.Publish(stream.Audit, subject.AuditAdd, admin.Audit{
			AdminId:   a.ID,
			IpAddress: a.IP,
			EventName: admin.RunEvent,
			OldData:   oldEvent,
			NewData:   *selectedEvent,
			CreatedAt: time.Now(),
		})
		if err != nil {
			s.sentry.CaptureException(err)
			s.logger.Error("failed to publish audit event", zap.Error(err))
		}
	}

	var event = new(Event)
	event.toService(selectedEvent)
	err = s.nats.Publish(stream.Notifications, subject.NotificationsTopicUsersUnsubscribed, event)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("cannot publish message", zap.Error(err), zap.Int("id", id))
		return nil, err
	}

	return serviceResponse, nil
}
