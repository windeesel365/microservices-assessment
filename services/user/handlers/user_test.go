package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	pb "github.com/windeesel365/microservices-assessment/services/user/userpb"
	"github.com/windeesel365/microservices-assessment/services/user/utils"
)

func TestServer_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	server := &Server{DB: sqlxDB}

	req := &pb.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
		Role:     "user",
	}

	passwordHash, _ := utils.HashPassword(req.Password)

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(req.Username, req.Email, passwordHash, req.Role).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, time.Now(), time.Now()))

	resp, err := server.CreateUser(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, req.Username, resp.Username)
	assert.Equal(t, req.Email, resp.Email)
	assert.Equal(t, req.Role, resp.Role)
}

func TestServer_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	server := &Server{DB: sqlxDB}

	req := &pb.GetUserRequest{Id: 1}

	mock.ExpectQuery(`SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = $1`).
		WithArgs(req.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "role", "created_at", "updated_at"}).
			AddRow(1, "testuser", "test@example.com", "user", time.Now(), time.Now()))

	resp, err := server.GetUser(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, "testuser", resp.Username)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.Equal(t, "user", resp.Role)
}

func TestServer_UpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	server := &Server{DB: sqlxDB}

	req := &pb.UpdateUserRequest{
		Id:       1,
		Username: "updateduser",
		Email:    "updated@example.com",
		Password: "newpassword",
	}

	passwordHash, _ := utils.HashPassword(req.Password)

	mock.ExpectQuery(`UPDATE users SET username = \$1, email = \$2, password_hash = \$3, updated_at = \$4 WHERE id = \$5 RETURNING id, role, created_at, updated_at`).
		WithArgs(req.Username, req.Email, passwordHash, sqlmock.AnyArg(), req.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role", "created_at", "updated_at"}).
			AddRow(1, "user", time.Now(), time.Now()))

	resp, err := server.UpdateUser(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, req.Username, resp.Username)
	assert.Equal(t, req.Email, resp.Email)
}

func TestServer_DeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	server := &Server{DB: sqlxDB}

	req := &pb.DeleteUserRequest{Id: 1}

	mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs(req.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	resp, err := server.DeleteUser(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
}
