package handlers

import (
	"database/sql"
	"encoding/json"
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
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Fetch tasks for the user
	rows, err := config.DB.Query("SELECT id, title, description, created_at, updated_at, completed FROM tasks WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to fetch tasks", err))
		return
	}
	defer rows.Close()

	// Parse and collect tasks
	var tasks []models.TaskResponse
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.UpdatedAt, &task.Completed); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to scan task", err))
			return
		}

		tasks = append(tasks, models.TaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			Completed:   task.Completed,
		})
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.NewSuccessResponse("Tasks retrieved successfully", tasks))
}

// GetOneTask retrieves a single task for the authenticated user
func GetOneTask(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Task ID is required", nil))
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.NewErrorResponse("Task not found", nil))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to fetch task", err))
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.NewSuccessResponse("Task retrieved successfully", models.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Completed:   task.Completed,
	}))
}

// DeleteTask deletes a single task for the authenticated user
func DeleteTask(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Task ID is required", nil))
		return
	}

	// Check if task exists
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE user_id = ? AND id = ?)", userID, taskID).Scan(&exists)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to check task existence", err))
		return
	}

	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Task not found", nil))
		return
	}

	// Delete the task
	_, err = config.DB.Exec("DELETE FROM tasks WHERE user_id = ? AND id = ?", userID, taskID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to delete task", err))
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.NewSuccessResponse("Task deleted successfully", nil))
}

// UpdateTask updates a task for the authenticated user
func UpdateTask(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Task ID is required", nil))
		return
	}

	// Check if task exists
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE user_id = ? AND id = ?)", userID, taskID).Scan(&exists)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to check task existence", err))
		return
	}

	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Task not found", nil))
		return
	}

	// Check if request body is empty
	if r.ContentLength == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Request body is required", nil))
		return
	}

	// Parse request body
	var req models.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Invalid request body", err))
		return
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Validation failed", err))
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to update task", err))
		return
	}

	// // Get the updated task ID
	// lastID, err := result.LastInsertId()
	// if err != nil {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to get task ID", err))
	// 	return
	// }

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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to fetch created task", err))
		return
	}

	// Return success response with the complete task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.NewSuccessResponse("Task updated successfully", models.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Completed:   task.Completed,
	}))
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Task ID is required", nil))
		return
	}

	// Check if task exists
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE user_id = ? AND id = ?)", userID, taskID).Scan(&exists)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to check task existence", err))
		return
	}

	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Task not found", nil))
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to mark task as completed", err))
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.NewSuccessResponse("Task marked as completed successfully", nil))
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Request body is required", nil))
		return
	}

	// Parse request body
	var req models.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Invalid request body", err))
		return
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Validation failed", err))
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to create task", err))
		return
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to get task ID", err))
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse("Failed to fetch created task", err))
		return
	}

	// Return success response with the complete task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.NewSuccessResponse("Task created successfully", models.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Completed:   task.Completed,
	}))
}
