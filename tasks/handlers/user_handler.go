package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eokwukwe/golearn/tasks/config"
	"github.com/eokwukwe/golearn/tasks/models"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	// Check that method is post
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Method not allowed. Use POST", nil))
		return
	}

	// Check if request body is empty
	if r.ContentLength == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Request body is empty", nil))
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Invalid request body", err))
		return
	}

	// Validate input using validator
	err := validateUser(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Validation failed", err))
		return
	}

	// Check if email already exists
	var count int
	if err := config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.Email).Scan(&count); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to check email", err))
		return
	}

	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Email already exists", nil))
		return
	}

	// Hash password securely
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to hash password", err))
		return
	}

	// Create user
	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	// Insert user into database and get result
	result, err := config.DB.Exec(`INSERT INTO users (name, email, password) VALUES (?, ?, ?)`,
		user.Name, user.Email, user.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to create user", nil))
		return
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to get user ID", nil))
		return
	}

	// Prepare response
	response := models.RegisterResponse{
		ID:        int(lastID),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	// Write success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(models.NewSuccessResponse("User created successfully", response)); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to encode response", nil))
	}
}

// hashPassword securely hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// validateUser validates user registration input
func validateUser(user *models.RegisterRequest) error {
	validate := validator.New()
	return validate.Struct(user)
}
