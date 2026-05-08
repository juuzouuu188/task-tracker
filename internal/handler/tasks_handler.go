package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/juuzouuu188/task-tracker/internal/model"
)

var tasks = map[string]model.Task{}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	json.NewDecoder(r.Body).Decode(&input)

	task := model.Task{
		ID:     generateID(),
		Title:  input.Title,
		Status: model.StatusTodo,
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
