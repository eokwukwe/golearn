package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/eokwukwe/golearn/tasks/config"
	"github.com/eokwukwe/golearn/tasks/handlers"
	"github.com/eokwukwe/golearn/tasks/models"
)

func TestRegister(t *testing.T) {
	// Initialize test database
	if err := config.InitDB(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	// Test cases
	type testCase struct {
		name         string
		request      models.RegisterRequest
		wantStatus   int
		wantResponse models.Response
	}

	testCases := []testCase{
		{
			name: "valid registration",
			request: models.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "securepassword123",
			},
			wantStatus: http.StatusCreated,
			wantResponse: models.Response{
				Status:  "success",
				Message: "User created successfully",
			},
		},
		{
			name: "invalid email",
			request: models.RegisterRequest{
				Name:     "Jane Doe",
				Email:    "invalid-email",
				Password: "securepassword123",
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantResponse: models.Response{
				Status:  "error",
				Message: "Validation failed",
				Errors: map[string]string{
					"email": "email must be a valid email address",
				},
			},
		},
		{
			name: "duplicate email",
			request: models.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "securepassword123",
			},
			wantStatus: http.StatusConflict,
			wantResponse: models.Response{
				Status:  "error",
				Message: "Email already exists",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request body
			reqBody, err := json.Marshal(tc.request)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rec := httptest.NewRecorder()

			// Call the handler
			handlers.Register(rec, req)

			// Check status code
			if rec.Code != tc.wantStatus {
				t.Errorf("%s: expected status %d, got %d", tc.name, tc.wantStatus, rec.Code)
			}

			// Check response body
			var resp models.Response
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Errorf("%s: failed to decode response: %v", tc.name, err)
				return
			}

			// Check response fields
			if resp.Status != tc.wantResponse.Status {
				t.Errorf("%s: expected status %q, got %q", tc.name, tc.wantResponse.Status, resp.Status)
			}
			if resp.Message != tc.wantResponse.Message {
				t.Errorf("%s: expected message %q, got %q", tc.name, tc.wantResponse.Message, resp.Message)
			}
			if tc.wantResponse.Errors != nil {
				if resp.Errors == nil {
					t.Errorf("%s: expected errors, got nil", tc.name)
				} else if !reflect.DeepEqual(resp.Errors, tc.wantResponse.Errors) {
					t.Errorf("%s: expected errors %v, got %v", tc.name, tc.wantResponse.Errors, resp.Errors)
				}
			}
		})
	}

	// Clean up test database
	if _, err := config.DB.Exec("DELETE FROM users"); err != nil {
		t.Errorf("Failed to clean up test database: %v", err)
	}
}
