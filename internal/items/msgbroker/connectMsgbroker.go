package msgbroker

import (
	"budgeting-service/internal/items/config"

	"github.com/segmentio/kafka-go"
)

type MsgBrokers struct {
	TransactionCreated  *kafka.Reader
	BudgetUpdated       *kafka.Reader
	GoalProgressUpdated *kafka.Reader
	NotificationCreated *kafka.Reader
}

func InitMessageBroker(config *config.Config) *MsgBrokers {
	readers := map[string]*kafka.Reader{
		"transaction_created": kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{config.Kafka.Broker},
			Topic:   "transaction_created",
			GroupID: "budgeting_service",
		}),
		"budget_updated": kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{config.Kafka.Broker},
			Topic:   "budget_updated",
			GroupID: "budgeting_service",
		}),
		"goal_progress_updated": kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{config.Kafka.Broker},
			Topic:   "goal_progress_updated",
			GroupID: "budgeting_service",
		}),
		"notification_created": kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{config.Kafka.Broker},
			Topic:   "notification_created",
			GroupID: "budgeting_service",
		}),
	}

	return &MsgBrokers{
		TransactionCreated:  readers["transaction_created"],
		BudgetUpdated:       readers["budget_updated"],
		GoalProgressUpdated: readers["goal_progress_updated"],
		NotificationCreated: readers["notification_created"],
	}
}
