package main

import (
	"go-rabbitmq-server1/api"
	"go-rabbitmq-server1/internal/models"
	"go-rabbitmq-server1/internal/rabbitmq"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	CONFIG_FILE = "config/config.yml"
	TEST_QUEUE  = "test"
)

func main() {
	ConfigData := &models.ConfigData{}
	if err := ConfigData.LoadConfig(CONFIG_FILE); err != nil {
		log.Fatal("Error : Loading config file", err)
	}

	config := ConfigData.GetConfig()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Fatal("Error : Failed to connect to RabbitMQ", err)
	}

	defer conn.Close()

	producer, err := rabbitmq.NewProducer(conn)
	if err != nil {
		log.Fatal("Error : Failed to setup new rabbitMq channel", err)
	}

	defer producer.Channel.Close()

	q, err := producer.DeclareQueue(TEST_QUEUE, false, false, false, false, nil)
	if err != nil {
		log.Fatal("Error : Failed to declare queue : ", TEST_QUEUE, err)
	}

	app := &api.App{
		Producer:       &producer,
		QueueTest:      &q,
		HandlerTimeout: config.HandlerTimeout,
	}

	srv := &http.Server{
		Addr:         config.Port,
		Handler:      app.NewRouter(),
		ReadTimeout:  time.Second * time.Duration(config.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(config.WriteTimeout),
	}

	log.Printf("INFO : %s started on port %s", ConfigData.Server.Name, config.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error : Failed to start server at port %s", config.Port)
	}
}
