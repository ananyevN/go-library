package rabbit

import (
	"fmt"
	"github.com/bxcodec/library/mail"
	"github.com/bxcodec/library/message_broker"
	"github.com/streadway/amqp"
	"log"
)

const RabbitMqUrl = "amqp://guest:guest@localhost:5672/"

type rabbitMqService struct {
}

func NewRabbitMqService() message_broker.MessageBroker {
	return &rabbitMqService{}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (r rabbitMqService) Send(eventType message_broker.EventType, content string) error {
	conn, err := amqp.Dial(RabbitMqUrl)
	message_broker.FailOnError(err, FailedToConnect)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"crud_exchange", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	log.Printf("Publishing to %s topic", string(eventType))

	err = ch.Publish(
		"crud_exchange",   // exchange
		string(eventType), // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(content),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", content)

	return nil
}

func (r rabbitMqService) Receive(eventType message_broker.EventType) (chan string, error) {
	conn, err := amqp.Dial(RabbitMqUrl)
	message_broker.FailOnError(err, FailedToConnect)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"crud_exchange", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	log.Printf("Binding queue %s to exchange %s with routing key %s",
		q.Name, "logs_topic", string(eventType))

	err = ch.QueueBind(
		q.Name,            // queue name
		string(eventType), // routing key
		"crud_exchange",   // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	func() {
		for d := range msgs {
			go func(del amqp.Delivery) {
				body := fmt.Sprintf("%s", del.Body)
				event := message_broker.Event{Content: body}
				event.Subject = string(eventType)
				useCase := mail.NewMailUseCase()
				useCase.SendEmail(event)
			}(d)
		}
	}()

	<-forever

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	//close(emailChan)
	return nil, nil
}
