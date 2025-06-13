package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"notifications/pkg/lib/config"
	"notifications/pkg/lib/db/postgres"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Params struct {
	fx.In
	fx.Lifecycle

	Config config.Config
	Logger logger.Logger
}

type databaseParams struct {
	MasterKey  string
	ReplicaKey string
	Master     string
	Replica    string
	MaxConn    int
}

type conn struct {
	config      config.Config
	logger      logger.Logger
	serviceName string
	pools       map[string]*pgxpool.Pool
}

// Database alias for the service
const (
	Notifications = "notifications"
	Rom           = "rom"
)

// specify service names for the databases
const (
	_notifications = "notifications"
)

// dbNames is a list of databases that are used by the service
const (
	_notificationsDB string = "notificationsDB"
	_romDB           string = "romDb"
)

const (
	master  = "master"
	replica = "replica"
	maxConn = "maxConn"
)

func New(p Params) QueryExecutor {
	var (
		dbNames    = []string{_notificationsDB, _romDB}
		connection = &conn{
			serviceName: _notifications,
			config:      p.Config,
			logger:      p.Logger,
			pools:       make(map[string]*pgxpool.Pool),
		}
	)

	for _, dbName := range dbNames {
		var dbParams = connection.getDBParams(_notifications, dbName)

		pgMaster, err := postgres.New(_notifications, dbParams.Master, postgres.MaxPoolSize(dbParams.MaxConn))
		if err != nil {
			p.Logger.Error("err occurred during connection to master db",
				zap.Error(err),
				zap.String("service", _notifications),
				zap.String("database", dbName))
		} else {
			p.Logger.Info("Successfully connected to master db", zap.String("service", _notifications), zap.String("database", dbName))
		}

		connection.pools[dbParams.MasterKey] = pgMaster.Pool
	}

	p.Logger.Info("Service connected to all databases", zap.String("service", _notifications), zap.Any("databases", dbNames))

	p.Lifecycle.Append(
		fx.Hook{
			OnStop: func(_ context.Context) error {
				for _, pool := range connection.pools {
					pool.Close()
				}
				return nil
			},
		},
	)

	return connection
}

func (c *conn) getDBParams(service, dbName string) databaseParams {
	var (
		dnsMasterKey  = fmt.Sprintf("databases.%s.%s.%s", service, dbName, master)
		dnsReplicaKey = fmt.Sprintf("databases.%s.%s.%s", service, dbName, replica)
		maxConnKey    = fmt.Sprintf("databases.%s.%s.%s", service, dbName, maxConn)
	)

	return databaseParams{
		MasterKey:  dnsMasterKey,
		ReplicaKey: dnsReplicaKey,
		Master:     c.config.GetString(dnsMasterKey),
		Replica:    c.config.GetString(dnsReplicaKey),
		MaxConn:    c.config.GetInt(maxConnKey),
	}
}
