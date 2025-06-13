package user

import "time"

type User struct {
	UserID            int
	Phone             string
	PersonExternalRef string
	Token             string
	Status            string
	Language          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	CountryID         int8
	PushEnabled       bool
}

type EventRelation struct {
	UserID int
	Topic  string
	Lang   string
}

const _cols = `
		user_id, 
		phone, 
		person_external_ref,
		token, 
		status,
		push_enabled,
		language,
		country_id,
		created_at, 
		updated_at`

func fields(u *User) []any {
	return []any{
		&u.UserID,
		&u.Phone,
		&u.PersonExternalRef,
		&u.Token,
		&u.Status,
		&u.PushEnabled,
		&u.Language,
		&u.CountryID,
		&u.CreatedAt,
		&u.UpdatedAt,
	}
}
