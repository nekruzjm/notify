package rom

import (
	"notifications/internal/lib/language"
)

type Event struct {
	ID        int
	Title     language.Language
	Body      language.Language
	Image     language.Language
	ExtraData map[string]string
}

type Inbox struct {
	ID        int
	UserID    int
	CountryID int8
	Type      string
	Title     language.Language
	Body      language.Language
	ExtraData map[string]string
}
