package api

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MockProducer struct {
	SetupCalled        bool
	DeclareQueueCalled bool
	PushCalled         bool
	GetChannelCalled   bool
	MockChannel        *amqp.Channel
}

func (m *MockProducer) DeclareQueue(name string, durable bool, deleteWhenUnused bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	m.PushCalled = true
	return amqp.Queue{}, nil
}

func (m *MockProducer) Push(ctx context.Context, qName string, body string, exchangeName string, mandatory bool, immediate bool) error {
	m.DeclareQueueCalled = true
	return nil
}
func (m *MockProducer) GetChannel() *amqp.Channel {
	return m.MockChannel
}

// func TestAppSendMessageToQueue(t *testing.T) {
// 	MockProducer := new(MockProducer)
// 	app := &App{
// 		Producer:       MockProducer,
// 		QueueTest:      &amqp.Queue{Name: "test-queue"},
// 		HandlerTimeout: 30,
// 	}

// }
