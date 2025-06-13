package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"notifications/internal/lib/ctxman"
	"notifications/pkg/util/strset"
)

type QueryExecutor interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
	Begin(context.Context) (pgx.Tx, error)
	CopyFrom(ctx context.Context, identifier pgx.Identifier, cols []string, rows pgx.CopyFromSource) (int64, error)
}

func (c *conn) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	c.logger.Info("Exec query",
		zap.String("service", c.serviceName),
		zap.String("query", strset.RemoveSpecialChars(sql)),
		zap.Any("arguments", args))

	return c.getPool(ctx).Exec(ctx, sql, args...)
}

func (c *conn) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if len(args) < 100 {
		c.logger.Info("Query",
			zap.String("service", c.serviceName),
			zap.String("query", strset.RemoveSpecialChars(sql)),
			zap.Any("arguments", args))
	}

	return c.getPool(ctx).Query(ctx, sql, args...)
}

func (c *conn) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	c.logger.Info("QueryRow",
		zap.String("service", c.serviceName),
		zap.String("query", strset.RemoveSpecialChars(sql)),
		zap.Any("arguments", args))

	return c.getPool(ctx).QueryRow(ctx, sql, args...)
}

func (c *conn) SendBatch(ctx context.Context, batch *pgx.Batch) pgx.BatchResults {
	c.logger.Info("SendBatch query", zap.String("service", c.serviceName))
	return c.getPool(ctx).SendBatch(ctx, batch)
}

func (c *conn) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.getPool(ctx).Begin(ctx)
}

func (c *conn) CopyFrom(ctx context.Context, identifier pgx.Identifier, cols []string, rows pgx.CopyFromSource) (int64, error) {
	c.logger.Info("CopyFrom query", zap.Any("identifier", identifier), zap.Any("columns", cols))
	return c.getPool(ctx).CopyFrom(ctx, identifier, cols, rows)
}

func (c *conn) getPool(ctx context.Context) *pgxpool.Pool {
	var ctxInfo = ctxman.Get(ctx)
	return c.pools[c.getPoolKey(ctxInfo.DBName, ctxInfo.IsReplica)]
}

func (c *conn) getPoolKey(dbName string, isReplica bool) (key string) {
	var replKey = master
	if isReplica {
		replKey = replica
	}

	switch dbName {
	case Notifications:
		key = fmt.Sprintf("databases.%s.%s.%s", c.serviceName, _notificationsDB, replKey)
	case Rom:
		key = fmt.Sprintf("databases.%s.%s.%s", c.serviceName, _romDB, replKey)
	default:
		key = fmt.Sprintf("databases.%s.%s.%s", c.serviceName, _notificationsDB, replKey)
	}
	return
}
