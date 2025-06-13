package push

import (
	"context"

	"go.uber.org/fx"

	"notifications/internal/db"
)

var Module = fx.Provide(New)

type Repo interface {
	Insert(context.Context, *Push) (*Push, error)
	DeleteByIDs(ctx context.Context, ids []int) error
	Clean(ctx context.Context) error
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
