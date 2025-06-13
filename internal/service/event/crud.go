package event

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"notifications/internal/api/resp"
	"notifications/internal/api/transport/broker/stream"
	"notifications/internal/api/transport/broker/subject"
	"notifications/internal/lib/language"
	"notifications/internal/repo/event"
	"notifications/internal/repo/repomodel"
	"notifications/internal/service/admin"
	"notifications/pkg/lib/tinypng"
	"notifications/pkg/util/strset"
	"notifications/pkg/util/timeset"
)

func (s *service) GetEvents(ctx context.Context, filter Filter) ([]*Event, error) {
	list, err := s.eventRepo.GetByFilter(ctx, event.Filter{
		ID:     filter.ID,
		Topic:  filter.Topic,
		Status: filter.Status,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	})
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("failed to get events by filter", zap.Error(err))
			return nil, err
		}
		return nil, resp.Wrap(resp.ErrNotFound, "events not found")
	}

	var events = make([]*Event, 0, len(list))
	for _, item := range list {
		var serviceEvent = new(Event)
		serviceEvent.toService(item)
		s.setImgURL(serviceEvent)
		events = append(events, serviceEvent)
	}

	return events, nil
}

func (s *service) Create(ctx context.Context, a admin.Admin, request *Request) (*Event, error) {
	if strset.IsEmpty(request.Topic) {
		return nil, resp.Wrap(resp.ErrBadRequest, "topic cannot be empty")
	}

	err := s.validate(request)
	if err != nil {
		return nil, err
	}

	var eventItem = &event.Event{
		ID:          int(s.idGenerator.Generate().Int64()),
		Status:      request.Status,
		Title:       request.Title,
		Body:        request.Body,
		Topic:       request.Topic,
		Category:    request.Category,
		Link:        request.Link,
		ExtraData:   request.ExtraData,
		ScheduledAt: request.ScheduledAtTime,
	}

	createdEvent, err := s.eventRepo.Create(ctx, eventItem)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to create events", zap.Error(err))
		return nil, err
	}

	var item = new(Event)
	item.toService(createdEvent)
	s.setImgURL(item)

	err = s.nats.Publish(stream.Audit, subject.AuditAdd, admin.Audit{
		AdminId:   a.ID,
		IpAddress: a.IP,
		EventName: admin.CreateNotificationsEvent,
		OldData:   event.Event{},
		NewData:   *eventItem,
		CreatedAt: time.Now(),
	})
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to publish audit event", zap.Error(err))
	}

	return item, nil
}

func (s *service) Update(ctx context.Context, a admin.Admin, request *Request) (*Event, error) {
	err := s.validate(request)
	if err != nil {
		return nil, err
	}

	selectedEvent, err := s.eventRepo.GetByID(ctx, request.ID)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("failed to get events by id", zap.Error(err), zap.Int("id", request.ID))
			return nil, err
		}
		return nil, resp.Wrap(resp.ErrNotFound, "event not found")
	}

	if selectedEvent.Status == _sent {
		return nil, resp.Wrap(resp.ErrBadRequest, "cannot change event with sent status")
	}

	var oldEvent = *selectedEvent

	var _emptyLang language.Language

	if request.Title != _emptyLang {
		selectedEvent.Title = request.Title
	}
	if request.Body != _emptyLang {
		selectedEvent.Body = request.Body
	}
	if !strset.IsEmpty(request.Status) && selectedEvent.Status != request.Status {
		selectedEvent.Status = request.Status
	}
	if !strset.IsEmpty(request.Category) && selectedEvent.Category != request.Category {
		selectedEvent.Category = request.Category
	}
	if !strset.IsEmpty(request.Link) && selectedEvent.Link != request.Link {
		selectedEvent.Link = request.Link
	}
	if !strset.IsEmpty(request.ScheduledAt) && !selectedEvent.ScheduledAt.Equal(request.ScheduledAtTime) {
		selectedEvent.ScheduledAt = request.ScheduledAtTime
	}

	selectedEvent.ExtraData = request.ExtraData

	updatedEvent, err := s.eventRepo.Update(ctx, selectedEvent)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to update events", zap.Error(err), zap.Int("id", request.ID))
		return nil, err
	}

	var item = new(Event)
	item.toService(updatedEvent)
	s.setImgURL(item)

	err = s.nats.Publish(stream.Audit, subject.AuditAdd, admin.Audit{
		AdminId:   a.ID,
		IpAddress: a.IP,
		EventName: admin.UpdateNotificationsEvent,
		OldData:   oldEvent,
		NewData:   *selectedEvent,
		CreatedAt: time.Now(),
	})
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to publish audit event", zap.Error(err))
	}

	return item, nil
}

