package nats

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func (n *natsConn) Reply(subj, qGroup string, handler nats.MsgHandler) {
	_, err := n.conn.QueueSubscribe(subj, qGroup, handler)
	if err != nil {
		n.logger.Error("err on js.conn.QueueSubscribe", zap.String("subject", subj), zap.Error(err))
		return
	}

	n.logger.Info("Listening on Request-Reply", zap.String("subject", subj))
}
