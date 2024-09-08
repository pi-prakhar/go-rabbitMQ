package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQProducer interface {
	DeclareQueue(name string, durable bool, deleteWhenUnused bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error)
	Push(ctx context.Context, qName string, body string, exchangeName string, mandatory bool, immediate bool) error
	GetChannel() *amqp.Channel
}

type Producer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func (p *Producer) setup() error {
	channel, err := p.connection.Channel()

	if err != nil {
		return err
	}
	p.channel = channel
	return nil
}

func (p *Producer) GetChannel() *amqp.Channel {
	return p.channel
}

func (p *Producer) DeclareQueue(name string, durable bool, deleteWhenUnused bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	q, err := p.channel.QueueDeclare(
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
	err := p.channel.PublishWithContext(ctx,
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

func NewProducer(conn *amqp.Connection) (RMQProducer, error) {
	producer := &Producer{
		connection: conn,
	}

	err := producer.setup()
	if err != nil {
		return nil, err
	}

	return producer, err
}
