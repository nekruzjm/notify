package sms

import "time"

type Request struct {
	Phone         string    `json:"phoneNumber"`
	Text          string    `json:"text"`
	SenderAddress string    `json:"senderAddress"`
	Priority      int       `json:"priority"`
	ExpiresIn     int       `json:"expiresIn"`
	SmsType       int       `json:"smsType"`
	ScheduledAt   time.Time `json:"scheduledAt"`
}
