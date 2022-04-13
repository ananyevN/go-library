package rabbit

import (
	mb "github.com/bxcodec/library/message_broker"
	"github.com/streadway/amqp"
)

const RabbitMqUrl = "amqp://guest:guest@host.docker.internal:5672/"

type rabbitMqService struct {
	queue string
}

func NewRabbitMqService(q string) mb.MessageBroker {
	return &rabbitMqService{queue: q}
}

func (r rabbitMqService) Send(event mb.Event) error {
	conn, err := amqp.Dial(RabbitMqUrl)
	mb.FailOnError(err, FailedToConnect)
	defer conn.Close()

	ch, err := conn.Channel()
	mb.FailOnError(err, FailedToOpenChannel)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		r.queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	mb.FailOnError(err, FailedToOpenQueue)

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        event.Marshal(),
		})
	mb.FailOnError(err, FailedToPublishMessage)

	return nil
}

func (r rabbitMqService) Receive(emailChan chan []byte) error {
	conn, err := amqp.Dial(RabbitMqUrl)
	mb.FailOnError(err, FailedToConnect)
	defer conn.Close()

	ch, err := conn.Channel()
	mb.FailOnError(err, FailedToOpenChannel)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		r.queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	mb.FailOnError(err, FailedToOpenQueue)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	mb.FailOnError(err, FailedToRegisterConsumer)

	for d := range msgs {
		emailChan <- d.Body
	}

	return nil
}
