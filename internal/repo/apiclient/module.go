package apiclient

import (
	"context"

	"go.uber.org/fx"

	"notifications/internal/db"
)

var Module = fx.Provide(New)

type Repo interface {
	GetByUserID(ctx context.Context, userID string) (APIClient, error)
}

type Params struct {
	fx.In

	DB db.QueryExecutor
}

type repo struct {
	db db.QueryExecutor
}

func New(p Params) Repo {
	return &repo{
		db: p.DB,
	}
}
