package rabbit

import (
	"fmt"
	mb "github.com/bxcodec/library/message_broker"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
)

var rabbitMqUrl string

func init() {
	viper.SetConfigFile(`config.json`)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	rabbitHost := viper.GetString(`rabbit.host`)
	rabbitPort := viper.GetInt(`rabbit.port`)
	rabbitUser := viper.GetString(`rabbit.user`)
	rabbitPass := viper.GetString(`rabbit.pass`)
	rabbitMqUrl = fmt.Sprintf("amqp://%s:%s@%s:%d/", rabbitUser, rabbitPass, rabbitHost, rabbitPort)
}

type rabbitMqService struct {
	queue string
}

func NewRabbitMqService(q string) mb.MessageBroker {
	return &rabbitMqService{queue: q}
}

func (r rabbitMqService) Send(event mb.Event) error {
	conn, err := amqp.Dial(rabbitMqUrl)
	if err != nil {
		log.Printf("%s: %s", FailedToConnect, err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("%s: %s", FailedToOpenChannel, err)
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		r.queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Printf("%s: %s", FailedToOpenQueue, err)
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        event.Marshal(),
		})
	if err != nil {
		log.Printf("%s: %s", FailedToPublishMessage, err)
		return err
	}

	return nil
}

func (r rabbitMqService) Receive(emailChan chan []byte) error {
	conn, err := amqp.Dial(rabbitMqUrl)
	if err != nil {
		log.Printf("%s: %s", FailedToConnect, err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("%s: %s", FailedToOpenChannel, err)
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		r.queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Printf("%s: %s", FailedToOpenQueue, err)
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Printf("%s: %s", FailedToRegisterConsumer, err)
		return err
	}

	for d := range msgs {
		emailChan <- d.Body
	}

	return nil
}
