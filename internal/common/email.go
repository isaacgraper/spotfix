package common

import (
	"log"
	"net/smtp"
)

type Email struct {
	From string
	Pwd  string
	To   []string

	SmtpHost string
	SmtpPort string

	Message []byte
}

func NewEmail(from, pwd string, to []string, smtpHost, smtpPort string, message []byte) *Email {
	return &Email{
		From:     from,
		Pwd:      pwd,
		To:       to,
		SmtpHost: smtpHost,
		SmtpPort: smtpPort,
		Message:  message,
	}
}

func SendEmail(e *Email) {
	auth := smtp.PlainAuth("", e.From, e.Pwd, e.SmtpHost)

	err := smtp.SendMail(e.SmtpHost+":"+e.SmtpPort, auth, e.From, e.To, e.Message)
	if err != nil {
		log.Println("error sending email: ", err)
		return
	}

	log.Println("report send to: ", e.To)
}
