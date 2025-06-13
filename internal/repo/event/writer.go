package event

import (
	"context"

	"notifications/internal/db"
	"notifications/internal/lib/ctxman"
	"notifications/internal/lib/language"
)

func (r *repo) Create(ctx context.Context, event *Event) (*Event, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	var e = new(Event)
	err := r.db.QueryRow(ctx, `
				INSERT INTO events (id, topic, status, title, body, image, category, link, extra_data, scheduled_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING `+_cols,
		event.ID,
		event.Topic,
		event.Status,
		event.Title,
		event.Body,
		event.Image,
		event.Category,
		event.Link,
		event.ExtraData,
		event.ScheduledAt).Scan(fields(e)...)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *repo) Update(ctx context.Context, event *Event) (*Event, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	var e = new(Event)
	err := r.db.QueryRow(ctx, `
			UPDATE events SET 
				status = $1,
				title = $2, 
				body = $3,
				category = $4,
				link = $5,
				extra_data = $6,
				scheduled_at = $7,
				updated_at = now()
			WHERE id = $8 RETURNING `+_cols,
		event.Status,
		event.Title,
		event.Body,
		event.Category,
		event.Link,
		event.ExtraData,
		event.ScheduledAt,
		event.ID).Scan(fields(e)...)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *repo) UpdateStatus(ctx context.Context, id int, status string) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `UPDATE events SET status = $1, updated_at = now() WHERE id = $2`, status, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) UpdateExtraData(ctx context.Context, id int, extraData map[string]string) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `UPDATE events SET extra_data = extra_data || $1, updated_at = now() WHERE id = $2`, extraData, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) UpdateImage(ctx context.Context, id int, image language.Language) (*Event, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	var event = new(Event)
	err := r.db.QueryRow(ctx, `UPDATE events SET image = $1, updated_at = now() WHERE id = $2 RETURNING `+_cols, image, id).Scan(fields(event)...)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (r *repo) Delete(ctx context.Context, id int) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
