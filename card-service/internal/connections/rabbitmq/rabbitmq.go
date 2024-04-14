package rabbitmq

import (
	"github.com/sm888sm/halten-backend/card-service/internal/config"
	"github.com/streadway/amqp"
)

var (
	RabbitMQChannel *amqp.Channel
	conn            *amqp.Connection
)

func Connect(cfg *config.RabbitMQConfig) error {
	var err error
	conn, err = amqp.Dial(cfg.URL)
	if err != nil {
		return err
	}

	RabbitMQChannel, err = conn.Channel()
	if err != nil {
		return err
	}

	err = RabbitMQChannel.ExchangeDeclare(
		"halten", // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	return nil
}

func GetConnection() *amqp.Connection {
	return conn
}
