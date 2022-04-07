package mail

import "github.com/bxcodec/library/message_broker"

type MailService interface {
	SendEmail(event message_broker.Event) error
}
