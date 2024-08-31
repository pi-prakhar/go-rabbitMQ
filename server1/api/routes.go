package api

import (
	"go-rabbitmq-server1/internal/rabbitmq"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type App struct {
	Producer       *rabbitmq.Producer
	QueueTest      *amqp.Queue
	HandlerTimeout int
}

func (app *App) NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/test", app.handleTest)
	mux.HandleFunc("/push", app.sendMessageToQueue)

	return app.corsMiddleware(mux)
}
