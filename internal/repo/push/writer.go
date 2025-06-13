package push

import (
	"context"

	"notifications/internal/db"
	"notifications/internal/lib/ctxman"
)

func (r *repo) Insert(ctx context.Context, push *Push) (*Push, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{DBName: db.Notifications, IsReplica: false})
	createdPush := new(Push)

	err := r.db.QueryRow(ctx, `
				INSERT INTO push (id, user_id, type, status, title, body, api_client) 
				VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING `+_cols,
		push.ID,
		push.UserID,
		push.Type,
		push.Status,
		push.Title,
		push.Body,
		push.APIClient).Scan(fields(createdPush)...)
	if err != nil {
		return nil, err
	}

	return push, nil
}

func (r *repo) DeleteByIDs(ctx context.Context, ids []int) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, "DELETE FROM push WHERE id = ANY($1)", ids)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) Clean(ctx context.Context) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, "DELETE FROM push WHERE status = 'approved' AND created_at < now() - INTERVAL '10 days'")
	if err != nil {
		return err
	}

	return nil
}
