package mail

import (
	"github.com/bxcodec/library/message_broker"
)

type Email struct {
	from     string
	password string
	toEmail  []string
	host     string
	port     string
}

type Sender interface {
	SendEmail(event message_broker.Event) error
}
