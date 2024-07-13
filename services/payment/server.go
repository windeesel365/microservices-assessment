package main

import (
	"log"
	"net"
	"os"

	pb "github.com/windeesel365/microservices-assessment/services/payment/paymentpb"

	orderpb "github.com/windeesel365/microservices-assessment/services/order/orderpb"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/windeesel365/microservices-assessment/services/payment/config"
	"github.com/windeesel365/microservices-assessment/services/payment/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg := config.LoadConfig()

	orderServiceAddr := os.Getenv("ORDER_SERVICE_PORT")
	if orderServiceAddr == "" {
		log.Fatalf("ORDER_SERVICE_PORT is not set in the environment")
	}

	//ต่อกับ orderService
	orderConn, err := grpc.Dial(orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to order service: %v", err)
	}

	defer orderConn.Close()

	orderClient := orderpb.NewOrderServiceClient(orderConn)

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	pb.RegisterPaymentServiceServer(s, &handlers.Server{
		DB:          db,
		OrderClient: orderClient,
	})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
