package calculations

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"calculationFactorial/calculations"

	"github.com/julienschmidt/httprouter"
)

func TestFactorial(t *testing.T) {
	tests := []struct {
		input  int
		output uint64
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 6},
		{5, 120},
		{10, 3628800},
	}

	for _, test := range tests {
		result := calculations.Factorial(test.input)
		if result != test.output {
			t.Errorf("Factorial of %d returned %d, expected %d", test.input, result, test.output)
		}
	}
}

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		requestBody  string
		expectedCode int
		expectedBody string
	}{
		{
			`{"a": 5, "b": 3}`,
			http.StatusOK,
			`{"a_factorial":120,"b_factorial":6}`,
		},
		{
			`{"a": 10, "b": 0}`,
			http.StatusOK,
			`{"a_factorial":3628800,"b_factorial":1}`,
		},
		{
			`{"a": 0, "b": 5}`,
			http.StatusOK,
			`{"a_factorial":1,"b_factorial":120}`,
		},
		{
			`{"a": 1, "b": 1}`,
			http.StatusOK,
			`{"a_factorial":1,"b_factorial":1}`,
		},
		{
			`{"a": -1, "b": 5}`,
			http.StatusBadRequest,
			`Incorrect input`,
		},
		{
			`{"a": 5, "b": -1}`,
			http.StatusBadRequest,
			`Incorrect input`,
		},
	}

	router := httprouter.New()
	router.POST("/calculate", calculations.CalculateHandler)

	for _, test := range tests {
		req, err := http.NewRequest("POST", "/calculate", strings.NewReader(test.requestBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != test.expectedCode {
			t.Errorf("Expected status code %d, but got %d", test.expectedCode, rr.Code)
		}

		gotBody := strings.TrimRight(rr.Body.String(), "\n")
		if gotBody != test.expectedBody {
			t.Errorf("Expected response body '%s', but got '%s'", test.expectedBody, gotBody)
		}
	}
}
