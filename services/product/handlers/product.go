package handlers

import (
	"context"

	"github.com/jmoiron/sqlx"
	pb "github.com/windeesel365/microservices-assessment/services/product/productpb"
)

type Server struct {
	pb.UnimplementedProductServiceServer
	DB *sqlx.DB
}

func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	var id int
	query := "INSERT INTO products (name, description, price, category_id, stock) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := s.DB.QueryRowx(query, req.GetName(), req.GetDescription(), req.GetPrice(), req.GetCategoryId(), req.GetStock()).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &pb.CreateProductResponse{Id: int32(id)}, nil
}

func (s *Server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	var product pb.GetProductResponse
	query := "SELECT id, name, description, price, category_id, stock, created_at, updated_at FROM products WHERE id = $1"
	err := s.DB.Get(&product, query, req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{
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

func (s *Server) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	query := "UPDATE products SET name = $1, description = $2, price = $3, category_id = $4, stock = $5 WHERE id = $6"
	_, err := s.DB.Exec(query, req.GetName(), req.GetDescription(), req.GetPrice(), req.GetCategoryId(), req.GetStock(), req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.UpdateProductResponse{Success: true}, nil
}

func (s *Server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	query := "DELETE FROM products WHERE id = $1"
	_, err := s.DB.Exec(query, req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteProductResponse{Success: true}, nil
}
