package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eokwukwe/golearn/tasks/config"
	"github.com/eokwukwe/golearn/tasks/handlers"
	"github.com/eokwukwe/golearn/tasks/models"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	// Set up test database
	config.DB = config.InitTestDB()
	defer config.DB.Close()

	// Create a test user
	testUser := models.User{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "correct_password",
	}

	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testUser.Password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	_, err = config.DB.Exec(
		"INSERT INTO users (email, name, password) VALUES (?, ?, ?)",
		testUser.Email,
		testUser.Name,
		hashedPassword,
	)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test cases
	tests := []struct {
		name        string
		requestBody models.LoginRequest
		wantStatus  int
		wantMessage string
	}{
		{
			name: "valid login",
			requestBody: models.LoginRequest{
				Email:    testUser.Email,
				Password: testUser.Password,
			},
			wantStatus:  http.StatusOK,
			wantMessage: "Login successful",
		},
		{
			name: "invalid email",
			requestBody: models.LoginRequest{
				Email:    "invalid",
				Password: "password",
			},
			wantStatus:  http.StatusUnprocessableEntity,
			wantMessage: "Validation failed",
		},
		{
			name: "missing password",
			requestBody: models.LoginRequest{
				Email: "test@example.com",
			},
			wantStatus:  http.StatusUnprocessableEntity,
			wantMessage: "Validation failed",
		},
		{
			name: "wrong password",
			requestBody: models.LoginRequest{
				Email:    "test@example.com",
				Password: "wrong_password",
			},
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Invalid credentials",
		},
		{
			name: "non-existent user",
			requestBody: models.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password",
			},
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			reqBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rec := httptest.NewRecorder()

			// Call handler
			handlers.Login(rec, req)

			// Check status code
			if rec.Code != tt.wantStatus {
				t.Errorf("Login() status code = %v, want %v", rec.Code, tt.wantStatus)
			}

			// Check response body
			var response models.Response
			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Errorf("Failed to unmarshal response: %v", err)
			}

			if response.Status != "error" && response.Status != "success" {
				t.Errorf("Login() response status = %v, want 'error' or 'success'", response.Status)
			}

			if response.Message != tt.wantMessage {
				t.Errorf("Login() response message = %v, want %v", response.Message, tt.wantMessage)
			}
		})
	}
}
