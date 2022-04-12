// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocks

import (
	"github.com/bxcodec/library/message_broker"
	mock "github.com/stretchr/testify/mock"
)

// MessageBroker is an autogenerated mock type for the MessageBroker type
type MessageBroker struct {
	mock.Mock
}

// Receive provides a mock function with given fields: eventType, emailChan
func (_m *MessageBroker) Receive(eventType message_broker.EventType, emailChan chan []byte) error {
	ret := _m.Called(eventType, emailChan)

	var r0 error
	if rf, ok := ret.Get(0).(func(message_broker.EventType, chan []byte) error); ok {
		r0 = rf(eventType, emailChan)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Send provides a mock function with given fields: eventType, content
func (_m *MessageBroker) Send(eventType message_broker.EventType, content string) error {
	ret := _m.Called(eventType, content)

	var r0 error
	if rf, ok := ret.Get(0).(func(message_broker.EventType, string) error); ok {
		r0 = rf(eventType, content)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
