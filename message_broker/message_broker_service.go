package message_broker

import (
	"log"
)

type EventType string

const (
	GetById EventType = "get.by.id.sql"
	UPDATE  EventType = "update.sql"
	DELETE  EventType = "delete.sql"
	ADD     EventType = "add.sql"
	FETCH   EventType = "fetch.sql"
)

type Event struct {
	Content string
	Subject string
}

type MessageBroker interface {
	Send(eventType EventType, content string) error
	Receive(eventType EventType, emailChan chan []byte) error
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
