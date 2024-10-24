package main

import (
	"WebApp/store"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
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

	// Check if the response contains the expected output (HTML content).
	expectedHTML := "<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Todo App</title>\r\n</head>\r\n<body>\r\n    <h1>Todo List</h1>\r\n    \r\n    \r\n    <form method=\"POST\" action=\"/add\">\r\n        <input type=\"text\" name=\"item\" placeholder=\"Todo item\" required>\r\n        <input type=\"text\" name=\"status\" placeholder=\"Status\" required>\r\n        <button type=\"submit\">Add</button>\r\n    </form>\r\n\r\n    <h2>Current Todo Items</h2>\r\n    <ul>\r\n        \r\n        <li>\r\n            \r\n            Task 1 - Pending\r\n\r\n            \r\n            <form method=\"POST\" action=\"/delete\" style=\"display:inline;\">\r\n                <input type=\"hidden\" name=\"item\" value=\"Task 1\">\r\n                <button type=\"submit\">Delete</button>\r\n            </form>\r\n            \r\n            <form method=\"POST\" action=\"/update\" style=\"display:inline;\">\r\n                <input type=\"hidden\" name=\"item\" value=\"Task 1\">\r\n                <input type=\"text\" name=\"status\" placeholder=\"New status\" required>\r\n                <button type=\"submit\">Update</button>\r\n            </form>\r\n        </li>\r\n        \r\n    </ul>\r\n</body>\r\n</html>\r\n"
	fmt.Println("\n", rr.Body.String())
	if rr.Body.String() != expectedHTML {
		t.Errorf("handler returned unexpected body: got %q \n \n want %q", rr.Body.String(), expectedHTML)
	}
}
