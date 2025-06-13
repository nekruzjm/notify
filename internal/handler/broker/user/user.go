package user

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"

	"notifications/internal/service/user"
)

func (h *handler) UserCreated(msg jetstream.Msg) {
	h.logger.Info("UserCreated msg", zap.ByteString("data", msg.Data()))

	var (
		ctx  = context.Background()
		data = struct {
			UserID            int    `json:"userID"`
			CountryID         int8   `json:"countryID"`
			Phone             string `json:"phone"`
			PersonExternalRef string `json:"personExternalRef"`
			Token             string `json:"token"`
			Status            string `json:"status"`
			Language          string `json:"language"`
		}{}
	)

	err := sonic.Unmarshal(msg.Data(), &data)
	if err != nil {
		h.logger.Error("sonic.Unmarshal error", zap.Error(err))
		return
	}

	err = h.service.CreateUser(ctx, user.User{
		UserID:            data.UserID,
		CountryID:         data.CountryID,
		Phone:             data.Phone,
		PersonExternalRef: data.PersonExternalRef,
		Token:             data.Token,
		Status:            data.Status,
		Language:          data.Language,
	})
	if err != nil {
		h.logger.Error("CreateUser error", zap.Error(err))
		return
	}

	err = msg.Ack()
	if err != nil {
		h.logger.Error("msg ack error", zap.Error(err))
		return
	}
}

func (h *handler) StatusUpdated(msg jetstream.Msg) {
	h.logger.Info("StatusUpdated msg", zap.ByteString("data", msg.Data()))

	var (
		ctx  = context.Background()
		data = struct {
			UserID int    `json:"userID"`
			Status string `json:"status"`
		}{}
	)

	err := sonic.Unmarshal(msg.Data(), &data)
	if err != nil {
		h.logger.Error("sonic.Unmarshal error", zap.Error(err))
		return
	}

	err = h.service.UpdateStatus(ctx, data.UserID, data.Status)
	if err != nil {
		h.logger.Error("UpdateStatus error", zap.Error(err))
		return
	}

	err = msg.Ack()
	if err != nil {
		h.logger.Error("msg ack error", zap.Error(err))
		return
	}
}

func (h *handler) PhoneUpdated(msg jetstream.Msg) {
	h.logger.Info("PhoneUpdated msg", zap.ByteString("data", msg.Data()))

	var (
		ctx  = context.Background()
		data = struct {
			UserID int    `json:"userID"`
			Phone  string `json:"phone"`
		}{}
	)

	err := sonic.Unmarshal(msg.Data(), &data)
	if err != nil {
		h.logger.Error("sonic.Unmarshal error", zap.Error(err))
		return
	}

	err = h.service.UpdatePhone(ctx, data.UserID, data.Phone)
	if err != nil {
		h.logger.Error("UpdatePhone error", zap.Error(err))
		return
	}

	err = msg.Ack()
	if err != nil {
		h.logger.Error("msg ack error", zap.Error(err))
		return
	}
}

func (h *handler) PersonExternalRefUpdated(msg jetstream.Msg) {
	h.logger.Info("PersonExternalRefUpdated msg", zap.ByteString("data", msg.Data()))

	var (
		ctx  = context.Background()
		data = struct {
			UserID            int    `json:"userID"`
			PersonExternalRef string `json:"personExternalRef"`
		}{}
	)

	err := sonic.Unmarshal(msg.Data(), &data)
	if err != nil {
		h.logger.Error("sonic.Unmarshal error", zap.Error(err))
		return
	}

	err = h.service.UpdatePersonExternalRef(ctx, data.UserID, data.PersonExternalRef)
	if err != nil {
		h.logger.Error("UpdatePersonExternalRef error", zap.Error(err))
		return
	}

	err = msg.Ack()
	if err != nil {
		h.logger.Error("msg ack error", zap.Error(err))
		return
	}
}

func (h *handler) SettingsUpdated(msg jetstream.Msg) {
	h.logger.Info("SettingsUpdated msg", zap.ByteString("data", msg.Data()))

	var (
		ctx  = context.Background()
		data = struct {
			UserID    int    `json:"userID"`
			Language  string `json:"language"`
			IsEnabled *bool  `json:"isEnabled"`
		}{}
	)

	err := sonic.Unmarshal(msg.Data(), &data)
	if err != nil {
		h.logger.Error("sonic.Unmarshal error", zap.Error(err))
		return
	}

	err = h.service.UpdateUserSettings(ctx, data.UserID, data.Language, data.IsEnabled)
	if err != nil {
		h.logger.Error("UpdateUserSettings error", zap.Error(err))
		return
	}

	err = msg.Ack()
	if err != nil {
		h.logger.Error("msg ack error", zap.Error(err))
		return
	}
}

func (h *handler) TokenUpdated(msg jetstream.Msg) {
	h.logger.Info("TokenUpdated msg", zap.ByteString("data", msg.Data()))

	var (
		ctx  = context.Background()
		data = struct {
			UserID int    `json:"userID"`
			Token  string `json:"token"`
		}{}
	)

	err := sonic.Unmarshal(msg.Data(), &data)
	if err != nil {
		h.logger.Error("sonic.Unmarshal error", zap.Error(err))
		return
	}

	err = h.service.UpdateToken(ctx, data.UserID, data.Token)
	if err != nil {
		h.logger.Error("UpdateToken error", zap.Error(err))
		return
	}

	err = msg.Ack()
	if err != nil {
		h.logger.Error("msg ack error", zap.Error(err))
		return
	}
}
