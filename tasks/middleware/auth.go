package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/eokwukwe/golearn/tasks/config"
	"github.com/eokwukwe/golearn/tasks/models"
)

type contextKey string

const (
	ContextUserIDKey contextKey = "user_id"
	ContextUserKey   contextKey = "user"
)

// AuthMiddleware checks for valid token and adds user to context
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header and remove Bearer prefix
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.NewErrorResponse("Authorization token required", nil))
			return
		}

		// Extract token by removing "Bearer " prefix
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.NewErrorResponse("Authorization token required", nil))
			return
		}

		// Get user ID and expires_at from sessions table
		var userID int
		var expiresAt time.Time
		var err error
		if err = config.DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE token = ?", token).Scan(&userID, &expiresAt); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.NewErrorResponse("Invalid token", nil))
			return
		}

		// Check if token has expired
		if expiresAt.Before(time.Now()) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.NewErrorResponse("Token has expired", nil))
			return
		}

		// Get user from database
		var user models.User
		if err = config.DB.QueryRow("SELECT id, email, created_at FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Email, &user.CreatedAt); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.NewErrorResponse("User not found", nil))
			return
		}

		// Add user and user ID to context
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextUserIDKey, userID)
		ctx = context.WithValue(ctx, ContextUserKey, user)

		// Call next handler with updated context
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

// GetUserIDFromContext retrieves user ID from request context
func GetUserIDFromContext(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(ContextUserIDKey).(int)
	return userID, ok
}

// GetUserFromContext retrieves user from request context
func GetUserFromContext(r *http.Request) (*models.User, bool) {
	user, ok := r.Context().Value(ContextUserKey).(*models.User)
	return user, ok
}
