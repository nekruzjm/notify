package push

import (
	"context"
	"errors"
	"time"

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
	"notifications/pkg/lib/cache"
	"notifications/pkg/lib/notifier/firebase"
	"notifications/pkg/lib/observer/logger"
	"notifications/pkg/lib/observer/sentry"
	"notifications/pkg/util/strset"
)

type internal struct {
	logger      logger.Logger
	sentry      sentry.Sentry
	nats        nats.Event
	cache       cache.Cache
	fcmSender   firebase.Sender
	userRepo    user.Repo
	pushRepo    push.Repo
	romRepo     rom.Repo
	idGenerator *snowflake.Node
}

func (i *internal) Send(ctx context.Context, request *Request) (string, error) {
	selectedUser, err := i.userRepo.GetByUserID(ctx, request.InternalRequest.UserID)
	if err != nil {
		if errors.Is(err, repomodel.ErrNotFound) {
			i.logger.Warning("user not found", zap.Int("userID", request.InternalRequest.UserID))
			return "", nil
		}
		i.sentry.CaptureException(err)
		i.logger.Error("err occurred during getting user", zap.Error(err), zap.Int("userID", request.InternalRequest.UserID))
		return "", err
	}

	if !strset.IsEmpty(selectedUser.Token) && !strset.IsEmpty(request.InternalRequest.Token) && selectedUser.Token != request.InternalRequest.Token {
		err = i.userRepo.UpdateToken(ctx, request.InternalRequest.UserID, request.InternalRequest.Token)
		if err != nil {
			i.sentry.CaptureException(err)
			i.logger.Error("err occurred during updating token", zap.Error(err), zap.Int("userID", request.InternalRequest.UserID))
		}
		selectedUser.Token = request.InternalRequest.Token
	}

	switch {
	case request.ShowInFeed:
		return i.sendStateful(ctx, selectedUser, request)
	case request.Sync:
		return i.sendStatelessSync(ctx, selectedUser, request)
	default:
		return i.sendStatelessAsync(ctx, selectedUser, request)
	}
}

func (i *internal) sendStateful(ctx context.Context, user *user.User, request *Request) (string, error) {
	var title, body = language.Language{}, language.Language{}
	title.SetAll(request.InternalRequest.Data[_title])
	body.SetAll(request.InternalRequest.Data[_message])

	savedPush, err := i.pushRepo.Insert(ctx, &push.Push{
		ID:        int(i.idGenerator.Generate().Int64()),
		UserID:    user.UserID,
		Status:    _approved,
		Title:     title,
		Body:      body,
		Type:      _push,
		APIClient: _defaultAPIClient,
	})
	if err != nil {
		i.sentry.CaptureException(err)
		i.logger.Error("error in push.Insert", zap.Error(err), zap.Int("userID", request.InternalRequest.UserID))
		return "", err
	}

	err = i.romRepo.InsertInbox(ctx, &rom.Inbox{
		ID:        savedPush.ID,
		UserID:    savedPush.UserID,
		CountryID: user.CountryID,
		Type:      savedPush.Type,
		Title:     savedPush.Title,
		Body:      savedPush.Body,
		ExtraData: map[string]string{},
	})
	if err != nil {
		// delete saved push if error occurred while inserting in rom in order to avoid inconsistency
		// because notification db and rom db are on different servers
		errP := i.pushRepo.DeleteByIDs(ctx, []int{savedPush.ID})
		if errP != nil {
			i.sentry.CaptureException(errP)
			i.logger.Error("error in push.DeleteByIDs", zap.Error(errP), zap.Int("userID", request.InternalRequest.UserID))
		}
		i.sentry.CaptureException(err)
		i.logger.Error("err in romRepo.InsertInbox", zap.Error(err), zap.Int("userID", request.InternalRequest.UserID))
		return "", err
	}

	// if user status is not active or push is disabled, save the state but do not send push
	if user.Status != _active {
		i.logger.Warning("user status is not active", zap.Int("userID", request.InternalRequest.UserID))
		return "", nil
	}
	if !user.PushEnabled || strset.IsEmpty(user.Token) {
		i.logger.Warning("user push is disabled or token is empty", zap.Int("userID", user.UserID))
		return "", nil
	}

	message := new(messaging.Message)
	message.Data = request.InternalRequest.Data
	message.Token = user.Token
	firebase.AndroidMSG(message, request.InternalRequest.Data, firebase.AndroidHighestPriority)
	firebase.IosMSG(message, request.InternalRequest.Data, firebase.ApnsHighestPriority)

	_, err = i.fcmSender.SendPush(ctx, message)
	if err != nil {
		if !firebase.IsValidationErr(err) {
			i.sentry.CaptureException(err)
			i.logger.Error("error in fcm.SendPush", zap.Error(err), zap.Int("userID", request.InternalRequest.UserID))
			return "", err
		}

		err = i.nats.Publish(stream.Notifications, subject.NotificationsFcmRegistrationTokenRemoved, user.UserID)
		if err != nil {
			i.sentry.CaptureException(err)
			i.logger.Error("error on publish event", zap.Error(err), zap.Int("userID", request.InternalRequest.UserID))
			return "", err
		}
	}

	return "", nil
}

