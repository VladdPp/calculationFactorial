package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	testCases := []struct {
		desc     string
		request  CalculationRequest
		expected CalculationResponse
	}{
		{
			desc:     "Test case 1",
			request:  CalculationRequest{A: 5, B: 3},
			expected: CalculationResponse{AFactorial: 120, BFactorial: 6},
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// Encode the request body
			body, err := json.Marshal(tc.request)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			// Create a request with the encoded body
			req, err := http.NewRequest("POST", "/calculate", bytes.NewReader(body))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Call the handler function directly and pass the request
			validateInput(calculateHandler)(rr, req, nil)

			// Check the response status code
			if rr.Code != http.StatusOK {
				t.Errorf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
			}

			// Decode the response body into CalculationResponse
			var response CalculationResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}

			// Check if the actual response matches the expected response
			if response != tc.expected {
				t.Errorf("Expected response %+v, but got %+v", tc.expected, response)
			}
		})
	}
}
