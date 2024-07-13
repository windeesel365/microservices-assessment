package utils

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Fatalf("expected hashed password to be non-empty")
	}

	//check hashed password matched original password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		t.Fatalf("expected password to match, got %v", err)
	}
}
