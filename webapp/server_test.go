package main

import (
	"WebApp/store"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MockUUIDGenerator implements store.UuidGenerator.
type MockUUIDGenerator struct{}

// GenerateUUID mocks the UUID generation.
func (g *MockUUIDGenerator) GenerateUUID() string {
	return "1"
}

// TestGetHandler tests the GetHandler function with InMemoryStore.
func TestGetHandler(t *testing.T) {
	// Create the context for running the store.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	TemplatePath = "../static/index.html"
	// Initialize the InMemoryStore.
	mockStore := store.NewInMemoryStore(ctx)

	// Insert mock data into the store.
	mockStore.ToDo = append(mockStore.ToDo, store.ToDoItem{Item: "Task 1", Status: "Pending"})
	mockStore.ToDo = append(mockStore.ToDo, store.ToDoItem{Item: "Task 2", Status: "Pending"})
	mockStore.ToDo = append(mockStore.ToDo, store.ToDoItem{Item: "Task 3", Status: "Pending"})

	// Mock the UUID generator.
	mockStore.SetUUIDGenerator(&MockUUIDGenerator{})

	// Override the global Store variable with the mock store.
	Store = mockStore

	// Create a mock HTTP request (GET request to "/").
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response.
	rr := httptest.NewRecorder()

	// Create a context with the taskID for the request.
	ctx = context.WithValue(req.Context(), taskID("taskID"), "test-trace-id")
	req = req.WithContext(ctx)

	// Call the handler function with the response recorder and the request.
	GetHandler(rr, req)

	// Check if the status code is what we expect (200 OK).
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedHTML := rr.Body.String()
	if !strings.Contains(expectedHTML, "Task 1 - Pending") && strings.Contains(expectedHTML, "Task 2 - Pending") {
		t.Errorf("handler returned unexpected body: got %q \n \n want %q", rr.Body.String(), expectedHTML)
	}
}
