package main

import (
	"context"
	pb "microecommerce/pb/productpb"

	"github.com/jmoiron/sqlx"
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
