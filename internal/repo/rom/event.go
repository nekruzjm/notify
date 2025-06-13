package rom

import (
	"context"

	"notifications/internal/db"
	"notifications/internal/lib/ctxman"
)

func (r *repo) DeleteEvent(ctx context.Context, id int) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Rom,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `DELETE FROM notification_events WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}
