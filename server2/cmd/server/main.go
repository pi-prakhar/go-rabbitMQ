package main

import (
	"go-rabbitmq-server2/api"
	"go-rabbitmq-server2/internal/models"
	"go-rabbitmq-server2/internal/rabbitmq"
	"log"
	"net/http"
	"time"
)

const (
	CONFIG_FILE        = "config/config.yml"
	TEST_QUEUE         = "test"
	MESSAGE_EXCHANGE   = "messages"
	BROADCAST_EXCHANGE = "broadcast"
)

func main() {
	ConfigData := &models.ConfigData{}
	if err := ConfigData.LoadConfig(CONFIG_FILE); err != nil {
		log.Fatal("Error: Loading config file", err)
	}

	config := ConfigData.GetConfig()

	conn, connStr, err := rabbitmq.Connect(ConfigData.RabbitMQ.Host, ConfigData.RabbitMQ.User, ConfigData.RabbitMQ.Password, ConfigData.RabbitMQ.RetryCount, ConfigData.RabbitMQ.RetrySleep)

	if err != nil {
		log.Fatalf("Error : Failed to connect to RabbitMQ Server : %s %v", connStr, err)
	}

	defer conn.Close()

	var consumer rabbitmq.RMQConsumer
	consumer, err = rabbitmq.NewConsumer(conn)
	if err != nil {
		log.Fatal("Error : Failed to setup new rabbitMq channel", err)
	}

	defer consumer.GetChannel().Close()

	err = consumer.DeclareExchange(MESSAGE_EXCHANGE, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Error : Failed to declare exchange : ", MESSAGE_EXCHANGE, err)
	}

	q, err := consumer.DeclareQueue(TEST_QUEUE, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Error : Failed to declare queue : ", TEST_QUEUE, err)
	}

	err = consumer.Listen(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal("Error : Failed to start listener for queue : ", q.Name)
	}

	qTemp1, err := consumer.DeclareQueue("", false, false, true, false, nil)
	if err != nil {
		log.Fatal("Error : Failed to declare queue temporary queue 1", err)
	}

	err = consumer.BindQueueToChannel(qTemp1.Name, "", MESSAGE_EXCHANGE, false, nil)
	if err != nil {
		log.Fatalf("Error : Failed to bind queue : %s to exchange : %s", qTemp1.Name, BROADCAST_EXCHANGE)
	}

	err = consumer.Listen(qTemp1.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal("Error : Failed to start listener for queue : ", qTemp1.Name)
	}

	qTemp2, err := consumer.DeclareQueue("", false, false, true, false, nil)
	if err != nil {
		log.Fatal("Error : Failed to declare queue temporary queue", err)
	}

	err = consumer.BindQueueToChannel(qTemp2.Name, "1", MESSAGE_EXCHANGE, false, nil)
	if err != nil {
		log.Fatalf("Error : Failed to bind queue : %s to exchange : %s", qTemp2.Name, BROADCAST_EXCHANGE)
	}

	err = consumer.ListenWithPriority(qTemp2.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal("Error : Failed to start listener for queue : ", qTemp2.Name)
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

	log.Printf("INFO : %s started on port %s in %s mode", ConfigData.Server.Name, config.Port, ConfigData.Server.Mode)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server at port %s", config.Port)
	}

}
