package handlers

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/windeesel365/microservices-assessment/services/user/models"
	pb "github.com/windeesel365/microservices-assessment/services/user/userpb"
	"github.com/windeesel365/microservices-assessment/services/user/utils"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	DB *sqlx.DB
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	var user models.UserResponse
	err = s.DB.QueryRowx(query, req.Username, req.Email, passwordHash, req.Role).StructScan(&user)
	if err != nil {
		return nil, err
	}

	user.Username = req.Username
	user.Email = req.Email
	user.Role = req.Role

	return &pb.UserResponse{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	query := `SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = $1`
	var user models.UserResponse
	err := s.DB.Get(&user, query, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	query := `UPDATE users SET username = $1, email = $2, password_hash = $3, updated_at = $4 WHERE id = $5 RETURNING id, role, created_at, updated_at`
	var user models.UserResponse
	err = s.DB.QueryRowx(query, req.Username, req.Email, passwordHash, time.Now(), req.Id).StructScan(&user)
	if err != nil {
		return nil, err
	}

	user.Username = req.Username
	user.Email = req.Email

	return &pb.UserResponse{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	query := `DELETE FROM users WHERE id = $1`
	result, err := s.DB.Exec(query, req.Id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return &pb.DeleteUserResponse{Success: false}, nil
	}

	return &pb.DeleteUserResponse{Success: true}, nil
}
