package rabbit

import (
	"github.com/bxcodec/library/message_brocker"
	"github.com/bxcodec/library/smtpp"
	"github.com/streadway/amqp"
	"log"
)

const RabbitMqUrl = "amqp://guest:guest@localhost:5672/"

type rabbitMqService struct {
	QueueName string
}

func NewRabbitMqService(QueueName string) message_brocker.MessageBroker {
	return &rabbitMqService{QueueName: QueueName}
}

func (r rabbitMqService) Send(content string) error {
	conn, err := amqp.Dial(RabbitMqUrl)
	message_brocker.FailOnError(err, FailedToConnect)
	defer conn.Close()

	ch, err := conn.Channel()
	message_brocker.FailOnError(err, FailedToOpenChannel)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		r.QueueName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	message_brocker.FailOnError(err, FailedToOpenQueue)

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(content),
		})
	message_brocker.FailOnError(err, FailedToPublishMessage)

	smtpp.SendEmail(message_brocker.Event{Content: content})

	return nil
}

func (r rabbitMqService) Receive() (message_brocker.Event, error) {
	conn, err := amqp.Dial(RabbitMqUrl)
	message_brocker.FailOnError(err, FailedToConnect)
	defer conn.Close()

	ch, err := conn.Channel()
	message_brocker.FailOnError(err, FailedToOpenChannel)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		r.QueueName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	message_brocker.FailOnError(err, FailedToOpenQueue)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	message_brocker.FailOnError(err, FailedToRegisterConsumer)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	return message_brocker.Event{}, nil
}
