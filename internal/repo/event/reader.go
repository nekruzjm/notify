package event

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"

	"notifications/internal/db"
	"notifications/internal/lib/ctxman"
	"notifications/internal/repo/repomodel"
	"notifications/pkg/util/strset"
)

func (r *repo) GetByID(ctx context.Context, eventID int) (*Event, error) {
	return r.selectEvent(ctx, "id = $1 AND (status = 'active' OR status = 'draft')", eventID)
}

func (r *repo) GetActiveByIDWithLock(ctx context.Context, eventID int) (*Event, error) {
	return r.selectEvent(ctx, "id = $1 AND (status = 'active' OR status = 'draft') FOR UPDATE", eventID)
}

func (r *repo) GetAllActive(ctx context.Context) ([]*Event, error) {
	return r.selectEvents(ctx, "status = 'active'")
}

func (r *repo) GetSent(ctx context.Context) ([]*Event, error) {
	return r.selectEvents(ctx, "status = 'sent'")
}

func (r *repo) GetByFilter(ctx context.Context, filter Filter) ([]*Event, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	var (
		conditions = "1=1"
		args       []any
		idx        uint8
	)

	if filter.Limit == 0 {
		filter.Limit = 10
	}
	if filter.ID != 0 {
		idx++
		conditions += " AND id = $" + strset.IntToStr(int(idx))
		args = append(args, filter.ID)
	}
	if !strset.IsEmpty(filter.Topic) {
		idx++
		conditions += " AND topic = $" + strset.IntToStr(int(idx))
		args = append(args, filter.Topic)
	}
	if !strset.IsEmpty(filter.Status) {
		idx++
		conditions += " AND status = $" + strset.IntToStr(int(idx))
		args = append(args, filter.Status)
	}

	conditions += " ORDER BY updated_at DESC LIMIT $" + strset.IntToStr(int(idx+1)) + " OFFSET $" + strset.IntToStr(int(idx+2))
	args = append(args, filter.Limit, filter.Offset)

	var builder strings.Builder
	builder.WriteString("SELECT ")
	builder.WriteString(_cols)
	builder.WriteString(" FROM events WHERE ")
	builder.WriteString(conditions)
	query := builder.String()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events = make([]*Event, 0, filter.Limit)

	for rows.Next() {
		var event = new(Event)
		err = rows.Scan(fields(event)...)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if len(events) == 0 {
		return nil, repomodel.ErrNotFound
	}

	return events, nil
}

func (r *repo) selectEvent(ctx context.Context, condition string, args ...any) (*Event, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	var (
		event   = new(Event)
		builder strings.Builder
	)

	builder.WriteString("SELECT ")
	builder.WriteString(_cols)
	builder.WriteString(" FROM events WHERE ")
	builder.WriteString(condition)

	query := builder.String()

	err := r.db.QueryRow(ctx, query, args...).Scan(fields(event)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repomodel.ErrNotFound
		}
		return nil, err
	}

	return event, nil
}

func (r *repo) selectEvents(ctx context.Context, condition string, args ...any) (events []*Event, err error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	var builder strings.Builder
	builder.WriteString("SELECT ")
	builder.WriteString(_cols)
	builder.WriteString(" FROM events WHERE ")
	builder.WriteString(condition)

	query := builder.String()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event = new(Event)
		err = rows.Scan(fields(event)...)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if len(events) == 0 {
		return nil, repomodel.ErrNotFound
	}

	return events, nil
}
