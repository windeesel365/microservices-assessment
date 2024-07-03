package main

import (
	"context"
	"log"
	pb "microecommerce/pb/productpb"
	"net"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedProductServiceServer
	db *sqlx.DB
}

type ProductResponse struct {
	Id          int32   `db:"id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Price       float32 `db:"price"`
	CategoryId  int32   `db:"category_id"`
	Stock       int32   `db:"stock"`
	CreatedAt   string  `db:"created_at"`
	UpdatedAt   string  `db:"updated_at"`
}

func (s *server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	var id int
	query := "INSERT INTO products (name, description, price, category_id, stock) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := s.db.QueryRowx(query, req.GetName(), req.GetDescription(), req.GetPrice(), req.GetCategoryId(), req.GetStock()).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &pb.CreateProductResponse{Id: int32(id)}, nil
}

func (s *server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	var product pb.GetProductResponse
	query := "SELECT id, name, description, price, category_id, stock, created_at, updated_at FROM products WHERE id = $1"
	err := s.db.Get(&product, query, req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.ProductResponse{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CategoryId:  product.CategoryId,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

func (s *server) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	query := "UPDATE products SET name = $1, description = $2, price = $3, category_id = $4, stock = $5 WHERE id = $6"
	_, err := s.db.Exec(query, req.GetName(), req.GetDescription(), req.GetPrice(), req.GetCategoryId(), req.GetStock(), req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.UpdateProductResponse{Success: true}, nil
}

func (s *server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	query := "DELETE FROM products WHERE id = $1"
	_, err := s.db.Exec(query, req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteProductResponse{Success: true}, nil
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

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	db, err := sqlx.Connect("postgres", dbSource)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	pb.RegisterUserServiceServer(s, &server{db: db})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
