package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	connection *amqp.Connection
	Channel    *amqp.Channel
}

func (p *Producer) setup() error {
	channel, err := p.connection.Channel()

	if err != nil {
		return err
	}
	p.Channel = channel
	return nil
}

func (p *Producer) DeclareQueue(name string, durable bool, deleteWhenUnused bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	q, err := p.Channel.QueueDeclare(
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

func (p *Producer) Push(ctx context.Context, qName string, body string, exchangeName string, mandatory bool, immediate bool) error {
	err := p.Channel.PublishWithContext(ctx,
		exchangeName,
		qName,
		mandatory,
		immediate,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})

	if err != nil {
		return err
	}
	return nil
}

func NewProducer(conn *amqp.Connection) (Producer, error) {
	producer := Producer{
		connection: conn,
	}

	err := producer.setup()
	if err != nil {
		return Producer{}, err
	}

	return producer, err
}
