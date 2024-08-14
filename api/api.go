package api

import (
	"log"
	"net"

	"budgeting-service/internal/items/config"
	"budgeting-service/internal/items/service"

	account_pb "budgeting-service/genproto/account"
	budget_pb "budgeting-service/genproto/budget"
	category_pb "budgeting-service/genproto/category"
	goal_pb "budgeting-service/genproto/goal"
	notification_pb "budgeting-service/genproto/notification"
	report_pb "budgeting-service/genproto/report"
	transaction_pb "budgeting-service/genproto/transaction"

	"google.golang.org/grpc"
)

type API struct {
	service *service.Service
}

func New(service *service.Service) *API {
	return &API{
		service: service,
	}
}

func (a *API) RUN(config *config.Config, service *service.Service) error {
	listener, err := net.Listen("tcp", "budgeting"+config.Server.Port)
	if err != nil {
		return err
	}

	serverRegisterer := grpc.NewServer()

	account_pb.RegisterAccountServiceServer(serverRegisterer, service.AccountService)
	budget_pb.RegisterBudgetServiceServer(serverRegisterer, service.BudgetService)
	category_pb.RegisterCategoryServiceServer(serverRegisterer, service.CategoryService)
	goal_pb.RegisterGoalServiceServer(serverRegisterer, service.GoalService)
	notification_pb.RegisterNotificationServiceServer(serverRegisterer, service.NotificationService)
	report_pb.RegisterReportServiceServer(serverRegisterer, service.ReportService)
	transaction_pb.RegisterTransactionServiceServer(serverRegisterer, service.TransactionService)

	log.Println("Server has started running on port:", config.Server.Port)

	return serverRegisterer.Serve(listener)
}
