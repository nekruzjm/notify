package event

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"

	"notifications/internal/lib/language"
	"notifications/internal/repo/repomodel"
	"notifications/pkg/util/strset"
)

func (s *service) UnsubscribeUsers(ctx context.Context, event *Event) error {
	s.logger.Info("UnsubscribeUsers start", zap.Int("eventID", event.ID))

	userIDs, err := s.userRepo.GetUserIDsByEventID(ctx, event.ID)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during getting user relations", zap.Error(err), zap.Int("eventID", event.ID))
			return err
		}
		return nil
	}

	var userIDsStr = make([]string, 0, len(userIDs))
	for _, id := range userIDs {
		userIDsStr = append(userIDsStr, strset.IntToStr(id))
	}

	users, err := s.userRepo.GetTokensByUserIDs(ctx, userIDsStr)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during getting user ids", zap.Error(err), zap.Int("eventID", event.ID))
			return err
		}
		return nil
	}

	var tokens = make([]string, 0, len(users))
	for i := 0; i < len(users); i++ {
		tokens = append(tokens, users[i].Token)
	}

	var (
		wg        = new(sync.WaitGroup)
		languages = language.GetAll()
	)

	for _, lang := range languages {
		wg.Add(1)
		go func() {
			defer wg.Done()

			topicLang := buildTopic(event.Topic, lang)
			_, err = s.fcmTopicMan.UnsubscribeTokens(ctx, tokens, topicLang)
			if err != nil {
				s.sentry.CaptureException(err)
				s.logger.Error("err from UnsubscribeTokens", zap.Error(err), zap.String("topic", topicLang), zap.Int("eventID", event.ID))
				return
			}
		}()
	}

	wg.Wait()

	err = s.userRepo.DeleteRelationByEventID(ctx, event.ID)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("err occurred during deleting user relations", zap.Error(err), zap.Int("eventID", event.ID))
		return err
	}

	s.logger.Info("UnsubscribeUsers end", zap.Int("eventID", event.ID))

	return nil
}