func (i *internal) sendStatelessAsync(ctx context.Context, user *user.User, request *Request) (string, error) {
	if !user.PushEnabled || strset.IsEmpty(user.Token) {
		return "", nil
	}

	// implement deduplication for stateless push, to prevent sending multiple push for the same transaction
	trID := request.InternalRequest.Data[_trID]
	if !strset.IsEmpty(trID) {
		var (
			exists   bool
			cacheKey = ":stateless:" + trID
		)
		if err := i.cache.Get(ctx, cacheKey, &exists); err == nil {
			return "", nil
		}
		if err := i.cache.Set(ctx, cacheKey, true, 10*time.Minute); err != nil {
			i.sentry.CaptureException(err)
			i.logger.Error("error in cache.Set", zap.Error(err), zap.Int("userID", request.InternalRequest.UserID))
		}
	}

	var (
		pushType = request.InternalRequest.Data[_pushType]
		message  = &messaging.Message{
			Data:  request.InternalRequest.Data,
			Token: user.Token,
		}
	)

	if pushType != _silent {
		firebase.AndroidMSG(message, request.InternalRequest.Data, firebase.AndroidHighestPriority)
		firebase.IosMSG(message, request.InternalRequest.Data, firebase.ApnsHighestPriority)
	}

	_, err := i.fcmSender.SendPush(ctx, message)
	if err != nil {
		i.logger.Warning("error in fcm.SendPush", zap.Error(err), zap.Int("userID", request.InternalRequest.UserID))

		if !firebase.IsValidationErr(err) {
			i.sentry.CaptureException(err)
			i.logger.Error("error in fcm.SendPush", zap.Error(err), zap.Int("userID", request.InternalRequest.UserID))
			return "", err
		}

		err = i.nats.Publish(stream.Notifications, subject.NotificationsFcmRegistrationTokenRemoved, request.InternalRequest.UserID)
		if err != nil {
			i.sentry.CaptureException(err)
			i.logger.Error("error on publish event", zap.Error(err), zap.Int("userID", request.InternalRequest.UserID))
			return "", err
		}
	}

	return "", nil
}

func (i *internal) sendStatelessSync(ctx context.Context, user *user.User, request *Request) (string, error) {
	if !user.PushEnabled || strset.IsEmpty(user.Token) {
		i.logger.Error("user push is disabled or token is empty", zap.Int("userID", request.InternalRequest.UserID))
		return "", resp.ErrBadRequest
	}

	// implement deduplication for stateless push, to prevent sending multiple push for the same transaction
	trID := request.InternalRequest.Data[_trID]
	if !strset.IsEmpty(trID) {
		var (
			exists   bool
			cacheKey = ":stateless:" + trID
		)
		if err := i.cache.Get(ctx, cacheKey, &exists); err == nil {
			return "", nil
		}
		if err := i.cache.Set(ctx, cacheKey, true, 10*time.Minute); err != nil {
			i.sentry.CaptureException(err)
			i.logger.Error("error in cache.Set", zap.Error(err), zap.Int("userID", request.InternalRequest.UserID))
		}
	}

	var message = &messaging.Message{
		Data:  request.InternalRequest.Data,
		Token: user.Token,
	}

	firebase.AndroidMSG(message, request.InternalRequest.Data, firebase.AndroidHighestPriority)
	firebase.IosMSG(message, request.InternalRequest.Data, firebase.ApnsHighestPriority)

	messageID, err := i.fcmSender.SendPush(ctx, message)
	if err != nil {
		if !firebase.IsValidationErr(err) {
			i.sentry.CaptureException(err)
			i.logger.Error("error in fcm.SendPush", zap.Error(err), zap.Int("userID", request.InternalRequest.UserID))
		} else {
			i.logger.Warning("error in fcm.SendPush", zap.Error(err), zap.Int("userID", request.InternalRequest.UserID))
		}

		nErr := i.nats.Publish(stream.Notifications, subject.NotificationsFcmRegistrationTokenRemoved, request.InternalRequest.UserID)
		if nErr != nil {
			i.sentry.CaptureException(nErr)
			i.logger.Error("error on publish event", zap.Error(nErr), zap.Int("userID", request.InternalRequest.UserID))
			return "", nErr
		}
		return "", err
	}

	return messageID, nil
}

func (i *internal) Clean() {
	err := i.pushRepo.Clean(context.Background())
	if err != nil {
		i.sentry.CaptureException(err)
		i.logger.Error("err occurred during cleaning push", zap.Error(err))
	}
}
