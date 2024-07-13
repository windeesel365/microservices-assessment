package main

import (
	"context"

	pb "github.com/windeesel365/microservices-assessment/services/payment/paymentpb"

	orderpb "github.com/windeesel365/microservices-assessment/services/order/orderpb"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Server struct {
	pb.UnimplementedPaymentServiceServer
	DB          *sqlx.DB
	OrderClient orderpb.OrderServiceClient
}

func (s *Server) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	var id int
	err := s.DB.QueryRowContext(ctx, `
        INSERT INTO payments (order_id, amount, payment_method) 
        VALUES ($1, $2, $3) RETURNING id`, req.OrderId, req.Amount, req.PaymentMethod).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &pb.CreatePaymentResponse{Id: int32(id)}, nil
}

func (s *Server) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {
	payment := pb.GetPaymentResponse{}
	err := s.DB.GetContext(ctx, &payment, `
        SELECT id, order_id, amount, payment_method, status, transaction_id, created_at 
        FROM payments WHERE id = $1`, req.Id)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (s *Server) UpdatePayment(ctx context.Context, req *pb.UpdatePaymentRequest) (*pb.UpdatePaymentResponse, error) {
	result, err := s.DB.ExecContext(ctx, `
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

func (s *Server) DeletePayment(ctx context.Context, req *pb.DeletePaymentRequest) (*pb.DeletePaymentResponse, error) {
	result, err := s.DB.ExecContext(ctx, `DELETE FROM payments WHERE id = $1`, req.Id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &pb.DeletePaymentResponse{Success: rowsAffected > 0}, nil
}
