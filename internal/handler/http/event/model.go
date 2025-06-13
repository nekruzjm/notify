package event

import (
	"time"

	"notifications/internal/lib/language"
)

// Route keys
const (
	_id       = "id"
	_status   = "status"
	_topic    = "topic"
	_type     = "type"
	_limit    = "limit"
	_offset   = "offset"
	_image    = "image"
	_users    = "users"
	_language = "language"
	_userID   = "userID"
)

type request struct {
	Status      string            `json:"status"`
	Topic       string            `json:"topic"`
	Category    string            `json:"category"`
	Link        string            `json:"link"`
	ScheduledAt string            `json:"scheduledAt"`
	Title       language.Language `json:"title"`
	Body        language.Language `json:"body"`
	ExtraData   map[string]string `json:"extraData"`
}

var _ eventModel

type eventModel struct {
	ID          int               `json:"id"`
	Topic       string            `json:"topic"`
	Status      string            `json:"status"`
	Title       language.Language `json:"title"`
	Body        language.Language `json:"body"`
	Image       language.Language `json:"image"`
	Category    string            `json:"category"`
	Link        string            `json:"link"`
	ExtraData   map[string]string `json:"extraData"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	ScheduledAt time.Time         `json:"scheduledAt"`
}
