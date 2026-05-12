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

// -----------------------------
// Helpers (CONSISTENCY LAYER)
// -----------------------------

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{
		"error": msg,
	})
}

// -----------------------------
// CREATE TASK
// -----------------------------

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		writeError(w, http.StatusBadRequest, "title cannot be empty")
		return
	}

	task := model.Task{
		ID:     generateID(),
		Title:  input.Title,
		Status: model.StatusTodo,
	}

	tasks[task.ID] = task

	writeJSON(w, http.StatusCreated, task)
}

// -----------------------------
// LIST TASKS
// -----------------------------

func ListAllTasks(w http.ResponseWriter, r *http.Request) {
	taskList := []model.Task{}

	for _, task := range tasks {
		taskList = append(taskList, task)
	}

	writeJSON(w, http.StatusOK, taskList)
}

// -----------------------------
// GET TASK BY ID
// -----------------------------

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	// safer path handling
	path := strings.TrimSuffix(r.URL.Path, "/")
	id := strings.TrimPrefix(path, "/tasks/")

	if id == "" {
		writeError(w, http.StatusBadRequest, "missing task id")
		return
	}

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id format")
		return
	}

	task, exists := tasks[id]
	if !exists {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimSuffix(r.URL.Path, "/")
	path = path[len("/tasks/"):]
	id := path[:len(path)-len("/status")]

	var input struct {
		Status model.TaskStatus `json:"status"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	// validate status
	if input.Status != model.StatusTodo &&
		input.Status != model.StatusInProgress &&
		input.Status != model.StatusDone {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//check if it exists
	task, exists := tasks[id]
	if !exists {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	task.Status = input.Status
	tasks[id] = task

	writeJSON(w, http.StatusOK, task)
}

// -----------------------------
// ID GENERATOR
// -----------------------------

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
