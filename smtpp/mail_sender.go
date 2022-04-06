package smtpp

import (
	"fmt"
	"github.com/bxcodec/library/message_brocker"
	"net/smtp"
	"os"
)

func SendEmail(event message_brocker.Event) {
	from := "n.ananyev777@gmail.com"
	password := os.Getenv("MAIL_PASS")

	toEmail := []string{"life_love_asap@mail.ru"}
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	subject := "Subject: SMTP_GENERATED_MESSAGE\n"
	body := event.Content
	message := []byte(subject + body)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, toEmail, message)
	if err != nil {
		fmt.Println("ERROR")
	}
}
