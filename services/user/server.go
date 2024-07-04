package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	pb "ecommerce-user/userpb"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	db *sqlx.DB
}

type UserResponse struct {
	Id        int32  `db:"id"`
	Username  string `db:"username"`
	Email     string `db:"email"`
	Role      string `db:"role"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	var user UserResponse
	err = s.db.QueryRowx(query, req.Username, req.Email, passwordHash, req.Role).StructScan(&user)
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

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	query := `SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = $1`
	var user UserResponse
	err := s.db.Get(&user, query, req.Id)
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

func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	query := `UPDATE users SET username = $1, email = $2, password_hash = $3, updated_at = $4 WHERE id = $5 RETURNING id, role, created_at, updated_at`
	var user UserResponse
	err = s.db.QueryRowx(query, req.Username, req.Email, passwordHash, time.Now(), req.Id).StructScan(&user)
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

func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	query := `DELETE FROM users WHERE id = $1`
	result, err := s.db.Exec(query, req.Id)
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

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
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
