package event

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"notifications/internal/api/resp"
	"notifications/internal/api/transport/broker/stream"
	"notifications/internal/api/transport/broker/subject"
	"notifications/internal/repo/repomodel"
	"notifications/internal/service/admin"
	"notifications/pkg/util/strset"
)

func (s *service) LoadAllUsers(ctx context.Context, a admin.Admin, id int) (*Event, error) {
	selectedEvent, err := s.eventRepo.GetByID(ctx, id)
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

	var event = new(Event)
	delete(selectedEvent.ExtraData, _successCountKey)
	delete(selectedEvent.ExtraData, _failedCountKey)
	delete(selectedEvent.ExtraData, _failedReasonKey)
	selectedEvent.Status = _loadingUsers

	event.toService(selectedEvent)
	event.SubscribeAll = true
	err = s.nats.Publish(stream.Notifications, subject.NotificationsTopicUsersSubscribed, event)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("cannot publish message", zap.Error(err), zap.Int("id", id))
		return nil, err
	}

	_, err = s.eventRepo.Update(ctx, selectedEvent)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("cannot update event", zap.Error(err), zap.Int("id", id))
		return nil, err
	}

	err = s.nats.Publish(stream.Audit, subject.AuditAdd, admin.Audit{
		AdminId:   a.ID,
		IpAddress: a.IP,
		EventName: admin.LoadUsersEvent,
		OldData:   oldEvent,
		NewData:   *selectedEvent,
		CreatedAt: time.Now(),
	})
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("failed to publish audit event", zap.Error(err))
	}

	return event, nil
}

func (s *service) SubscribeAllUsers(ctx context.Context, event *Event) (err error) {
	s.logger.Info("SubscribeAllUsers start", zap.Int("eventID", event.ID))

	defer func() {
		if err != nil {
			event.Status = _failedLoading
			event.ExtraData["reason"] = err.Error()
			_, err = s.eventRepo.Update(ctx, toRepo(event))
		}
	}()

	var (
		resultCh = make(chan ChunkResult)
		wg       = new(sync.WaitGroup)
	)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	const _workerCount = 100

	wg.Add(_workerCount)
	for i := 1; i <= _workerCount; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					s.logger.Error("context canceled", zap.Error(ctx.Err()), zap.Int("id", event.ID))
					return
				default:
					res := s.processUsers(ctx, event.ID, event.Topic)
					if errors.Is(res.Err, repomodel.ErrNotFound) {
						return
					}
					resultCh <- res
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var response PrepareUsersResponse

	for res := range resultCh {
		if res.Err == nil {
			atomic.AddInt64(&response.SuccessCount, int64(res.SuccessCount))
			atomic.AddInt64(&response.FailedCount, int64(res.FailedCount))
			atomic.AddInt64(&response.ErrCount, int64(res.ErrCount))
		}
	}

	event.Status = _active
	event.ExtraData[_successCountKey] = strset.IntToStr(int(response.SuccessCount))
	event.ExtraData[_failedCountKey] = strset.IntToStr(int(response.FailedCount))

	if response.ErrCount > (response.FailedCount / 2) {
		event.ExtraData[_failedReasonKey] = "most of the failed to subscribe users have inactive tokens"
	}

	_, err = s.eventRepo.Update(ctx, toRepo(event))
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("cannot update event", zap.Error(err), zap.Int("eventID", event.ID))
		return err
	}

	err = s.cache.Delete(ctx, ":process-users")
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("err from Delete", zap.Error(err), zap.Int("eventID", event.ID))
		return err
	}

	s.logger.Info("SubscribeAllUsers end", zap.Int("eventID", event.ID))

	return nil
}

func (s *service) processUsers(ctx context.Context, eventID int, topic string) (result ChunkResult) {
	var (
		cacheKey = ":process-users"
		lastID   int
	)

	err := s.cache.Get(ctx, cacheKey, &lastID)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			s.logger.Error("err Get", zap.Error(err), zap.Int("eventID", eventID))
			return ChunkResult{Err: err}
		}
	}

	users, err := s.userRepo.GetTokensWithLimit(ctx, lastID)
	if err != nil {
		s.logger.Error("err from GetTokensByUserIDs", zap.Error(err), zap.Int("eventID", eventID))
		return ChunkResult{Err: err}
	}

	topics, relations := s.groupUsersByTopic(topic, users)

	for topicLang, tokens := range topics {
		if len(tokens) == 0 {
			continue
		}

		response, err := s.fcmTopicMan.SubscribeTokens(ctx, tokens, topicLang)
		if err != nil {
			s.logger.Error("err from SubscribeTokens", zap.Error(err), zap.String("topic", topicLang), zap.Int("eventID", eventID))
			return ChunkResult{Err: err}
		}

		if response.Errors != nil {
			for _, errDetail := range response.Errors {
				if errDetail.Reason == _fcmTokenErr {
					result.ErrCount++
				}
			}
		}
		result.SuccessCount += response.SuccessCount
		result.FailedCount += response.FailureCount
	}

	_, err = s.userRepo.BatchInsert(ctx, eventID, relations)
	if err != nil {
		s.logger.Error("err from BatchInsert", zap.Error(err), zap.Int("eventID", eventID))
		return ChunkResult{Err: err}
	}

	lastID = users[len(users)-1].UserID
	err = s.cache.Set(ctx, cacheKey, lastID, 0)
	if err != nil {
		s.logger.Error("err Set", zap.Error(err), zap.Int("eventID", eventID))
		return ChunkResult{Err: err}
	}

	return result
}
