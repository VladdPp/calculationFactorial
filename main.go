package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CalculationRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

type CalculationResponse struct {
	AFactorial uint64 `json:"a_factorial"`
	BFactorial uint64 `json:"b_factorial"`
}

type contextKey string

const calculationRequestKey contextKey = "calculationRequest"

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func calculateFactorial(n int) uint64 {
	if n == 0 {
		return 1
	}
	factorial := uint64(1)
	for i := 1; i <= n; i++ {
		factorial *= uint64(i)
	}
	return factorial
}

func calculateHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	request, ok := r.Context().Value(calculationRequestKey).(*CalculationRequest)
	if !ok {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response := CalculationResponse{
		AFactorial: calculateFactorial(request.A),
		BFactorial: calculateFactorial(request.B),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func validateInput(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var request CalculationRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Check if 'a' and 'b' are non-negative integers
		if request.A < 0 || request.B < 0 {
			http.Error(w, "Incorrect input", http.StatusBadRequest)
			return
		}

		// Copy the decoded request into the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, calculationRequestKey, &request)
		r = r.WithContext(ctx)

		next(w, r, ps)
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", index)
	router.POST("/calculate", validateInput(calculateHandler))

	fmt.Println("Server listening on port 8989")
	http.ListenAndServe(":8989", router)
}
