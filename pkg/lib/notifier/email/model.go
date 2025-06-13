package email

import "net/smtp"

type Email struct {
	plainAuth smtp.Auth
	host      string
	port      string
	sender    string
	receivers []string
	body      []byte
}

func New() *Email {
	return &Email{}
}

func (e *Email) SetHost(host string) *Email {
	e.host = host
	return e
}

func (e *Email) SetPort(port string) *Email {
	e.port = port
	return e
}

func (e *Email) SetSender(sender string) *Email {
	e.sender = sender
	return e
}

func (e *Email) SetReceiver(receivers []string) *Email {
	e.receivers = receivers
	return e
}

func (e *Email) SetPlainAuth(auth smtp.Auth) *Email {
	e.plainAuth = auth
	return e
}

func (e *Email) SetBody(body []byte) *Email {
	e.body = body
	return e
}

func (e *Email) Addr() string {
	return e.host + ":" + e.port
}
