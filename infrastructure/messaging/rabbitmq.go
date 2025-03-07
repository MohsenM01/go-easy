package messaging

import (
	"go-easy/config"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	router *MessageRouter
}

func NewMessageQueue() *MessageQueue {

	cfg := config.LoadConfig()

	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	router := NewMessageRouter()

	return &MessageQueue{conn: conn, ch: ch, router: router}
}

func (mq *MessageQueue) StartMessageQueue() {
	mq.declareQueues()
	mq.startConsumers()
}

func (mq *MessageQueue) Close() {
	mq.conn.Close()
	mq.ch.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
