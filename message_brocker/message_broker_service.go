package message_brocker

import "log"

type Event struct {
	Content string
}

type MessageBroker interface {
	Send(content string) error
	Receive() (Event, error)
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
