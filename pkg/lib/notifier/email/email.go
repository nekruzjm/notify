package email

import "net/smtp"

func (e *Email) Send() error {
	return smtp.SendMail(e.Addr(), e.plainAuth, e.sender, e.receivers, e.body)
}
