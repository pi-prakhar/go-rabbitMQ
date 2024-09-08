package main

import (
	"go-rabbitmq-server1/api"
	"go-rabbitmq-server1/internal/models"
	"go-rabbitmq-server1/internal/rabbitmq"
	"log"
	"net/http"
	"time"
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

	conn, connStr, err := rabbitmq.Connect(ConfigData.RabbitMQ.Host, ConfigData.RabbitMQ.User, ConfigData.RabbitMQ.Password, ConfigData.RabbitMQ.RetryCount, ConfigData.RabbitMQ.RetrySleep)

	if err != nil {
		log.Fatalf("Error : Failed to connect to RabbitMQ Server : %s %v", connStr, err)
	}

	defer conn.Close()

	producer, err := rabbitmq.NewProducer(conn)
	if err != nil {
		log.Fatal("Error : Failed to setup new rabbitMq channel", err)
	}

	defer producer.GetChannel().Close()

	q, err := producer.DeclareQueue(TEST_QUEUE, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Error : Failed to declare queue : ", TEST_QUEUE, err)
	}

	app := &api.App{
		Producer:       producer,
		QueueTest:      &q,
		HandlerTimeout: config.HandlerTimeout,
	}

	srv := &http.Server{
		Addr:         config.Port,
		Handler:      app.NewRouter(),
		ReadTimeout:  time.Second * time.Duration(config.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(config.WriteTimeout),
	}

	log.Printf("INFO : %s started on port %s in %s mode", ConfigData.Server.Name, config.Port, ConfigData.Server.Mode)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error : Failed to start server at port %s", config.Port)
	}
}
