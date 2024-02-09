package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Tests if everything runs accordingly for HTTP request in the API
func TestAssessTransactions_ValidInput(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()

	// Register the handler function for the /check_transactions route
	router.POST("/check_transactions", AssessTransactions)

	// Define a test case
	tests := []struct {
		name     string
		payload  string // JSON sent in the request
		expected int    // Expected HTTP status code
	}{
		{
			name: "Valid Post",
			payload: `{
				"transactions": [
				  {"id": 1, "user_id": 1, "amount_us_cents": 200000, "card_id": 1},
				  {"id": 2, "user_id": 1, "amount_us_cents": 600000, "card_id": 1},
				  {"id": 3, "user_id": 1, "amount_us_cents": 1100000, "card_id": 1},
				  {"id": 4, "user_id": 2, "amount_us_cents": 100000, "card_id": 2},
				  {"id": 5, "user_id": 2, "amount_us_cents": 100000, "card_id": 3},
				  {"id": 6, "user_id": 2, "amount_us_cents": 100000, "card_id": 4}
				]
			  }`,
			expected: http.StatusOK,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP request with the payload
			req, err := http.NewRequest("POST", "/check_transactions", bytes.NewBufferString(test.payload))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder to capture the response
			recorder := httptest.NewRecorder()

			// Serve the HTTP request to the router
			router.ServeHTTP(recorder, req)

			// Check if the response status code matches the expected status code
			assert.Equal(t, test.expected, recorder.Code, "status code mismatch")
		})
	}
}

func TestAssessTransactions_InvalidInput(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()

	// Register the handler function for the /check_transactions route
	router.POST("/check_transactions", AssessTransactions)

	// Define a test case
	tests := []struct {
		name     string
		payload  string // JSON sent in the request
		expected int    // Expected HTTP status code
	}{
		{
			name: "Invalid Post",
			payload: `
				"transactions": [
				  {"id": 1, "id": 1},
				]
			  }`,
			expected: http.StatusBadRequest,
		},
		{
			name:     "Empty Post",
			payload:  ``,
			expected: http.StatusBadRequest,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP request with the payload
			req, err := http.NewRequest("POST", "/check_transactions", bytes.NewBufferString(test.payload))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder to capture the response
			recorder := httptest.NewRecorder()

			// Serve the HTTP request to the router
			router.ServeHTTP(recorder, req)

			// Check if the response status code matches the expected status code
			assert.Equal(t, test.expected, recorder.Code, "status code mismatch")
		})
	}
}

func TestMain(t *testing.T) {
	// Run the main function in a separate goroutine
	go main()

	// Sleep for a short time to allow the server to start
	time.Sleep(100 * time.Millisecond)

	// Send a request to the server to check if it's running
	resp, err := http.Get("http://localhost:9090")

	// Check for errors
	if err != nil {
		t.Errorf("Failed to send request to server: %v", err)
	}

	defer resp.Body.Close()
	// Check if the server responds with a success status code
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}
