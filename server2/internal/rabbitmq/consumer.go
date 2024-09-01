package rabbitmq

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	connection *amqp.Connection
	Channel    *amqp.Channel
}

func (c *Consumer) setup() error {
	channel, err := c.connection.Channel()

	if err != nil {
		return err
	}
	c.Channel = channel
	return nil
}

func (c *Consumer) DeclareQueue(name string, durable bool, deleteWhenUnused bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	q, err := c.Channel.QueueDeclare(
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
	msgs, err := c.Channel.Consume(
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
			time.Sleep(10 * time.Second)
			msg.Ack(false) // if this is not made false the message will still remain in the queue.
			log.Printf("INFO : Message : %s Successfully processed", message)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	return nil
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		connection: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, err
}
