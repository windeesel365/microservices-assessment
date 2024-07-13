package main

import (
	"log"
	"net"
	"os"

	pb "github.com/windeesel365/microservices-assessment/services/order/orderpb"

	productpb "github.com/windeesel365/microservices-assessment/services/product/productpb"
	userpb "github.com/windeesel365/microservices-assessment/services/user/userpb"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/windeesel365/microservices-assessment/services/order/config"
	"github.com/windeesel365/microservices-assessment/services/order/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg := config.LoadConfig()

	userServiceAddr := os.Getenv("USER_SERVICE_PORT")
	if userServiceAddr == "" {
		log.Fatalf("USER_SERVICE_PORT is not set in the environment")
	}

	productServiceAddr := os.Getenv("PRODUCT_SERVICE_PORT")
	if productServiceAddr == "" {
		log.Fatalf("PRODUCT_SERVICE_PORT is not set in the environment")
	}

	//ต่อกับ userService
	userConn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}

	defer userConn.Close()

	//ต่อกับ productService
	productConn, err := grpc.Dial(productServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to product service: %v", err)
	}

	defer productConn.Close()

	userClient := userpb.NewUserServiceClient(userConn)
	productClient := productpb.NewProductServiceClient(productConn)

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	pb.RegisterOrderServiceServer(s, &handlers.Server{
		db:            db,
		userClient:    userClient,
		productClient: productClient,
	})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
