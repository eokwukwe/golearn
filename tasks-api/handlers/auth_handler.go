package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/eokwukwe/golearn/tasks/config"
	"github.com/eokwukwe/golearn/tasks/models"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

const (
	SessionDuration = 24 * time.Hour // 24 hours
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Check that method is post
	if r.Method != http.MethodPost {
		config.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed. Use POST", nil)
		return
	}

	// Check if request body is empty
	if r.ContentLength == 0 {
		config.WriteErrorResponse(w, http.StatusBadRequest, "Request body is empty", nil)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		config.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate input using validator
	err := validateUserLogin(&req)
	if err != nil {
		config.WriteErrorResponse(w, http.StatusUnprocessableEntity, "Validation failed", err)
		return
	}

	// Find user by email
	var user models.User
	err = config.DB.QueryRow("SELECT id, name, email, password, created_at FROM users WHERE email = ?", req.Email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		config.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		config.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	// Generate session token
	token, err := generateToken()
	if err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	// Create new session
	session := models.Session{
		UserID:    user.ID,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(SessionDuration),
	}

	// Insert session into database
	result, err := config.DB.Exec(
		"INSERT INTO sessions (user_id, token, created_at, expires_at) VALUES (?, ?, ?, ?)",
		session.UserID,
		session.Token,
		session.CreatedAt,
		session.ExpiresAt,
	)
	if err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create session", err)
		return
	}

	// Get session ID
	_, err = result.LastInsertId()
	if err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get session ID", err)
		return
	}

	// Prepare response
	response := models.LoginResponse{
		Token: token,
		User: struct {
			ID        int       `json:"id"`
			Name      string    `json:"name"`
			Email     string    `json:"email"`
			CreatedAt time.Time `json:"created_at"`
		}{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	// Set success response
	config.WriteSuccessResponse(w, "Login successful", response)
}

func validateUserLogin(user *models.LoginRequest) error {
	validate := validator.New()
	return validate.Struct(user)
}

func generateToken() (string, error) {
	// Generate 32 random bytes (256 bits)
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Encode to base64 to get a URL-safe string
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
