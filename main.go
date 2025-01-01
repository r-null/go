package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type RequestBody struct {
	Message string `json:"message"`
}

var requestBody = RequestBody{}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello", requestBody.Message)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	json.NewDecoder(r.Body).Decode(&requestBody)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/hello", HelloHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/hello", UpdateTaskHandler).Methods(http.MethodPost)

	http.ListenAndServe(":8080", router)
}
