package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	pb "ecommerce-order/orderpb"

	productpb "github.com/windeesel365/microservices-assessment/services/product/productpb"
	userpb "github.com/windeesel365/microservices-assessment/services/user/userpb"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedOrderServiceServer
	db            *sqlx.DB
	userClient    userpb.UserServiceClient
	productClient productpb.ProductServiceClient
}

type OrderItem struct {
	ProductId int32   `db:"product_id"`
	Quantity  int32   `db:"quantity"`
	Price     float32 `db:"price"`
}

type OrderResponse struct {
	Id          int32       `db:"id"`
	UserId      int32       `db:"user_id"`
	Items       []OrderItem `db:"items"`
	TotalAmount float32     `db:"total_amount"`
	Status      string      `db:"status"`
	CreatedAt   string      `db:"created_at"`
	UpdatedAt   string      `db:"updated_at"`
}

func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	query := `INSERT INTO orders (user_id, total_amount, status) VALUES ($1, $2, 'PENDING') RETURNING id`
	var orderId int32
	err := s.db.QueryRowx(query, req.UserId, req.TotalAmount).Scan(&orderId)
	if err != nil {
		return nil, err
	}

	for _, item := range req.Items {
		query := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)`
		_, err = s.db.Exec(query, orderId, item.ProductId, item.Quantity, item.Price)
		if err != nil {
			return nil, err
		}
	}

	return &pb.CreateOrderResponse{UserId: req.UserId}, nil
}

func (s *server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	query := `SELECT id, user_id, total_amount, status, created_at, updated_at FROM orders WHERE id = $1`
	var order OrderResponse
	err := s.db.Get(&order, query, req.Id)
	if err != nil {
		return nil, err
	}

	query = `SELECT product_id, quantity, price FROM order_items WHERE order_id = $1`
	var items []OrderItem
	err = s.db.Select(&items, query, req.Id)
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

func (s *server) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := s.db.Exec(query, req.Status, time.Now(), req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateOrderResponse{Success: true}, nil
}

func (s *server) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := s.db.Exec(query, req.Id)
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

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	db, err := sqlx.Connect("postgres", dbSource)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	pb.RegisterOrderServiceServer(s, &server{
		db:            db,
		userClient:    userClient,
		productClient: productClient,
	})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
