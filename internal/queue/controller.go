package queue

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	Addr     = "amqp://guest:guest@localhost:5672/"
	MafiaKey = "mafia"
	AllKey   = "all"
)

type Config struct {
	Addr                 string
	RoutingKeys          []string
	ProducerExchangeName []string
	ConsumerExchangeName string
}

type Controller struct {
	config *Config

	ChatChan <-chan string

	produceQueue   map[string]*amqp.Channel
	consumeChannel <-chan amqp.Delivery
	conn           *amqp.Connection
}

func NewController(config *Config, chatChan chan string) (*Controller, error) {
	ctl := &Controller{
		config:       config,
		ChatChan:     chatChan,
		produceQueue: map[string]*amqp.Channel{},
	}
	conn, err := amqp.Dial(config.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	ctl.conn = conn

	if err != nil {
		return nil, fmt.Errorf("failed to make new produce channel: %v", err)
	}
	ctl.consumeChannel, err = NewConsumeChannel(conn, config)
	if err != nil {
		return nil, fmt.Errorf("failed to make new consume channel: %v", err)
	}

	return ctl, nil
}

func (c *Controller) AddProducer(producerExchangeName string) error {
	ch, err := NewProducerChannel(c.conn, producerExchangeName)
	if err != nil {
		return fmt.Errorf("failed to make producer channel: %v", err)
	}
	c.produceQueue[producerExchangeName] = ch
	return nil
}

func (c *Controller) Push(producerExchangeName string, key string, msg string) error {
	err := c.produceQueue[producerExchangeName].PublishWithContext(
		context.Background(),
		producerExchangeName,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish msg: %v", err)
	}
	return nil
}

func (c *Controller) StartConsume(callback func(delivery amqp.Delivery) error) error {
	for d := range c.consumeChannel {
		err := callback(d)
		if err != nil {
			fmt.Printf("failed to call func: %v", err)
		}
	}
	return nil
}

func (c *Controller) Close() error {
	for _, queue := range c.produceQueue {
		if err := queue.Close(); err != nil {
			return fmt.Errorf("failed to close producer queue: %v", err)
		}
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed to close conn: %v", err)
	}
	return nil
}
