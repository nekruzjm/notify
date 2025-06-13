package rom

import (
	"context"

	"notifications/internal/db"
	"notifications/internal/lib/ctxman"
)

func (r *repo) InsertInbox(ctx context.Context, inbox *Inbox) error {
	ctx = ctxman.Save(ctx, ctxman.Info{DBName: db.Rom, IsReplica: false})
	_, err := r.db.Exec(ctx, `INSERT INTO notification_inbox (id, user_id, country_id, type, title, body, extra_data) 
									VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		inbox.ID, inbox.UserID, inbox.CountryID, inbox.Type, inbox.Title, inbox.Body, inbox.ExtraData)
	if err != nil {
		return err
	}

	return nil
}
