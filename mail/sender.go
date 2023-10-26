package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender interface {
	SendEmail(
		To []string,
		Cc []string,
		Subject string,
		Content string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name, fromEmailAddres, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddres,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(to []string, cc []string, subject string, content string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.To = to
	e.Cc = cc
	e.HTML = []byte(content)

	smtp_auth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, "smtp.gmail.com")
	return e.Send("smtp.gmail.com:587", smtp_auth)
}