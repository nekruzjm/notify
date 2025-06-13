package event

import (
	"time"

	"notifications/internal/lib/language"
)

type Event struct {
	ID          int
	Topic       string
	Status      string
	Title       language.Language
	Body        language.Language
	Image       language.Language
	Category    string
	Link        string
	ExtraData   map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ScheduledAt time.Time
}

type Filter struct {
	Topic  string
	Status string
	ID     uint
	Limit  uint
	Offset uint
}

const _cols = `
			id, 
			topic, 
			status,
			title,
			body,
			image,
			category,
			link,
			extra_data,
			created_at, 
			updated_at, 
			scheduled_at`

func fields(e *Event) []any {
	return []any{
		&e.ID,
		&e.Topic,
		&e.Status,
		&e.Title,
		&e.Body,
		&e.Image,
		&e.Category,
		&e.Link,
		&e.ExtraData,
		&e.CreatedAt,
		&e.UpdatedAt,
		&e.ScheduledAt,
	}
}
