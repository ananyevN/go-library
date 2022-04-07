package rabbit

import (
	"github.com/bxcodec/library/message_broker"
	"github.com/streadway/amqp"
	"log"
)

const RabbitMqUrl = "amqp://guest:guest@host.docker.internal:5672/"

type rabbitMqService struct {
	QueueName string
}

func NewRabbitMqService(QueueName string) message_broker.MessageBroker {
	return &rabbitMqService{QueueName: QueueName}
}

func (r rabbitMqService) Send(content string) error {
	conn, err := amqp.Dial(RabbitMqUrl)
	message_broker.FailOnError(err, FailedToConnect)
	defer conn.Close()

	ch, err := conn.Channel()
	message_broker.FailOnError(err, FailedToOpenChannel)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		r.QueueName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	message_broker.FailOnError(err, FailedToOpenQueue)

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(content),
		})
	message_broker.FailOnError(err, FailedToPublishMessage)

	return nil
}

func (r rabbitMqService) Receive() ([]message_broker.Event, error) {
	conn, err := amqp.Dial(RabbitMqUrl)
	message_broker.FailOnError(err, FailedToConnect)
	defer conn.Close()

	ch, err := conn.Channel()
	message_broker.FailOnError(err, FailedToOpenChannel)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		r.QueueName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	message_broker.FailOnError(err, FailedToOpenQueue)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	message_broker.FailOnError(err, FailedToRegisterConsumer)

	events := make([]message_broker.Event, 0)
	go func() {
		for d := range msgs {
			events = append(events, message_broker.Event{Content: string(d.Body)})
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	return events, nil
}
