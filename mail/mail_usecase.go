package mail

import (
	"fmt"
	"github.com/bxcodec/library/message_broker"
	"github.com/spf13/viper"
	"net/smtp"
	"os"
)

type emailUseCase struct {
	email Email
}

func init() {
	viper.SetConfigFile(`config.json`)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func NewSender() Sender {
	return &emailUseCase{email: Email{
		from:     viper.GetString(`mail.from`),
		password: os.Getenv("MAIL_PASS"),
		toEmail:  []string{viper.GetString(`mail.to`)},
		host:     viper.GetString(`mail.host`),
		port:     viper.GetString(`mail.port`),
	}}
}

func (e emailUseCase) SendEmail(event message_broker.Event) error {
	address := e.email.host + ":" + e.email.port
	auth := smtp.PlainAuth("", e.email.from, e.email.password, e.email.host)
	email := email{event}.compress()
	err := smtp.SendMail(address, auth, e.email.from, e.email.toEmail, email)
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
