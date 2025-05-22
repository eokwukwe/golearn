package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eokwukwe/golearn/tasks/config"
	"github.com/eokwukwe/golearn/tasks/handlers"
	"github.com/eokwukwe/golearn/tasks/middleware"
	"github.com/eokwukwe/golearn/tasks/models"
	"github.com/stretchr/testify/assert"
)

func setupTest() (*httptest.ResponseRecorder, *http.Request) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/tasks", nil)
	req.Header.Set("Authorization", "Bearer test-session-token")
	req = req.WithContext(context.WithValue(req.Context(), middleware.ContextUserIDKey, 1))
	return recorder, req
}

// setupTestData creates test data for the test user
func setupTestData() {
	// Create test user
	config.DB.Exec(`INSERT INTO users (id, email, password_hash) VALUES (?, ?, ?)`, 1, "test@example.com", "hashed-password")
	// Create test session
	config.DB.Exec(`INSERT INTO sessions (user_id, token, created_at, expires_at) VALUES (?, ?, datetime('now'), datetime('now', '+1 day'))`, 1, "test-session-token")
	// Create test task
	config.DB.Exec(`INSERT INTO tasks (id, user_id, title, description, created_at, updated_at, completed) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'), ?)`, 1, 1, "Test Task", "Test Description", false)
}

// cleanupTestData removes test data
func cleanupTestData() {
	config.DB.Exec("DELETE FROM sessions WHERE user_id = ?", 1)
	config.DB.Exec("DELETE FROM tasks WHERE user_id = ?", 1)
	config.DB.Exec("DELETE FROM users WHERE id = ?", 1)
	config.DB.Close()
}

func setupTestWithBody(body []byte) (*httptest.ResponseRecorder, *http.Request) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-session-token")
	req = req.WithContext(context.WithValue(req.Context(), middleware.ContextUserIDKey, 1))
	return recorder, req
}

func TestGetTasks(t *testing.T) {
	// Set up test database
	config.DB = config.InitTestDB()

	// Set up test data
	setupTestData()
	defer cleanupTestData()

	recorder, req := setupTest()
	handlers.GetTasks(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	var response models.Response
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))

	assert.NotNil(t, response.Data)
	assert.NotEmpty(t, response.Data)
}

func TestGetOneTask(t *testing.T) {
	// Set up test database
	config.DB = config.InitTestDB()

	// Set up test data
	setupTestData()
	defer cleanupTestData()

	// Test with valid task ID
	recorder, req := setupTest()
	req.URL.Path = "/api/v1/tasks/1"
	handlers.GetOneTask(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	var response models.Response
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.NotNil(t, response.Data)

	// Test with invalid task ID
	recorder, req = setupTest()
	req.URL.Path = "/api/v1/tasks/invalid"
	handlers.GetOneTask(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	var errorResponse models.Response
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &errorResponse))
	assert.NotNil(t, errorResponse.Message, "Task not found")
}

func TestDeleteTask(t *testing.T) {
	// Set up test database
	config.DB = config.InitTestDB()

	// Set up test data
	setupTestData()
	defer cleanupTestData()

	// Test with valid task ID
	recorder, req := setupTest()
	req.URL.Path = "/api/v1/tasks/1"
	req.Method = http.MethodDelete
	handlers.DeleteTask(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	var response models.Response
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.NotNil(t, response.Message, "Task deleted successfully")

	// Test with invalid task ID
	recorder, req = setupTest()
	req.URL.Path = "/api/v1/tasks/invalid"
	req.Method = http.MethodDelete
	handlers.DeleteTask(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	var errorResponse models.Response
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &errorResponse))
	assert.NotNil(t, errorResponse.Message, "Task not found")
}

func TestUpdateTask(t *testing.T) {
	// Set up test database
	config.DB = config.InitTestDB()

	// Set up test data
	setupTestData()
	defer cleanupTestData()

	// Test successful update
	task := models.Task{
		Title:       "Updated Task",
		Description: "Updated description",
	}
	body, err := json.Marshal(task)
	assert.NoError(t, err)

	recorder, req := setupTestWithBody(body)
	req.URL.Path = "/api/v1/tasks/1"
	req.Method = http.MethodPut
	handlers.UpdateTask(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	var response models.Response
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.NotNil(t, response.Data)
	assert.Equal(t, "Task updated successfully", response.Message)

	// Test with invalid task ID
	recorder, req = setupTestWithBody(body)
	req.URL.Path = "/api/v1/tasks/invalid"
	req.Method = http.MethodPut
	handlers.UpdateTask(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	var errorResponse models.Response
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &errorResponse))
	assert.NotNil(t, errorResponse.Message, "Task not found")
}

func TestCompleteTask(t *testing.T) {
	// Set up test database
	config.DB = config.InitTestDB()

	// Set up test data
	setupTestData()
	defer cleanupTestData()

	// Test successful completion
	recorder, req := setupTest()
	req.URL.Path = "/api/v1/tasks/1"
	req.Method = http.MethodPatch
	handlers.CompleteTask(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	var response models.Response
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.NotNil(t, response.Message)

	// Test with invalid task ID
	recorder, req = setupTest()
	req.URL.Path = "/api/v1/tasks/invalid"
	req.Method = http.MethodPatch
	handlers.CompleteTask(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	var errorResponse models.Response
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &errorResponse))
	assert.NotNil(t, errorResponse.Message, "Task not found")
}

func TestCreateTask(t *testing.T) {

	// Set up test database
	config.DB = config.InitTestDB()

	// Set up test data
	setupTestData()
	defer cleanupTestData()

	// Test successful creation
	task := models.Task{
		Title:       "Test Task",
		Description: "Test description",
	}
	body, err := json.Marshal(task)
	assert.NoError(t, err)

	recorder, req := setupTestWithBody(body)
	req.Method = http.MethodPost
	handlers.CreateTask(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	var response models.Response
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.NotNil(t, response.Data)
	assert.Equal(t, "Task created successfully", response.Message)

	// Test with invalid request body
	recorder, req = setupTestWithBody([]byte("invalid json"))
	req.Method = http.MethodPost
	handlers.CreateTask(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	var errorResponse models.Response
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &errorResponse))
	assert.NotNil(t, errorResponse.Message)
}
