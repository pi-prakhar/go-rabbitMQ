package main

import (
	"go-rabbitmq-server2/api"
	"go-rabbitmq-server2/internal/models"
	"go-rabbitmq-server2/internal/rabbitmq"
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
		log.Fatal("Error: Loading config file", err)
	}

	config := ConfigData.GetConfig()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Fatal("Error : Failed to connect to RabbitMQ", err)
	}

	defer conn.Close()

	consumer, err := rabbitmq.NewConsumer(conn)
	if err != nil {
		log.Fatal("Error : Failed to setup new rabbitMq channel", err)
	}

	defer consumer.Channel.Close()

	q, err := consumer.DeclareQueue(TEST_QUEUE, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Error : Failed to declare queue : ", TEST_QUEUE, err)
	}

	err = consumer.Listen(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal("Error : Failed to register a consumer")
	}

	app := &api.App{
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
		log.Fatalf("Failed to start server at port %s", config.Port)
	}

}
