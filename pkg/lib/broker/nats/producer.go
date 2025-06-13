package nats

import (
	"errors"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func (n *natsConn) Publish(stream, subj string, msg any) (err error) {
	var body []byte
	if msg != nil {
		body, err = sonic.Marshal(msg)
		if err != nil {
			n.logger.Error("err occurred during serialization", zap.Error(err), zap.Any("msg", msg))
			return err
		}
	}

	_, err = n.js.Stream(n.ctx, stream)
	if err != nil {
		if !errors.Is(err, jetstream.ErrStreamNotFound) {
			n.logger.Error("can't get stream", zap.Error(err), zap.String("stream", stream))
		}
		return err
	}

	_, err = n.js.Publish(n.ctx, subj, body)
	if err != nil {
		n.logger.Error("err occurred during publish message", zap.Error(err), zap.String("stream", stream), zap.String("subject", subj))
		return err
	}

	n.logger.Info("Message published", zap.String("stream", stream), zap.String("subject", subj), zap.ByteString("data", body))

	return nil
}
