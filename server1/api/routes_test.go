package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestApp_handleTest(t *testing.T) {
	// Setup
	mockProducer := new(MockProducer)
	app := App{
		Producer: mockProducer,
	}

	// Create request and recorder
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(app.handleTest)
	handler.ServeHTTP(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Test successful", rr.Body.String())
}

func TestApp_sendMessageToQueue(t *testing.T) {
	// Setup
	mockProducer := new(MockProducer)
	mockQueue := &amqp.Queue{Name: "test-queue"}
	app := App{
		Producer:  mockProducer,
		QueueTest: mockQueue,
	}

	// Expect the producer to call Push with specific arguments
	mockProducer.On("Push", mock.Anything, "test-queue", "message body", "", false, false).Return(nil)

	// Create request and recorder
	req := httptest.NewRequest("POST", "/push", nil)
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(app.sendMessageToQueue)
	handler.ServeHTTP(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Message sent", rr.Body.String())

	// Assert that the mock was called as expected
	mockProducer.AssertExpectations(t)
}
