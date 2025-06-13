package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"notifications/internal/db"
	"notifications/internal/lib/ctxman"
	"notifications/internal/repo/repomodel"
)

func (r *repo) GetTopicsByUserID(ctx context.Context, userID int) ([]EventRelation, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	rows, err := r.db.Query(ctx, `SELECT e.topic, uer.language FROM user_event_relations uer LEFT JOIN events e ON uer.event_id = e.id WHERE uer.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eventRel = make([]EventRelation, 0)

	for rows.Next() {
		var e EventRelation
		err = rows.Scan(&e.Topic, &e.Lang)
		if err != nil {
			return nil, err
		}
		eventRel = append(eventRel, e)
	}

	if len(eventRel) == 0 {
		return nil, repomodel.ErrNotFound
	}

	return eventRel, nil
}

func (r *repo) InsertUserEventRelation(ctx context.Context, userID, eventID int, lang string) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `INSERT INTO user_event_relations (user_id, event_id, language) VALUES ($1, $2, $3)`, userID, eventID, lang)
	if err != nil {
		var pgErr = new(pgconn.PgError)
		errors.As(err, &pgErr)
		if pgerrcode.UniqueViolation == pgErr.Code {
			return repomodel.ErrUniqueViolation
		}
		return err
	}

	return nil
}

func (r *repo) DeleteRelationByEventID(ctx context.Context, eventID int) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `DELETE FROM user_event_relations WHERE event_id = $1`, eventID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) DeleteRelationByUserID(ctx context.Context, userID int) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `DELETE FROM user_event_relations WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) DeleteRelation(ctx context.Context, userID, eventID int) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	res, err := r.db.Exec(ctx, `DELETE FROM user_event_relations WHERE user_id = $1 AND event_id = $2`, userID, eventID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return repomodel.ErrNotFound
	}

	return nil
}

func (r *repo) UpdateRelationLanguage(ctx context.Context, userID int, lang string) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `UPDATE user_event_relations SET language = $1 WHERE user_id = $2 `, lang, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) BatchInsert(ctx context.Context, eventID int, relations []EventRelation) (int64, error) {
	const maxRetries = 3
	return r.batchInsertWithRetries(ctx, eventID, relations, 0, maxRetries)
}

func (r *repo) batchInsertWithRetries(ctx context.Context, eventID int, relations []EventRelation, attempt, maxRetries int) (int64, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	const (
		_tableName   = "user_event_relations"
		_eventIDCol  = "event_id"
		_userIDCol   = "user_id"
		_languageCol = "language"
	)

	rowsCount, err := r.db.CopyFrom(
		ctx,
		pgx.Identifier{_tableName},
		[]string{_eventIDCol, _userIDCol, _languageCol},
		pgx.CopyFromSlice(len(relations), func(i int) ([]any, error) {
			return []any{eventID, relations[i].UserID, relations[i].Lang}, nil
		}))
	if err != nil {
		var pgErr = new(pgconn.PgError)
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			if attempt >= maxRetries {
				return 0, fmt.Errorf("max retries reached: %w", err)
			}

			_, err = r.db.Exec(ctx, `DELETE FROM user_event_relations WHERE event_id = $1`, eventID)
			if err != nil {
				return 0, err
			}
			return r.batchInsertWithRetries(ctx, eventID, relations, attempt+1, maxRetries)
		}
		return 0, err
	}

	return rowsCount, nil
}
