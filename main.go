package main

import (
	"calculationFactorial/calculations"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	router := httprouter.New()

	router.GET("/", index)
	router.POST("/calculate", calculations.CalculateHandler)

	fmt.Println("Server listening on port 8989")
	http.ListenAndServe(":8989", router)
}
