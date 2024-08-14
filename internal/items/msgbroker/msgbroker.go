package msgbroker

import (
	"context"
	"log/slog"
	"sync"

	"budgeting-service/internal/items/service"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	budget_pb "budgeting-service/genproto/budget"
	goal_pb "budgeting-service/genproto/goal"
	notification_pb "budgeting-service/genproto/notification"
	transaction_pb "budgeting-service/genproto/transaction"
)

type MsgBroker struct {
	service *service.Service
	readers *MsgBrokers
	logger  *slog.Logger
	wg      *sync.WaitGroup
}

func New(service *service.Service, logger *slog.Logger, readers *MsgBrokers, wg *sync.WaitGroup) *MsgBroker {
	return &MsgBroker{
		service: service,
		readers: readers,
		logger:  logger,
		wg:      wg,
	}
}

func (m *MsgBroker) StartToConsume(ctx context.Context) {
	m.wg.Add(4)

	consumerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go m.consumeMessages(consumerCtx, m.readers.TransactionCreated, "transaction_created")
	go m.consumeMessages(consumerCtx, m.readers.BudgetUpdated, "budget_updated")
	go m.consumeMessages(consumerCtx, m.readers.GoalProgressUpdated, "goal_progress_updated")
	go m.consumeMessages(consumerCtx, m.readers.NotificationCreated, "notification_created")

	<-consumerCtx.Done()
	m.logger.Info("All consumers have stopped")
}

func (m *MsgBroker) consumeMessages(ctx context.Context, reader *kafka.Reader, logPrefix string) {
	defer m.wg.Done()
	for {
		select {
		case <-ctx.Done():
			m.logger.Info("Context done, stopping consumer", "consumer", logPrefix)
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				m.logger.Error("Error reading message", "error", err, "topic", logPrefix)
				return
			}

			var response proto.Message
			var errUnmarshal error

			switch logPrefix {
			case "transaction_created":
				var req transaction_pb.CreateTransactionRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = m.service.TransactionService.CreateTransaction(ctx, &req)
			case "budget_updated":
				var req budget_pb.UpdateBudgetRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = m.service.BudgetService.UpdateBudget(ctx, &req)
			case "goal_progress_updated":
				var req goal_pb.UpdateGoalRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = m.service.GoalService.UpdateGoal(ctx, &req)
			case "notification_created":
				var req notification_pb.CreateNotificationRequest
				errUnmarshal = protojson.Unmarshal(msg.Value, &req)
				response, err = m.service.NotificationService.CreateNotification(ctx, &req)
			}

			if errUnmarshal != nil {
				m.logger.Error("Error while unmarshaling data", "error", errUnmarshal)
				continue
			}

			if err != nil {
				m.logger.Error("Failed in %s: %s\n", logPrefix, err.Error())
				continue
			}

			_, err = proto.Marshal(response)
			if err != nil {
				m.logger.Error("Failed to marshal response", "error", err)
				continue
			}

			m.logger.Info("Successfully processed message", "topic", logPrefix)
		}
	}
}
