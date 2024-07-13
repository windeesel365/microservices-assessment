package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/windeesel365/microservices-assessment/services/payment/paymentpb"

	orderpb "github.com/windeesel365/microservices-assessment/services/order/orderpb"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedPaymentServiceServer
	db          *sqlx.DB
	orderClient orderpb.UserServiceClient
}

type PaymentResponse struct {
	Id            int32   `db:"id"`
	OrderId       int32   `db:"order_id"`
	Amount        float32 `db:"amount"`
	PaymentMethod string  `db:"payment_method"`
	Status        string  `db:"status"`
	TransactionId string  `db:"transaction_id"`
	CreatedAt     string  `db:"created_at"`
}

func (s *server) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	var id int
	err := s.db.QueryRowContext(ctx, `
        INSERT INTO payments (order_id, amount, payment_method) 
        VALUES ($1, $2, $3) RETURNING id`, req.OrderId, req.Amount, req.PaymentMethod).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &pb.CreatePaymentResponse{Id: int32(id)}, nil
}

func (s *server) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {
	payment := pb.GetPaymentResponse{}
	err := s.db.GetContext(ctx, &payment, `
        SELECT id, order_id, amount, payment_method, status, transaction_id, created_at 
        FROM payments WHERE id = $1`, req.Id)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (s *server) UpdatePayment(ctx context.Context, req *pb.UpdatePaymentRequest) (*pb.UpdatePaymentResponse, error) {
	result, err := s.db.ExecContext(ctx, `
        UPDATE payments SET status = $1, transaction_id = $2 WHERE id = $3`,
		req.Status, req.TransactionId, req.Id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &pb.UpdatePaymentResponse{Success: rowsAffected > 0}, nil
}

func (s *server) DeletePayment(ctx context.Context, req *pb.DeletePaymentRequest) (*pb.DeletePaymentResponse, error) {
	result, err := s.db.ExecContext(ctx, `DELETE FROM payments WHERE id = $1`, req.Id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &pb.DeletePaymentResponse{Success: rowsAffected > 0}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT is not set in the environment")
	}

	dbSource := os.Getenv("DATABASE_URL")
	if dbSource == "" {
		log.Fatalf("DATABASE_URL is not set in the environment")
	}

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

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	db, err := sqlx.Connect("postgres", dbSource)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	pb.RegisterPaymentServiceServer(s, &server{
		db:          db,
		orderClient: orderClient,
	})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
