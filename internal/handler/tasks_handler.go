package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/juuzouuu188/task-tracker/internal/model"
)

var tasks = map[string]model.Task{}

// HTTP POST
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title string `json:"title"`
	}

	json.NewDecoder(r.Body).Decode(&input)

	if input.Title == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

// HTTP GET
func ListAllTasks(w http.ResponseWriter, r *http.Request) {
	taskList := []model.Task{}

	for _, task := range tasks {
		taskList = append(taskList, task)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskList)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")

	w.Header().Set("Content-Type", "application/json")

	// empty ID
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "missing task id",
		})
		return
	}

	//format check (numeric IDs only)
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid task id format",
		})
		return
	}

	task, exists := tasks[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "task not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
