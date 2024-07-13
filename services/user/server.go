package main

import (
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/windeesel365/microservices-assessment/services/user/config"
	"github.com/windeesel365/microservices-assessment/services/user/handlers"
	pb "github.com/windeesel365/microservices-assessment/services/user/userpb"

	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg := config.LoadConfig()

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	pb.RegisterUserServiceServer(s, &handlers.Server{DB: db})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
