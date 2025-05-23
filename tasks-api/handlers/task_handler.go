package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/eokwukwe/golearn/tasks/config"
	"github.com/eokwukwe/golearn/tasks/middleware"
	"github.com/eokwukwe/golearn/tasks/models"
	"github.com/go-playground/validator/v10"
)

// GetTasks retrieves all tasks for the authenticated user
func GetTasks(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "user ID not found in context", fmt.Errorf("user id not found in context"))
		return
	}

	// Fetch tasks for the user
	rows, err := config.DB.Query("SELECT id, title, description, created_at, updated_at, completed FROM tasks WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch tasks", err)
		return
	}

	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.UpdatedAt, &task.Completed); err != nil {
			config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to scan task", err)
			return
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch tasks", err)
		return
	}

	config.WriteSuccessResponse(w, "Tasks retrieved successfully", tasks)
}

// GetOneTask retrieves a single task for the authenticated user
func GetOneTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "User ID not found in context", fmt.Errorf("user id not found in context"))
		return
	}

	// Get task ID from URL path
	path := r.URL.Path
	taskID := strings.TrimPrefix(path, "/api/v1/tasks/")
	if taskID == "" {
		config.WriteErrorResponse(w, http.StatusBadRequest, "Task ID is required", nil)
		return
	}

	// Fetch the task
	var task models.Task
	err := config.DB.QueryRow("SELECT id, title, description, created_at, updated_at, completed FROM tasks WHERE user_id = ? AND id = ?", userID, taskID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.Completed,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			config.WriteErrorResponse(w, http.StatusNotFound, "Task not found", nil)
			return
		}
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch task", err)
		return
	}

	// Return success response
	config.WriteSuccessResponse(w, "Task retrieved successfully", models.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Completed:   task.Completed,
	})
}

// DeleteTask deletes a single task for the authenticated user
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "User ID not found in context", fmt.Errorf("user id not found in context"))
		return
	}

	// Get task ID from URL path
	path := r.URL.Path
	taskID := strings.TrimPrefix(path, "/api/v1/tasks/")
	if taskID == "" {
		config.WriteErrorResponse(w, http.StatusBadRequest, "Task ID is required", nil)
		return
	}

	// Check if task exists
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE user_id = ? AND id = ?)", userID, taskID).Scan(&exists)
	if err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to check task existence", err)
		return
	}

	if !exists {
		config.WriteErrorResponse(w, http.StatusNotFound, "Task not found", nil)
		return
	}

	// Delete the task
	_, err = config.DB.Exec("DELETE FROM tasks WHERE user_id = ? AND id = ?", userID, taskID)
	if err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete task", err)
		return
	}

	// Return success response
	config.WriteSuccessResponse(w, "Task deleted successfully", nil)
}

// UpdateTask updates a task for the authenticated user
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "User ID not found in context", fmt.Errorf("user id not found in context"))
		return
	}

	// Get task ID from URL path
	path := r.URL.Path
	taskID := strings.TrimPrefix(path, "/api/v1/tasks/")
	if taskID == "" {
		config.WriteErrorResponse(w, http.StatusBadRequest, "Task ID is required", nil)
		return
	}

	// Check if task exists
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE user_id = ? AND id = ?)", userID, taskID).Scan(&exists)
	if err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to check task existence", err)
		return
	}

	if !exists {
		config.WriteErrorResponse(w, http.StatusNotFound, "Task not found", nil)
		return
	}

	// Check if request body is empty
	if r.ContentLength == 0 {
		config.WriteErrorResponse(w, http.StatusBadRequest, "Request body is required", nil)
		return
	}

	// Parse request body
	var req models.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		config.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		config.WriteErrorResponse(w, http.StatusUnprocessableEntity, "Validation failed", err)
		return
	}

	// Update task
	_, err = config.DB.Exec(
		"UPDATE tasks SET title = ?, description = ? WHERE user_id = ? AND id = ?",
		req.Title,
		req.Description,
		userID,
		taskID,
	)
	if err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update task", err)
		return
	}

	// Fetch the created task
	var task models.Task
	if err := config.DB.QueryRow("SELECT id, title, description, created_at, updated_at, completed FROM tasks WHERE id = ?", taskID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.Completed,
	); err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch created task", err)
		return
	}

	// Return success response with the complete task
	config.WriteSuccessResponse(w, "Task updated successfully", models.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Completed:   task.Completed,
	})
}

// CompleteTask marks a task as completed for the authenticated user
func CompleteTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Get task ID from URL path
	path := r.URL.Path
	taskID := strings.TrimPrefix(path, "/api/v1/tasks/")
	if taskID == "" {
		config.WriteErrorResponse(w, http.StatusBadRequest, "Task ID is required", nil)
		return
	}

	// Check if task exists
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE user_id = ? AND id = ?)", userID, taskID).Scan(&exists)
	if err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to check task existence", err)
		return
	}

	if !exists {
		config.WriteErrorResponse(w, http.StatusNotFound, "Task not found", nil)
		return
	}

	// Mark task as completed
	_, err = config.DB.Exec(
		"UPDATE tasks SET completed = ? WHERE user_id = ? AND id = ?",
		true,
		userID,
		taskID,
	)
	if err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to mark task as completed", err)
		return
	}

	config.WriteSuccessResponse(w, "Task marked as completed successfully", nil)
}

// CreateTask creates a new task for the authenticated user
func CreateTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Check if request body is empty
	if r.ContentLength == 0 {
		config.WriteErrorResponse(w, http.StatusBadRequest, "Request body is required", nil)
		return
	}

	// Parse request body
	var req models.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		config.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		config.WriteErrorResponse(w, http.StatusUnprocessableEntity, "Validation failed", err)
		return
	}

	// Insert task
	result, err := config.DB.Exec(
		"INSERT INTO tasks (user_id, title, description) VALUES (?, ?, ?)",
		userID,
		req.Title,
		req.Description,
	)
	if err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create task", err)
		return
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get task ID", err)
		return
	}

	// Fetch the created task
	var task models.Task
	if err := config.DB.QueryRow("SELECT id, title, description, created_at, updated_at, completed FROM tasks WHERE id = ?", lastID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.Completed,
	); err != nil {
		config.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch created task", err)
		return
	}

	// Return success response with the complete task
	config.WriteCreatedResponse(w, "Task created successfully", models.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Completed:   task.Completed,
	})
}
