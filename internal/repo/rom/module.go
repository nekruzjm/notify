package rom

import (
	"context"

	"go.uber.org/fx"

	"notifications/internal/db"
)

var Module = fx.Provide(New)

type Repo interface {
	InsertInbox(context.Context, *Inbox) error
	BatchInsert(ctx context.Context, eventID int, userIDs []int, event *Event) error
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
