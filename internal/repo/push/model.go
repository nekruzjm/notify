package push

import (
	"time"

	"notifications/internal/lib/language"
)

type Push struct {
	ID        int
	UserID    int
	Status    string
	Type      string
	APIClient string
	Title     language.Language
	Body      language.Language
	CreatedAt time.Time
	UpdatedAt time.Time
}

const _cols = `
			id, 
			user_id, 
			type,
			status,
			title,
			body,
			api_client,
			created_at, 
			updated_at`

func fields(p *Push) []any {
	return []any{
		&p.ID,
		&p.UserID,
		&p.Type,
		&p.Status,
		&p.Title,
		&p.Body,
		&p.APIClient,
		&p.CreatedAt,
		&p.UpdatedAt,
	}
}
