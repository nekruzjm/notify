package push

import (
	"context"
	"errors"
	"strings"

	"firebase.google.com/go/v4/messaging"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"

	"notifications/internal/api/resp"
	"notifications/internal/api/transport/broker/stream"
	"notifications/internal/api/transport/broker/subject"
	"notifications/internal/lib/language"
	"notifications/internal/repo/push"
	"notifications/internal/repo/repomodel"
	"notifications/internal/repo/rom"
	"notifications/internal/repo/user"
	"notifications/pkg/lib/broker/nats"
	"notifications/pkg/lib/notifier/firebase"
	"notifications/pkg/lib/observer/logger"
	"notifications/pkg/lib/observer/sentry"
	"notifications/pkg/util/strset"
)

type external struct {
	logger      logger.Logger
	sentry      sentry.Sentry
	nats        nats.Event
	fcmSender   firebase.Sender
	userRepo    user.Repo
	pushRepo    push.Repo
	romRepo     rom.Repo
	idGenerator *snowflake.Node
}

func (e *external) Clean() {}

func (e *external) Send(ctx context.Context, request *Request) (string, error) {
	err := request.validate()
	if err != nil {
		e.logger.Warning("invalid request", zap.Error(err), zap.Any("request", request), zap.String("requestID", request.ExternalRequest.ID))
		return "", resp.Wrap(resp.ErrBadRequest, err.Error())
	}

	var selectedUser *user.User
	if !strset.IsEmpty(request.ExternalRequest.Phone) {
		selectedUser, err = e.userRepo.GetActiveByPhone(ctx, request.ExternalRequest.Phone)
	} else {
		selectedUser, err = e.userRepo.GetActiveByPersonExternalRef(ctx, request.ExternalRequest.PersonExternalRef)
	}
	if err != nil {
		if errors.Is(err, repomodel.ErrNotFound) {
			return "", resp.Wrap(resp.ErrNotFound, "user not found")
		}
		e.logger.Error("cannot get user", zap.Error(err), zap.String("requestID", request.ExternalRequest.ID))
		return "", err
	}

	if request.ShowInFeed {
		return e.sendStateful(ctx, selectedUser, request)
	}

	return e.sendStateless(ctx, selectedUser, request)
}

func (e *external) sendStateless(ctx context.Context, user *user.User, request *Request) (string, error) {
	var (
		title = request.ExternalRequest.Title.Get(user.Language)
		body  = request.ExternalRequest.Body.Get(user.Language)
	)
	if strset.IsEmpty(title) {
		title = request.ExternalRequest.Title.Get(language.RU)
	}
	if strset.IsEmpty(body) {
		body = request.ExternalRequest.Body.Get(language.RU)
	}

	var (
		message = new(messaging.Message)
		data    = map[string]string{
			_title:   title,
			_comment: title,
			_message: body,
			_badge:   _badge0,
		}
	)

	message.Data = data
	message.Token = user.Token
	firebase.AndroidMSG(message, data, firebase.AndroidHighestPriority)
	firebase.IosMSG(message, data, firebase.ApnsHighestPriority)

	messageID, err := e.fcmSender.SendPush(ctx, message)
	if err != nil {
		if !firebase.IsValidationErr(err) {
			e.sentry.CaptureException(err)
			e.logger.Error("error in fcm.SendPush", zap.Error(err), zap.String("requestID", request.ExternalRequest.ID))
			return "", err
		}
		err = e.nats.Publish(stream.Notifications, subject.NotificationsFcmRegistrationTokenRemoved, user.UserID)
		if err != nil {
			e.sentry.CaptureException(err)
			e.logger.Error("error on publish event", zap.Error(err), zap.String("requestID", request.ExternalRequest.ID))
			return "", err
		}
	}

	messageID = messageID[strings.LastIndex(messageID, _slashDelim)+1:]

	return messageID, nil
}

func (e *external) sendStateful(ctx context.Context, user *user.User, request *Request) (messageID string, err error) {
	savedPush, err := e.pushRepo.Insert(ctx, &push.Push{
		ID:        int(e.idGenerator.Generate().Int64()),
		UserID:    user.UserID,
		Status:    _approved,
		Title:     request.ExternalRequest.Title,
		Body:      request.ExternalRequest.Body,
		Type:      request.ExternalRequest.PushType,
		APIClient: request.ExternalRequest.APIClient,
	})
	if err != nil {
		e.sentry.CaptureException(err)
		e.logger.Error("error in push.Insert", zap.Error(err), zap.String("requestID", request.ExternalRequest.ID))
		return "", err
	}

	err = e.romRepo.InsertInbox(ctx, &rom.Inbox{
		ID:        savedPush.ID,
		UserID:    savedPush.UserID,
		CountryID: user.CountryID,
		Type:      savedPush.Type,
		Title:     savedPush.Title,
		Body:      savedPush.Body,
		ExtraData: map[string]string{},
	})
	if err != nil {
		e.sentry.CaptureException(err)
		e.logger.Error("err in romRepo.InsertInbox", zap.Error(err), zap.String("requestID", request.ExternalRequest.ID))
		return "", err
	}

	if user.Status != _active {
		e.logger.Warning("user status is not active", zap.String("requestID", request.ExternalRequest.ID))
		return _inactiveUserMessageID, nil
	}
	if !user.PushEnabled {
		e.logger.Warning("user push is disabled", zap.String("requestID", request.ExternalRequest.ID))
		return _disabledPushMessageID, nil
	}

	var (
		title = request.ExternalRequest.Title.Get(user.Language)
		body  = request.ExternalRequest.Body.Get(user.Language)
	)
	if strset.IsEmpty(title) {
		title = request.ExternalRequest.Title.Get(language.RU)
	}
	if strset.IsEmpty(body) {
		body = request.ExternalRequest.Body.Get(language.RU)
	}

	data := make(map[string]string)
	data[_title] = title
	data[_comment] = title
	data[_message] = body
	data[_badge] = _badge1
	data[_category] = _defaultCategory
	data[_sectionName] = _defaultSectionName

	message := new(messaging.Message)
	message.Data = data
	message.Token = user.Token
	firebase.AndroidMSG(message, data, firebase.AndroidHighestPriority)
	firebase.IosMSG(message, data, firebase.ApnsHighestPriority)

	msgID, err := e.fcmSender.SendPush(ctx, message)
	if err != nil {
		if !firebase.IsValidationErr(err) {
			e.sentry.CaptureException(err)
			e.logger.Error("error in fcm.SendPush", zap.Error(err), zap.String("requestID", request.ExternalRequest.ID))
			return _fcmPushMessageID, nil
		}
		err = e.nats.Publish(stream.Notifications, subject.NotificationsFcmRegistrationTokenRemoved, user.UserID)
		if err != nil {
			e.sentry.CaptureException(err)
			e.logger.Error("error on publish event", zap.Error(err), zap.String("requestID", request.ExternalRequest.ID))
		}
		return _fcmPushMessageID, nil
	}

	messageID = msgID[strings.LastIndex(msgID, _slashDelim)+1:]

	return messageID, nil
}
