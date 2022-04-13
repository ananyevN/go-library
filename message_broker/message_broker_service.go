package message_broker

import (
	"encoding/json"
	"log"
)

type EventType string

type MessageBroker interface {
	Send(event Event) error
	Receive(emailChan chan []byte) error
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

const (
	GetById EventType = "get.by.id.sql"
	UPDATE  EventType = "update.sql"
	DELETE  EventType = "delete.sql"
	ADD     EventType = "add.sql"
	FETCH   EventType = "fetch.sql"
)

type Event struct {
	Content string `json:"content"`
	Subject string `json:"subject"`
}

func (e *Event) Marshal() []byte {
	res, err := json.Marshal(e)
	if err != nil {
		log.Fatalf("Error marshal to json")
	}
	return res
}

func (e *Event) Unmarshal(body []byte) *Event {
	err := json.Unmarshal(body, e)
	if err != nil {
		return &Event{}
	}
	return e
}
