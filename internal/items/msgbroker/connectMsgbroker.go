package msgbroker

import (
	"budgeting-service/internal/items/config"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MsgBrokers struct {
	Transaction_created   <-chan amqp.Delivery
	Budget_updated        <-chan amqp.Delivery
	Goal_progress_updated <-chan amqp.Delivery
	Notification_created  <-chan amqp.Delivery
}

func InitMessageBroker(config *config.Config) (*MsgBrokers, *amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(config.RabbitMQ.RabbitMQ)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	return &MsgBrokers{
		Transaction_created:   subscribeToQueue(ch, "transaction_created"),
		Budget_updated:        subscribeToQueue(ch, "budget_updated"),
		Goal_progress_updated: subscribeToQueue(ch, "goal_progress_updated"),
		Notification_created:  subscribeToQueue(ch, "notification_created"),
	}, conn, ch
}

func subscribeToQueue(ch *amqp.Channel, queueName string) <-chan amqp.Delivery {
	queue, err := getQueue(ch, queueName)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := getMessageQueue(ch, queue)
	if err != nil {
		log.Fatal(err)
	}

	return msgs
}

func getQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}

func getMessageQueue(ch *amqp.Channel, q amqp.Queue) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}
