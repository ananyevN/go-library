package rabbit

import (
	mb "github.com/bxcodec/library/message_broker"
	"github.com/streadway/amqp"
	"log"
)

const RabbitMqUrl = "amqp://guest:guest@localhost:5672/"

type rabbitMqService struct {
	exchange string
}

func NewRabbitMqService(ex string) mb.MessageBroker {
	return &rabbitMqService{exchange: ex}
}

func (r rabbitMqService) Send(eventType mb.EventType, content string) error {
	conn, err := amqp.Dial(RabbitMqUrl)
	mb.FailOnError(err, FailedToConnect)
	defer conn.Close()

	ch, err := conn.Channel()
	mb.FailOnError(err, FailedToOpenChannel)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		r.exchange, // name
		"topic",    // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	mb.FailOnError(err, "Failed to declare an exchange")

	log.Printf("Publishing to %s topic", string(eventType))

	err = ch.Publish(
		r.exchange,        // exchange
		string(eventType), // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(content),
		})
	mb.FailOnError(err, FailedToPublishMessage)

	log.Printf(" [x] Sent %s", content)

	return nil
}

func (r rabbitMqService) Receive(eventType mb.EventType, emailChan chan []byte) error {
	conn, err := amqp.Dial(RabbitMqUrl)
	mb.FailOnError(err, FailedToConnect)
	defer conn.Close()

	ch, err := conn.Channel()
	mb.FailOnError(err, FailedToOpenChannel)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		r.exchange, // name
		"topic",    // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	mb.FailOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	mb.FailOnError(err, FailedToOpenQueue)

	log.Printf("Binding queue %s to exchange %s with routing key %s",
		q.Name, "logs_topic", string(eventType))

	err = ch.QueueBind(
		q.Name,            // queue name
		string(eventType), // routing key
		r.exchange,        // exchange
		false,
		nil)
	mb.FailOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	mb.FailOnError(err, FailedToRegisterConsumer)

	for d := range msgs {
		emailChan <- d.Body
	}

	return nil
}