func (s *service) validate(request *Request) error {
	if !strset.IsEmpty(request.Topic) {
		ok, err := regexp.MatchString(_topicRegex, request.Topic)
		if err != nil {
			s.logger.Warning("failed to match topic regex", zap.Error(err))
			return resp.Wrap(resp.ErrBadRequest, err.Error())
		}
		if !ok {
			s.logger.Warning("topic is not valid", zap.String("topic", request.Topic))
			return resp.Wrap(resp.ErrBadRequest, "topic is not valid")
		}
	}

	if strset.IsEmpty(request.Status) {
		return resp.Wrap(resp.ErrBadRequest, "status cannot be empty")
	}

	if !slices.Contains([]string{_active, _draft}, request.Status) {
		s.logger.Warning("status is not valid", zap.String("status", request.Status))
		return resp.Wrap(resp.ErrBadRequest, "status is not valid")
	}

	scheduledAt, err := time.Parse(timeset.Layout, request.ScheduledAt)
	if err != nil {
		s.logger.Warning("scheduled at is not valid", zap.String("status", request.ScheduledAt))
		return resp.Wrap(resp.ErrBadRequest, "scheduled at is not valid")
	}

	scheduledAt = scheduledAt.UTC()
	if scheduledAt.Before(time.Now()) || scheduledAt.Equal(time.Now()) {
		s.logger.Warning("scheduled at cannot be before now", zap.String("status", request.ScheduledAt))
		return resp.Wrap(resp.ErrBadRequest, "scheduled at cannot be before or equal to current time")
	}

	request.ScheduledAtTime = scheduledAt
	return nil
}

func (s *service) Delete(ctx context.Context, a admin.Admin, id int) error {
	tx := s.transactor.New()
	err := tx.Begin(ctx)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("err on tx.Begin", zap.Error(err))
		return err
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

	selectedEvent, err := tx.EventRepo().GetByID(ctx, id)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("failed to get events by id", zap.Error(err), zap.Int("eventID", id))
			return err
		}
		return resp.Wrap(resp.ErrNotFound, "event not found")
	}

	for _, img := range selectedEvent.Image.GetAll() {
		if !strset.IsEmpty(img) {
			var eg errgroup.Group
			for _, size := range tinypng.Sizes() {
				eg.Go(func() error {
					if err = s.fileManager.Remove(&s.bucket, s.directory+size.Format, img); err != nil {
						s.sentry.CaptureException(err)
						s.logger.Error("failed to remove image", zap.Error(err))
						return err
					}
					return nil
				})
			}
			if err = eg.Wait(); err != nil {
				return resp.Wrap(resp.ErrInternalErr, err.Error())
			}
		}
	}

	err = tx.EventRepo().Delete(ctx, id)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to delete events", zap.Error(err), zap.Int("eventID", id))
		return err
	}

	err = s.nats.Publish(stream.Audit, subject.AuditAdd, admin.Audit{
		AdminId:   a.ID,
		IpAddress: a.IP,
		EventName: admin.DeleteNotificationsEvent,
		OldData:   *selectedEvent,
		NewData:   event.Event{},
		CreatedAt: time.Now(),
	})
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to publish audit event", zap.Error(err))
	}

	return nil
}
