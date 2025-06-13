package nats

import (
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

const (
	_infiniteReconnects = -1
	_reconnectWait      = time.Second * 5
	_pingInterval       = time.Second * 10
	_pingsOutstanding   = 5
)

func (n *natsConn) _connect() (*nats.Conn, error) {
	var (
		url        = n.config.GetString("jetStream.nats.url")
		user       = n.config.GetString("jetStream.nats.user")
		pass       = n.config.GetString("jetStream.nats.pass")
		clientName = n.config.GetString("jetStream.nats.clientName")
	)

	conn, err := nats.Connect(url,
		nats.Name(clientName),
		nats.UserInfo(user, pass),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(_infiniteReconnects),
		nats.ReconnectWait(_reconnectWait),
		nats.PingInterval(_pingInterval),
		nats.MaxPingsOutstanding(_pingsOutstanding),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			n.logger.Warning("client disconnected err", zap.Error(err))
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			n.logger.Info("client reconnected to the server", zap.String("connectedURL", nc.ConnectedUrl()))
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			n.logger.Error("client close handler", zap.Error(nc.LastError()))
		}),
	)
	if err != nil {
		n.logger.Error("err on Connect NATS JetStream", zap.Error(err))
		return nil, err
	}

	n.replicas = len(strings.Split(url, ","))
	n.conn = conn
	n.js, _ = jetstream.New(conn)

	n.logger.Info("Connected to NATS", zap.String("url", conn.ConnectedUrl()), zap.String("clientName", clientName))

	return conn, nil
}
