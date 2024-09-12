package api

import (
	"bytes"
	"context"
	"encoding/json"
	"go-rabbitmq-server1/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestApp_handleTest(t *testing.T) {
	// Setup
	mockProducer := new(MockProducer)
	app := App{
		Producer:  mockProducer,
		QueueTest: &amqp.Queue{Name: "test"},
	}
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	rr := httptest.NewRecorder()
	handler := app.NewRouter()

	// Test
	handler.ServeHTTP(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Hello"

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestApp_handlePushSuccess(t *testing.T) {
	// Setup
	mockProducer := new(MockProducer)
	app := App{
		Producer:  mockProducer,
		QueueTest: &amqp.Queue{Name: "test"},
	}

	// Create a test message
	message := models.Message{
		Message: "test message",
	}
	body, _ := json.Marshal(message)

	// Create a request with the message as body
	req, err := http.NewRequest("POST", "/push", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	// req.Header.Set("Content-Type", "application/json") // Set content type if needed

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the function being tested
	handler := app.NewRouter()

	// Test
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"message":"Succesfull Pushed Message to queue test","code":200,"data":"test message"}`
	got := strings.TrimSpace(rr.Body.String())

	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", got, expected)
	}
}

func TestApp_handlePushFailBadRequest(t *testing.T) {
	// Setup
	mockProducer := new(MockProducer)
	app := App{
		Producer:  mockProducer,
		QueueTest: &amqp.Queue{Name: "test"},
	}

	// Create a test message
	incorrectBody := make(map[string]interface{})
	incorrectBody["msg"] = "test_message"

	body, _ := json.Marshal(incorrectBody)

	// Create a request with the message as body
	req, err := http.NewRequest("POST", "/push", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	// req.Header.Set("Content-Type", "application/json") // Set content type if needed

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the function being tested
	handler := app.NewRouter()

	// Test
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check the response body
	expected := `{"message":"Failed to decode request body","code":400,"error":"json: unknown field \"msg\""}`
	got := strings.TrimSpace(rr.Body.String())

	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", got, expected)
	}
}

func TestApp_handlePushFailMethodNOtAllowed(t *testing.T) {
	// Setup
	mockProducer := new(MockProducer)
	app := App{
		Producer:  mockProducer,
		QueueTest: &amqp.Queue{Name: "test"},
	}

	// Create a test message
	incorrectBody := make(map[string]interface{})
	incorrectBody["msg"] = "test_message"

	body, _ := json.Marshal(incorrectBody)

	// Create a request with the message as body
	req, err := http.NewRequest("GET", "/push", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	// req.Header.Set("Content-Type", "application/json") // Set content type if needed

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the function being tested
	handler := app.NewRouter()

	// Test
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
