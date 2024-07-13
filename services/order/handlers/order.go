package handlers

import (
	"context"
	"time"

	pb "github.com/windeesel365/microservices-assessment/services/order/orderpb"

	productpb "github.com/windeesel365/microservices-assessment/services/product/productpb"
	userpb "github.com/windeesel365/microservices-assessment/services/user/userpb"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/windeesel365/microservices-assessment/services/order/models"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	DB            *sqlx.DB
	userClient    userpb.UserServiceClient
	productClient productpb.ProductServiceClient
}

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	query := `INSERT INTO orders (user_id, total_amount, status) VALUES ($1, $2, 'PENDING') RETURNING id`
	var orderId int32
	err := s.DB.QueryRowx(query, req.UserId, req.TotalAmount).Scan(&orderId)
	if err != nil {
		return nil, err
	}

	for _, item := range req.Items {
		query := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)`
		_, err = s.DB.Exec(query, orderId, item.ProductId, item.Quantity, item.Price)
		if err != nil {
			return nil, err
		}
	}

	return &pb.CreateOrderResponse{UserId: req.UserId}, nil

}

func (s *Server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	query := `SELECT id, user_id, total_amount, status, created_at, updated_at FROM orders WHERE id = $1`
	var order models.OrderResponse
	err := s.DB.Get(&order, query, req.Id)
	if err != nil {
		return nil, err
	}

	query = `SELECT product_id, quantity, price FROM order_items WHERE order_id = $1`
	var items []models.OrderItem
	err = s.DB.Select(&items, query, req.Id)
	if err != nil {
		return nil, err
	}

	order.Items = items

	return &pb.GetOrderResponse{
		Id:          order.Id,
		UserId:      order.UserId,
		Items:       convertOrderItemsToPB(items),
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}, nil
}

func (s *Server) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := s.DB.Exec(query, req.Status, time.Now(), req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateOrderResponse{Success: true}, nil
}

func (s *Server) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := s.DB.Exec(query, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteOrderResponse{Success: true}, nil
}

func convertOrderItemsToPB(items []OrderItem) []*pb.OrderItem {
	var pbItems []*pb.OrderItem
	for _, item := range items {
		pbItems = append(pbItems, &pb.OrderItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}
	return pbItems
}
