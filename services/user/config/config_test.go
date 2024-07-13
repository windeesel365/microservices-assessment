package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name          string
		port          string
		databaseURL   string
		expectError   bool
		expectedPort  string
		expectedDBURL string
	}{
		{
			name:          "Valid environment variables",
			port:          "8080",
			databaseURL:   "postgres://user:pass@localhost:5432/dbname",
			expectError:   false,
			expectedPort:  "8080",
			expectedDBURL: "postgres://user:pass@localhost:5432/dbname",
		},
		{
			name:        "Missing PORT",
			port:        "",
			databaseURL: "postgres://user:pass@localhost:5432/dbname",
			expectError: true,
		},
		{
			name:        "Missing DATABASE_URL",
			port:        "8080",
			databaseURL: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("PORT", tt.port)
			os.Setenv("DATABASE_URL", tt.databaseURL)

			defer func() {
				_ = os.Unsetenv("PORT")
				_ = os.Unsetenv("DATABASE_URL")
			}()

			if tt.expectError {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected error but got none")
					}
				}()
			}

			config := LoadConfig()

			if !tt.expectError {
				if config.Port != tt.expectedPort {
					t.Errorf("Expected port %v, got %v", tt.expectedPort, config.Port)
				}
				if config.DatabaseURL != tt.expectedDBURL {
					t.Errorf("Expected database URL %v, got %v", tt.expectedDBURL, config.DatabaseURL)
				}
			}
		})
	}
}
