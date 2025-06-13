package event

import (
	"sync"
	"time"

	"firebase.google.com/go/v4/messaging"

	"notifications/internal/lib/language"
	"notifications/internal/repo/event"
	"notifications/pkg/lib/notifier/firebase"
)

const (
	_underscoreDelim = "_"
	_dotDelim        = "."
	_empty           = ""
	_topicRegex      = "[a-zA-Z0-9-_.~%]+"
)

const (
	_active        = "active"
	_sent          = "sent"
	_failed        = "failed"
	_failedLoading = "failed_loading"
	_loadingUsers  = "loading_users"
	_draft         = "draft"
)

const (
	_successCountKey = "successCount"
	_failedCountKey  = "failedCount"
	_failedReasonKey = "failedReason"
)

const _userIDCsvHeader = "userID"

const (
	_title              = "title"
	_comment            = "comment"
	_message            = "message"
	_category           = "category"
	_image              = "image"
	_button             = "button"
	_sectionName        = "sectionName"
	_serviceID          = "serviceID"
	_cashBackID         = "cashBackID"
	_cashBackCategoryID = "cashBackCategoryID"
	_badge              = "badge"
)

const _topicSubCacheKey = ":topic-subscription:"
const _fcmTokenErr = "NOT_FOUND"

type Message struct {
	UserID int
	Data   map[string]string
}

type ChunkResult struct {
	SuccessCount int
	FailedCount  int
	ErrCount     int
	Err          error
}

type PrepareUsersResponse struct {
	SuccessCount int64
	FailedCount  int64
	ErrCount     int64
}

type FailedResponse struct {
	UserID int    `json:"user,omitempty"`
	Token  string `json:"token,omitempty"`
	Reason string `json:"reason,omitempty"`
}

type Filter struct {
	Topic  string
	Status string
	ID     uint
	Limit  uint
	Offset uint
}

type Event struct {
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

	SubscribeAll bool `json:"subscribeAll,omitempty"`
}

type Request struct {
	ID              int
	Topic           string
	Status          string
	Title           language.Language
	Body            language.Language
	Category        string
	Link            string
	ScheduledAt     string
	ScheduledAtTime time.Time
	ExtraData       map[string]string
}

func setupMessages(event *event.Event) []*messaging.Message {
	var (
		languages = language.GetAll()
		msgCh     = make(chan *messaging.Message, len(languages))
		wg        sync.WaitGroup
	)

	for _, lang := range languages {
		wg.Add(1)
		go func() {
			defer wg.Done()

			data := make(map[string]string)
			data[_title] = event.Title.Get(lang)
			data[_comment] = event.Title.Get(lang)
			data[_message] = event.Body.Get(lang)
			data[_image] = event.Image.Get(lang)
			data[_category] = event.Category
			data[_button] = event.Link
			data[_sectionName] = event.ExtraData[_sectionName]
			data[_serviceID] = event.ExtraData[_serviceID]
			data[_cashBackID] = event.ExtraData[_cashBackID]
			data[_cashBackCategoryID] = event.ExtraData[_cashBackCategoryID]
			data[_badge] = event.ExtraData[_badge]

			message := new(messaging.Message)
			message.Data = data
			message.Topic = buildTopic(event.Topic, lang)
			firebase.AndroidMSG(message, data, firebase.AndroidNormalPriority)
			firebase.IosMSG(message, data, firebase.ApnsNormalPriority)

			msgCh <- message
		}()
	}

	go func() {
		wg.Wait()
		close(msgCh)
	}()

	var messages = make([]*messaging.Message, 0, len(languages))
	for msg := range msgCh {
		messages = append(messages, msg)
	}

	return messages
}

func (e *Event) toService(event *event.Event) {
	e.ID = event.ID
	e.Topic = event.Topic
	e.Status = event.Status
	e.Title = event.Title
	e.Body = event.Body
	e.Image = event.Image
	e.Category = event.Category
	e.Link = event.Link
	e.ExtraData = event.ExtraData
	e.ScheduledAt = event.ScheduledAt
	e.CreatedAt = event.CreatedAt
	e.UpdatedAt = event.UpdatedAt
}

func toRepo(e *Event) *event.Event {
	return &event.Event{
		ID:          e.ID,
		Topic:       e.Topic,
		Status:      e.Status,
		Title:       e.Title,
		Body:        e.Body,
		Image:       e.Image,
		Category:    e.Category,
		ExtraData:   e.ExtraData,
		Link:        e.Link,
		ScheduledAt: e.ScheduledAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
