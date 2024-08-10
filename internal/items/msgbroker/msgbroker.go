package msgbroker

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"budgeting-service/internal/items/service"

	budget_pb "budgeting-service/genproto/budget"
	goal_pb "budgeting-service/genproto/goal"
	notification_pb "budgeting-service/genproto/notification"
	transaction_pb "budgeting-service/genproto/transaction"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type MsgBroker struct {
	service          *service.Service
	msgs             *MsgBrokers
	logger           *slog.Logger
	wg               *sync.WaitGroup
	numberOfServices int
}

func New(service *service.Service,
	logger *slog.Logger,
	msgs *MsgBrokers,
	wg *sync.WaitGroup,
	numberOfServices int) *MsgBroker {
	return &MsgBroker{
		service:          service,
		msgs:             msgs,
		logger:           logger,
		wg:               wg,
		numberOfServices: numberOfServices,
	}
}

func (m *MsgBroker) StartToConsume(ctx context.Context, contentType string) {
	m.wg.Add(m.numberOfServices)
	consumerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go m.consumeMessages(consumerCtx, m.msgs.Transaction_created, "transaction_created")
	go m.consumeMessages(consumerCtx, m.msgs.Budget_updated, "budget_updated")
	go m.consumeMessages(consumerCtx, m.msgs.Goal_progress_updated, "goal_progress_updated")
	go m.consumeMessages(consumerCtx, m.msgs.Notification_created, "notification_created")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	m.logger.Info("Shutting down, waiting for consumers to finish")
	cancel()
	m.wg.Wait()
	m.logger.Info("All consumers have stopped")
}

func (m *MsgBroker) consumeMessages(ctx context.Context, messages <-chan amqp.Delivery, logPrefix string) {
	defer m.wg.Done()
	for {
		select {
		case val := <-messages:
			var response proto.Message
			var err error

			switch logPrefix {
			case "transaction_created":
				var req transaction_pb.CreateTransactionRequest
				if err := json.Unmarshal(val.Body, &req); err != nil {
					m.logger.Error("Error while unmarshaling data", "error", err)
					val.Nack(false, false)
					continue
				}
				response, err = m.service.TransactionService.CreateTransaction(ctx, &req)
			case "budget_updated":
				var req budget_pb.UpdateBudgetRequest
				if err := json.Unmarshal(val.Body, &req); err != nil {
					m.logger.Error("Error while unmarshaling data", "error", err)
					val.Nack(false, false)
					continue
				}
				response, err = m.service.BudgetService.UpdateBudget(ctx, &req)
			case "goal_progress_updated":
				var req goal_pb.UpdateGoalRequest
				if err := json.Unmarshal(val.Body, &req); err != nil {
					m.logger.Error("Error while unmarshaling data", "error", err)
					val.Nack(false, false)
					continue
				}
				response, err = m.service.GoalService.UpdateGoal(ctx, &req)
			case "notification_created":
				var req notification_pb.GetNotificationsRequest
				if err := json.Unmarshal(val.Body, &req); err != nil {
					m.logger.Error("Error while unmarshaling data", "error", err)
					val.Nack(false, false)
					continue
				}
				response, err = m.service.NotificationService.GetNotifications(ctx, &req)
			}

			if err != nil {
				m.logger.Error("Failed in %s: %s\n", logPrefix, err.Error())
				val.Nack(false, false)
				continue
			}

			val.Ack(false)

			_, err = proto.Marshal(response)
			if err != nil {
				m.logger.Error("Failed to marshal response", "error", err)
				continue
			}

		case <-ctx.Done():
			m.logger.Info("Context done, stopping consumer", "consumer", logPrefix)
			return
		}
	}
}
