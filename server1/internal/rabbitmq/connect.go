package rabbitmq

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func getRabbitMQConnectionString(host string, user string, password string) string {
	connString := fmt.Sprintf("amqp://%s:%s@%s:5672/", user, password, host)
	return connString
}

func Connect(host string, user string, password string, retryCount int, retrySleep int) (*amqp.Connection, string, error) {
	connStr := getRabbitMQConnectionString(host, user, password)
	var conn *amqp.Connection
	var err error

	for i := 0; i < retryCount; i++ {
		conn, err = amqp.Dial(connStr)
		if err == nil {
			break
		}
		log.Printf("INFO : Failed to connect to RabbitMQ, retrying in 5 seconds... (%v)", err)
		time.Sleep(time.Second * time.Duration(retryCount))
	}

	return conn, connStr, err
}
