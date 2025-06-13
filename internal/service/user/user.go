package user

import (
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"

	"notifications/internal/repo/repomodel"
	"notifications/internal/repo/user"
	"notifications/pkg/util/strset"
)

func (s *service) CreateUser(ctx context.Context, request User) error {
	selectedUser, err := s.userRepo.GetByUserID(ctx, request.UserID)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during getting user", zap.Error(err), zap.Int("userID", request.UserID))
			return err
		}
	}

	if selectedUser == nil {
		err = s.userRepo.Create(ctx, user.User{
			UserID:            request.UserID,
			Phone:             request.Phone,
			PersonExternalRef: request.PersonExternalRef,
			Token:             request.Token,
			Status:            request.Status,
			Language:          request.Language,
			CountryID:         request.CountryID,
		})
		if err != nil {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during creating user", zap.Error(err), zap.Any("request", request))
			return err
		}
	}

	return nil
}

func (s *service) UpdateToken(ctx context.Context, userID int, token string) error {
	selectedUser, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during getting user", zap.Error(err), zap.Int("userID", userID))
			return err
		}
		s.logger.Warning("user not found", zap.Int("userID", userID))
		return nil
	}

	if selectedUser.Token == token {
		return nil
	}

	err = s.userRepo.UpdateToken(ctx, userID, token)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("err occurred during updating firebase token", zap.Error(err), zap.Int("userID", userID), zap.String("token", token))
		return err
	}

	if strset.IsEmpty(selectedUser.Token) {
		return nil
	}

	eventRelations, err := s.userRepo.GetTopicsByUserID(ctx, userID)
	if err != nil && errors.Is(err, repomodel.ErrNotFound) {
		return nil
	}

	for _, rel := range eventRelations {
		topic := buildTopic(rel.Topic, rel.Lang)
		response, err := s.fcmTopicMan.Unsubscribe(ctx, selectedUser.Token, topic)
		if err != nil {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during unsubscribing from topic", zap.Error(err), zap.String("topic", topic))
		}

		s.logger.Info("Unsubscribed from topic", zap.Any("response", response), zap.String("topic", topic))

		response, err = s.fcmTopicMan.Subscribe(ctx, token, topic)
		if err != nil {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during unsubscribing from topic", zap.Error(err), zap.String("topic", topic))
		}

		s.logger.Info("Subscribed to topic", zap.Any("response", response), zap.String("topic", topic))
	}

	return nil
}

func (s *service) UpdateUserSettings(ctx context.Context, userID int, language string, isEnabled *bool) error {
	selectedUser, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during getting user", zap.Error(err), zap.Int("userID", userID))
			return err
		}
		return nil
	}

	if isEnabled == nil && selectedUser.Language == language {
		return nil
	}

	if (selectedUser.Language != language && !strset.IsEmpty(language)) && !strset.IsEmpty(selectedUser.Token) {
		eventRelations, _ := s.userRepo.GetTopicsByUserID(ctx, userID)

		for _, rel := range eventRelations {
			topic := buildTopic(rel.Topic, rel.Lang)
			response, err := s.fcmTopicMan.Unsubscribe(ctx, selectedUser.Token, topic)
			if err != nil {
				s.sentry.CaptureException(err)
				s.logger.Error("err occurred during unsubscribing from topic", zap.Error(err), zap.String("topic", topic))
			}

			s.logger.Info("Unsubscribed from topic", zap.Any("response", response), zap.String("topic", topic))

			newTopic := replaceLastSegment(topic, language)
			response, err = s.fcmTopicMan.Subscribe(ctx, selectedUser.Token, newTopic)
			if err != nil {
				s.sentry.CaptureException(err)
				s.logger.Error("err occurred during subscribing from topic", zap.Error(err), zap.String("topic", topic))
			}

			s.logger.Info("Subscribed to topic", zap.Any("response", response), zap.String("newTopic", newTopic))
		}

		if len(eventRelations) != 0 {
			err = s.userRepo.UpdateRelationLanguage(ctx, userID, language)
			if err != nil {
				s.sentry.CaptureException(err)
				s.logger.Error("err occurred during updating relation language", zap.Error(err), zap.Int("userID", userID))
			}
		}
	}

	if isEnabled != nil {
		selectedUser.PushEnabled = *isEnabled
	}
	if selectedUser.Language != language && !strset.IsEmpty(language) {
		selectedUser.Language = language
	}

	err = s.userRepo.UpdateUserSettings(ctx, userID, selectedUser.Language, selectedUser.PushEnabled)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("err occurred during updating push setting", zap.Error(err), zap.Int("userID", userID))
		return err
	}

	return nil
}

func (s *service) UpdateStatus(ctx context.Context, userID int, status string) error {
	if strset.IsEmpty(status) {
		return nil
	}

	selectedUser, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during getting user", zap.Error(err), zap.Int("userID", userID))
			return err
		}
		return nil
	}

	if selectedUser.Status == status {
		return nil
	}

	if status == _deleted {
		if err = s.userRepo.DeleteRelationByUserID(ctx, userID); err != nil {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during deleting relations by userID", zap.Error(err), zap.Int("userID", userID))
		}

		if err = s.userRepo.DeleteByUserID(ctx, userID); err != nil {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during deleting user by id", zap.Error(err), zap.Int("userID", userID))
			return err
		}

		return nil
	}

	err = s.userRepo.UpdateStatus(ctx, userID, status)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("err occurred during updating status", zap.Error(err), zap.Int("userID", userID), zap.String("status", status))
		return err
	}

	return nil
}

func (s *service) UpdatePhone(ctx context.Context, userID int, phone string) error {
	if strset.IsEmpty(phone) {
		return nil
	}

	selectedUser, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during getting user", zap.Error(err), zap.Int("userID", userID))
			return err
		}
		return nil
	}

	if selectedUser.Phone == phone {
		return nil
	}

	err = s.userRepo.UpdatePhone(ctx, userID, phone)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("err occurred during updating status", zap.Error(err), zap.Int("userID", userID), zap.String("phone", phone))
		return err
	}

	return nil
}

func (s *service) UpdatePersonExternalRef(ctx context.Context, userID int, personExternalRef string) error {
	if strset.IsEmpty(personExternalRef) {
		return nil
	}

	selectedUser, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, repomodel.ErrNotFound) {
			s.sentry.CaptureException(err)
			s.logger.Error("err occurred during getting user", zap.Error(err), zap.Int("userID", userID))
			return err
		}
		return nil
	}

	if selectedUser.PersonExternalRef == personExternalRef {
		return nil
	}

	err = s.userRepo.UpdatePersonExternalRef(ctx, userID, personExternalRef)
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("err occurred during updating status", zap.Error(err), zap.Int("userID", userID), zap.String("personExternalRef", personExternalRef))
		return err
	}

	return nil
}

// replaceLastSegment replaces the last segment of the topic with the given language
// if the topic already has the language, it returns the topic as is
func replaceLastSegment(topic, lang string) string {
	lastUnderscore := strings.LastIndex(topic, "_")
	if lastUnderscore == -1 {
		return topic
	}
	return topic[:lastUnderscore+1] + lang
}

// buildTopic builds a topic with the given language
// if the topic already has the language, it returns the topic as is
func buildTopic(topic, lang string) string {
	if strings.HasSuffix(topic, lang) {
		return topic
	}
	return topic + "_" + lang
}
