package queue

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewProducerChannel(conn *amqp.Connection, producerExchangeName string) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %v", err)
	}
	err = ch.ExchangeDeclare(
		producerExchangeName,
		"direct",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to make queue decalre: %v", err)
	}
	//_, err = ch.QueueDeclare(
	//	config.ProducerExchangeName,
	//	false,
	//	true,
	//	false,
	//	false,
	//	nil,
	//)
	return ch, err
}

func NewConsumeChannel(conn *amqp.Connection, config *Config) (<-chan amqp.Delivery, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to make channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		config.ConsumerExchangeName,
		"direct",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to make exchange declare: %v", err)
	}
	q, err := ch.QueueDeclare(
		config.ConsumerExchangeName,
		false,
		true,
		false,
		false,
		nil,
	)
	for _, key := range config.RoutingKeys {
		err = ch.QueueBind(
			q.Name,
			key,
			config.ConsumerExchangeName,
			false,
			nil)
		if err != nil {
			return nil, fmt.Errorf("failed to bind queue to routing key %s: %v", key, err)
		}
	}
	msgs, err := ch.Consume(
		q.Name,
		"server",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start consume: %v", err)
	}

	return msgs, nil
}
