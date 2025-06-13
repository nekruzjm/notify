package rom

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"notifications/internal/db"
	"notifications/internal/lib/ctxman"
)

// BatchInsert in terms of performance is better than  Insert query, because it uses CopyFrom method.
// If there is a unique violation, it deletes all the rows with the same event_id and inserts them again, because CopyFrom does not support ON CONFLICT DO NOTHING.
// https://stackoverflow.com/questions/46715354/how-does-copy-work-and-why-is-it-so-much-faster-than-insert
func (r *repo) BatchInsert(ctx context.Context, eventID int, userIDs []int, event *Event) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Rom,
		IsReplica: false,
	})

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		} else {
			err = tx.Commit(ctx)
		}
	}()

	const (
		_tableName  = "notification_events_user_relation"
		_colEventID = "event_id"
		_colUserID  = "user_id"
	)

	_, err = tx.CopyFrom(ctx, pgx.Identifier{_tableName}, []string{_colEventID, _colUserID}, pgx.CopyFromSlice(len(userIDs), func(i int) ([]any, error) {
		return []any{eventID, userIDs[i]}, nil
	}))
	if err != nil {
		var pgErr = new(pgconn.PgError)
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			_, err = tx.Exec(ctx, `DELETE FROM notification_events_user_relation WHERE event_id = $1`, eventID)
			if err != nil {
				return err
			}
			_, err = tx.CopyFrom(ctx, pgx.Identifier{_tableName}, []string{_colEventID, _colUserID}, pgx.CopyFromSlice(len(userIDs), func(i int) ([]any, error) {
				return []any{eventID, userIDs[i]}, nil
			}))
			if err != nil {
				return err
			}
		}
		return err
	}

	_, err = tx.Exec(ctx, `
				INSERT INTO notification_events (id, title, body, image, extra_data) 
				VALUES ($1, $2, $3, $4, $5) 
				ON CONFLICT DO NOTHING`,
		event.ID, event.Title, event.Body, event.Image, event.ExtraData)
	if err != nil {
		return err
	}

	return nil
}
