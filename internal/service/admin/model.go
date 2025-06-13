package admin

import "time"

const CtxKey = "admin"

type Admin struct {
	ID        int
	CountryID int
	IP        string
	Username  string
	FullName  string
}

type Audit struct {
	AdminId   int
	EventName string
	IpAddress string
	CreatedAt time.Time
	OldData   any
	NewData   any
}

const (
	CreateNotificationsEvent = "create_notifications_event"
	UpdateNotificationsEvent = "update_notifications_event"
	DeleteNotificationsEvent = "delete_notifications_event"
	LoadUsersEvent           = "load_users_event"
	UploadImageEvent         = "upload_image_event"
	RemoveImageEvent         = "remove_image_event"
	RunEvent                 = "run_event"
)
