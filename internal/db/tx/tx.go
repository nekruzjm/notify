package tx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/fx"

	"notifications/internal/db"
	"notifications/internal/repo/event"
	"notifications/internal/repo/user"
)

var Module = fx.Provide(New)

type Transactor interface {
	New() Transactor

	Begin(context.Context) error
	Commit(context.Context) error
	Rollback(ctx context.Context) error

	repo
}

type repo interface {
	UserRepo() user.Repo
	EventRepo() event.Repo
}

type Params struct {
	fx.In

	DB        db.QueryExecutor
	UserRepo  user.Repo
	EventRepo event.Repo
}

type transactor struct {
	db        db.QueryExecutor
	tx        pgx.Tx
	userRepo  user.Repo
	eventRepo event.Repo
}

func New(p Params) Transactor {
	return &transactor{
		db:        p.DB,
		userRepo:  p.UserRepo,
		eventRepo: p.EventRepo,
	}
}

func (t *transactor) New() Transactor {
	return &transactor{
		db:        t.db,
		userRepo:  t.userRepo,
		eventRepo: t.eventRepo,
	}
}

func (t *transactor) Begin(ctx context.Context) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}
	t.tx = tx
	return nil
}

func (t *transactor) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *transactor) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (t *transactor) UserRepo() user.Repo {
	t.userRepo = user.New(user.Params{DB: t.tx})
	return t.userRepo
}

func (t *transactor) EventRepo() event.Repo {
	t.eventRepo = event.New(event.Params{DB: t.tx})
	return t.eventRepo
}
