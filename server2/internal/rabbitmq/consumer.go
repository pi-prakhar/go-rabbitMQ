package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQConsumer interface {
	DeclareExchange(name string, exchangeType string, durable bool, autoDeleted bool, internal bool, noWait bool, arguments amqp.Table) error
	DeclareQueue(name string, durable bool, deleteWhenUnused bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error)
	Listen(queueName string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) error
	ListenWithPriority(queueName string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) error
	BindQueueToChannel(queueName string, routingKey string, exchangeName string, noWait bool, args amqp.Table) error
	GetChannel() *amqp.Channel
}

type Consumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func (c *Consumer) DeclareExchange(name string, exchangeType string, durable bool, autoDeleted bool, internal bool, noWait bool, arguments amqp.Table) error {
	err := c.channel.ExchangeDeclare(
		name,
		exchangeType,
		durable,
		autoDeleted,
		internal,
		noWait,
		arguments,
	)

	return err
}

func (c *Consumer) DeclareQueue(name string, durable bool, deleteWhenUnused bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	q, err := c.channel.QueueDeclare(
		name,
		durable,
		deleteWhenUnused,
		exclusive,
		noWait,
		args,
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	return q, nil
}

func (c *Consumer) Listen(queueName string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) error {
	msgs, err := c.channel.Consume(
		queueName,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
	)

	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			message := string(msg.Body)
			log.Printf("INFO : Message : %s Recieved from queue : %s", message, queueName)
			//time.Sleep(10 * time.Second)
			msg.Ack(false) // if this is not made false the message will still remain in the queue.
			log.Printf("INFO : Message : %s from queue : %s Successfully processed", message, queueName)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	return nil
}

func (c *Consumer) ListenWithPriority(queueName string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) error {
	msgs, err := c.channel.Consume(
		queueName,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
	)

	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			message := string(msg.Body)
			log.Printf("INFO : Message : %s Recieved from queue : %s with priority", message, queueName)
			//time.Sleep(10 * time.Second)
			msg.Ack(false) // if this is not made false the message will still remain in the queue.
			log.Printf("INFO : Message : %s from queue : %s Successfully processed with priority", message, queueName)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	return nil
}

func (c *Consumer) GetChannel() *amqp.Channel {
	return c.channel
}

func (c *Consumer) BindQueueToChannel(queueName string, routingKey string, exchangeName string, noWait bool, args amqp.Table) error {
	err := c.channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		noWait,
		args,
	)
	return err
}

func (c *Consumer) setup() error {
	channel, err := c.connection.Channel()

	if err != nil {
		return err
	}
	c.channel = channel
	return nil
}

func NewConsumer(conn *amqp.Connection) (RMQConsumer, error) {
	consumer := &Consumer{
		connection: conn,
	}

	err := consumer.setup()
	if err != nil {
		return nil, err
	}

	return consumer, err
}
