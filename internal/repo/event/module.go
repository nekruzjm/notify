package event

import (
	"context"

	"go.uber.org/fx"

	"notifications/internal/db"
	"notifications/internal/lib/language"
)

var Module = fx.Provide(New)

type Repo interface {
	writer
	reader
}

type writer interface {
	Create(context.Context, *Event) (*Event, error)
	Update(context.Context, *Event) (*Event, error)
	UpdateImage(ctx context.Context, id int, image language.Language) (*Event, error)
	UpdateStatus(ctx context.Context, id int, status string) error
	UpdateExtraData(ctx context.Context, id int, extraData map[string]string) error
	Delete(ctx context.Context, id int) error
}

type reader interface {
	GetByID(ctx context.Context, id int) (*Event, error)
	GetActiveByIDWithLock(ctx context.Context, eventID int) (*Event, error)
	GetByFilter(context.Context, Filter) ([]*Event, error)
	GetAllActive(ctx context.Context) ([]*Event, error)
	GetSent(ctx context.Context) ([]*Event, error)
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
