package calculations

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

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

func Factorial(n int) uint64 {
	if n == 0 {
		return 1
	}
	factorial := uint64(1)
	for i := 1; i <= n; i++ {
		factorial *= uint64(i)
	}
	return factorial
}

func CalculateHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var request CalculationRequest

	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check if 'a' and 'b' are non-negative integers
	if request.A < 0 || request.B < 0 {
		http.Error(w, "Incorrect input", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var aFactorial, bFactorial uint64

	go func() {
		defer wg.Done()
		aFactorial = Factorial(request.A)
	}()

	go func() {
		defer wg.Done()
		bFactorial = Factorial(request.B)
	}()

	wg.Wait()

	response := CalculationResponse{
		AFactorial: aFactorial,
		BFactorial: bFactorial,
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
