package mail

import (
	"fmt"
	"github.com/bxcodec/library/message_broker"
	"net/smtp"
	"os"
)

type emailUseCase struct {
	from     string
	password string
	toEmail  []string
	host     string
	port     string
}

func NewMailUseCase() *emailUseCase {
	return &emailUseCase{
		from:     "n.ananyev777@gmail.com",
		password: os.Getenv("MAIL_PASS"),
		toEmail:  []string{"life_love_asap@mail.ru"},
		host:     "smtp.gmail.com",
		port:     "587",
	}
}

func (e emailUseCase) SendEmail(event message_broker.Event) (err error) {
	address := e.host + ":" + e.port
	auth := smtp.PlainAuth("", e.from, e.password, e.host)
	email := email{event}.compress()
	err = smtp.SendMail(address, auth, e.from, e.toEmail, email)
	if err != nil {
		err = fmt.Errorf("Error while sending email with %s subject. ", event.Subject)
	}
	return err
}

type email struct {
	event message_broker.Event
}

func (email email) compress() []byte {
	emailSubject := fmt.Sprintf("Subject: %s\r\n\r\n", email.event.Subject)
	emailContent := fmt.Sprintf(" %s\r\n", email.event.Content)
	return []byte(emailSubject + emailContent)
}
