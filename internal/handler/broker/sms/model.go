package sms

import "time"

type Request struct {
	Phone         string    `json:"phone"`
	Body          string    `json:"body"`
	SenderAddress string    `json:"senderAddress"`
	Priority      int       `json:"priority"`
	Type          int       `json:"type"`
	ExpiresIn     int       `json:"expiresIn"`
	ScheduledAt   time.Time `json:"scheduledAt"`
}
