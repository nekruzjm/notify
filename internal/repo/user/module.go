package user

import (
	"context"

	"go.uber.org/fx"

	"notifications/internal/db"
)

var Module = fx.Provide(New)

type Repo interface {
	writer
	reader
}

type writer interface {
	Create(context.Context, User) error
	UpdateToken(ctx context.Context, userID int, token string) error
	UpdateStatus(ctx context.Context, userID int, status string) error
	UpdatePhone(ctx context.Context, userID int, phone string) error
	UpdatePersonExternalRef(ctx context.Context, userID int, personExternalRef string) error
	UpdateUserSettings(ctx context.Context, userID int, language string, isEnabled bool) error
	DeleteByUserID(ctx context.Context, userID int) error

	InsertUserEventRelation(ctx context.Context, userID, eventID int, lang string) error
	DeleteRelationByEventID(ctx context.Context, eventID int) error
	DeleteRelationByUserID(ctx context.Context, userID int) error
	DeleteRelation(ctx context.Context, userID, eventID int) error
	BatchInsert(ctx context.Context, eventID int, relations []EventRelation) (int64, error)
	UpdateRelationLanguage(ctx context.Context, userID int, lang string) error
}

type reader interface {
	GetByUserID(ctx context.Context, userID int) (*User, error)
	GetActiveByPhone(ctx context.Context, phone string) (*User, error)
	GetActiveByPersonExternalRef(ctx context.Context, personExternalRef string) (*User, error)

	GetTokensByUserIDs(ctx context.Context, userIDs []string) ([]User, error)
	GetTopicsByUserID(ctx context.Context, userID int) ([]EventRelation, error)
	GetUserIDsByEventID(ctx context.Context, eventID int) ([]int, error)
	GetRelationsByEventID(ctx context.Context, eventID int) ([]EventRelation, error)
	GetTokensWithLimit(ctx context.Context, lastID int) ([]User, error)
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
