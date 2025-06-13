package event

import (
	"context"
	"encoding/csv"
	"errors"
	"mime/multipart"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"notifications/internal/api/resp"
	"notifications/internal/api/transport/broker/stream"
	"notifications/internal/api/transport/broker/subject"
	"notifications/internal/repo/repomodel"
	userrepo "notifications/internal/repo/user"
	"notifications/internal/service/admin"
	"notifications/pkg/lib/fileman"
	"notifications/pkg/util/strset"
)

func (s *service) LoadUsers(ctx context.Context, a admin.Admin, id int, file multipart.File, fileHeader *multipart.FileHeader) (*Event, error) {
	fileExt := fileman.GetFileExt(fileHeader.Filename)
	if fileExt != fileman.Csv {
		s.logger.Warning("invalid file extension", zap.String("fileName", fileHeader.Filename), zap.Int("id", id))
		return nil, resp.Wrap(resp.ErrBadRequest, "incorrect file type")
	}

	selectedEvent, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during getting event", zap.Error(err), zap.Int("id", id))
			return nil, err
		}
		return nil, resp.Wrap(resp.ErrNotFound, "event not found or not active")
	}

	var oldEvent = *selectedEvent

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		s.logger.Warning("cannot read csv header", zap.Error(err), zap.Int("id", id))
		return nil, resp.Wrap(resp.ErrBadRequest, "cannot read csv file")
	}

	if len(records) == 0 {
		s.logger.Warning("csv file is empty", zap.Int("id", id))
		return nil, resp.Wrap(resp.ErrBadRequest, "csv file is empty")
	}

	header := records[0]
	// if csv file has more headers but userID, it means file is invalid
	if len(header) != 1 || header[0] != _userIDCsvHeader {
		s.logger.Warning("incorrect csv header", zap.Any("header", header), zap.Int("id", id))
		return nil, resp.Wrap(resp.ErrBadRequest, "invalid csv header")
	}

	var cacheKey = _topicSubCacheKey + strset.IntToStr(id)

	if s.cache.Exists(ctx, cacheKey) {
		return nil, resp.Wrap(resp.ErrBadRequest, "user loading in process")
	}

	data := records[1:]

	err = s.cache.SetObj(ctx, cacheKey, data, 0)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("cannot save data in cache", zap.Error(err), zap.Int("id", id))
		return nil, err
	}

	var event = new(Event)
	delete(selectedEvent.ExtraData, _successCountKey)
	delete(selectedEvent.ExtraData, _failedCountKey)
	delete(selectedEvent.ExtraData, _failedReasonKey)
	selectedEvent.Status = _loadingUsers

	event.toService(selectedEvent)
	event.SubscribeAll = false
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

func (s *service) SubscribeUsers(ctx context.Context, event *Event) (err error) {
	s.logger.Info("SubscribeUsers start", zap.Int("eventID", event.ID))

	var (
		data     [][]string
		cacheKey = _topicSubCacheKey + strset.IntToStr(event.ID)
	)

	defer func() {
		if err != nil {
			event.Status = _failedLoading
			event.ExtraData["reason"] = err.Error()
			_, sErr := s.eventRepo.Update(ctx, toRepo(event))
			err = errors.Join(err, sErr)
		}
	}()

	err = s.cache.Get(context.Background(), cacheKey, &data)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("cannot get data from cache", zap.Error(err), zap.Int("eventID", event.ID))
		return err
	}

	var (
		jobCh    = make(chan []string)
		resultCh = make(chan ChunkResult)
		wg       = new(sync.WaitGroup)
	)

	const (
		_workerCount = 100
		_chunkSize   = 1000
	)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// By reading and processing the CSV file in chunks, we never hold more rows in memory than necessary
	// Only a single chunk is in memory (plus whatever our workers hold) at any time
	go func() {
		defer close(jobCh)
		var batch = make([]string, 0, _chunkSize)

		for _, record := range data {
			if len(record) == 0 {
				continue
			}

			batch = append(batch, record[0])
			// if the batch is equal to chunkSize (1000), send it to the workers and reset the batch
			if len(batch) == _chunkSize {
				select {
				case <-ctx.Done():
					s.logger.Error("context canceled", zap.Error(ctx.Err()), zap.Int("id", event.ID))
					return
				case jobCh <- batch:
				}
				batch = make([]string, 0, _chunkSize)
			}
		}
		// if there are records left in the batch or the records are less than the chunkSize, send them to the workers
		if len(batch) > 0 {
			jobCh <- batch
		}
	}()

	wg.Add(_workerCount)
	for i := 1; i <= _workerCount; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					s.logger.Error("context canceled", zap.Error(ctx.Err()), zap.Int("id", event.ID))
					return
				case chunk, ok := <-jobCh:
					if !ok {
						return
					}
					resultCh <- s.processChunk(ctx, event.ID, event.Topic, chunk)
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

	err = s.cache.Delete(ctx, cacheKey)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("cannot delete data from cache", zap.Error(err), zap.Int("eventID", event.ID))
		return err
	}

	s.logger.Info("SubscribeUsers end", zap.Int("eventID", event.ID))

	return nil
}

func (s *service) processChunk(ctx context.Context, eventID int, topic string, chunk []string) (result ChunkResult) {
	users, err := s.userRepo.GetTokensByUserIDs(ctx, chunk)
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

	return result
}

func (s *service) groupUsersByTopic(baseTopic string, users []userrepo.User) (map[string][]string, []userrepo.EventRelation) {
	var (
		topics    = make(map[string][]string)
		relations = make([]userrepo.EventRelation, 0, len(users))
	)

	for _, user := range users {
		if strset.IsEmpty(user.Token) {
			continue
		}
		userTopic := buildTopic(baseTopic, user.Language)
		topics[userTopic] = append(topics[userTopic], user.Token)
		relations = append(relations, userrepo.EventRelation{
			UserID: user.UserID,
			Lang:   user.Language,
		})
	}

	return topics, relations
}
