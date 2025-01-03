package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type RequestBody struct {
	Message string `json:"task"`
	IsDone  bool   `json:"is_done"`
}

func getAllMessages() ([]Message, error) {
	var messages []Message
	result := DB.Find(&messages)

	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}

func GetLastTaskHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := getAllMessages()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(messages) == 0 {
		fmt.Fprintln(w, "No tasks found")
		return
	}

	lastTask := messages[len(messages)-1]

	fmt.Fprintf(w, "Last task name: %s, isDone: %v\n", lastTask.Task, lastTask.IsDone)
	fmt.Fprintln(w, "============")
	fmt.Fprintln(w, "total tasks: ", len(messages))
}

func GetAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := getAllMessages()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(messages) == 0 {
		fmt.Fprintln(w, "No tasks found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestBody
	json.NewDecoder(r.Body).Decode(&reqBody)

	if reqBody.Message == "" {
		http.Error(w, "Empty task name", http.StatusBadRequest)

		return
	}

	DB.Create(&Message{Task: reqBody.Message, IsDone: reqBody.IsDone})
}

func main() {
	InitDB()

	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	router.HandleFunc("/api/todo/last", GetLastTaskHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/todo", GetAllTaskHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/todo", UpdateTaskHandler).Methods(http.MethodPost)

	http.ListenAndServe(":8080", router)
}
