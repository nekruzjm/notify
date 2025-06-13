package push

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"

	"notifications/internal/service/push"
)

func (h *handler) Sent(msg jetstream.Msg) {
	h.logger.Info("msg Sent", zap.ByteString("data", msg.Data()))

	var (
		ctx     = context.Background()
		message struct {
			UserID     int               `json:"userID"`
			Token      string            `json:"token"`
			Data       map[string]string `json:"data"`
			ShowInFeed bool              `json:"showInFeed"`
		}
	)

	err := sonic.Unmarshal(msg.Data(), &message)
	if err != nil {
		h.logger.Error("sonic.Unmarshal error", zap.Error(err), zap.ByteString("data", msg.Data()))
		return
	}

	err = msg.Ack()
	if err != nil {
		h.logger.Error("msg ack error", zap.Error(err), zap.Int("userID", message.UserID))
		return
	}

	var request = new(push.Request)
	request.InternalRequest.UserID = message.UserID
	request.InternalRequest.Token = message.Token
	request.InternalRequest.Data = message.Data
	request.ShowInFeed = message.ShowInFeed
	request.IsInternal = true

	_, err = h.service.Send(ctx, request)
	if err != nil {
		h.logger.Error("SendInternal error", zap.Error(err), zap.Int("userID", message.UserID))
		return
	}
}

func (h *handler) SyncSent(msg *nats.Msg) {
	h.logger.Info("msg SyncSent", zap.ByteString("data", msg.Data))

	var (
		ctx       = context.Background()
		err       error
		messageID string
		message   struct {
			UserID int               `json:"userID"`
			Token  string            `json:"token"`
			Data   map[string]string `json:"data"`
		}
		response struct {
			MessageID string `json:"messageID"`
			Error     string `json:"error"`
		}
	)

	defer func() {
		response.MessageID = messageID
		if err != nil {
			response.Error = err.Error()
		}

		respBytes, err := sonic.Marshal(response)
		if err != nil {
			h.logger.Error("sonic.Marshal error", zap.Error(err), zap.Any("response", response))
			return
		}

		_ = msg.Respond(respBytes)
	}()

	err = sonic.Unmarshal(msg.Data, &message)
	if err != nil {
		h.logger.Error("sonic.Unmarshal error", zap.Error(err), zap.ByteString("data", msg.Data))
		return
	}

	var request = new(push.Request)
	request.InternalRequest.UserID = message.UserID
	request.InternalRequest.Token = message.Token
	request.InternalRequest.Data = message.Data
	request.IsInternal = true
	request.Sync = true

	messageID, err = h.service.Send(ctx, request)
	if err != nil {
		h.logger.Error("SendInternal error", zap.Error(err), zap.Int("userID", message.UserID))
		return
	}
}

func (h *handler) Clean(msg jetstream.Msg) {
	err := msg.Ack()
	if err != nil {
		h.logger.Error("msg ack error", zap.Error(err))
		return
	}

	h.service.Clean()
}
