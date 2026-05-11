package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juuzouuu188/task-tracker/internal/model"
)

func TestCreateTask(t *testing.T) {
	tasks = map[string]model.Task{} // reset

	t.Run("valid task", func(t *testing.T) {
		body := bytes.NewBufferString(`{"title": "Buy groceries"}`)

		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		CreateTaskHandler(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("expected 201 got %d", rr.Code)
		}

		var task model.Task
		if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if task.Title != "Buy groceries" {
			t.Errorf("expected title 'Buy groceries' got '%s'", task.Title)
		}

		if task.Status != model.StatusTodo {
			t.Errorf("expected status 'todo' got '%s'", task.Status)
		}
	})

	t.Run("empty title", func(t *testing.T) {
		body := bytes.NewBufferString(`{}`)

		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		CreateTaskHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400 got %d", rr.Code)
		}
	})

	t.Run("whitespace title", func(t *testing.T) {
		body := bytes.NewBufferString(`{"title": ""}`)

		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		CreateTaskHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400 got %d", rr.Code)
		}

	})

	t.Run("invalid json", func(t *testing.T) {
		body := bytes.NewBufferString(`{invalid-json}`)

		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		CreateTaskHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400 got %d", rr.Code)
		}
	})

}

func TestListAllTasks(t *testing.T) {

	t.Run("multiple tasks", func(t *testing.T) {
		tasks = map[string]model.Task{}

		tasks["1"] = model.Task{ID: "1", Title: "Buy groceries", Status: model.StatusTodo}
		tasks["2"] = model.Task{ID: "2", Title: "Buy stationary", Status: model.StatusTodo}

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		ListAllTasks(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 got %d", rr.Code)
		}

		var result []model.Task
		if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 tasks got %d", len(result))
		}

		found := map[string]bool{}

		for _, task := range result {
			found[task.Title] = true
		}

		if !found["Buy groceries"] {
			t.Errorf("missing 'Buy groceries'")
		}

		if !found["Buy stationary"] {
			t.Errorf("missing 'Buy stationary'")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		tasks = map[string]model.Task{}

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		ListAllTasks(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 got %d", rr.Code)
		}

		var result []model.Task
		if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(result) != 0 {
			t.Errorf("expected 0 tasks got %d", len(result))
		}
	})
}

func TestGetTask(t *testing.T) {

	t.Run("existing task", func(t *testing.T) {
		tasks = map[string]model.Task{}

		tasks["1"] = model.Task{ID: "1", Title: "Buy groceries"}
		tasks["2"] = model.Task{ID: "2", Title: "Buy stationary"}

		req := httptest.NewRequest(http.MethodGet, "/tasks/2", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		GetTaskByID(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 got %d", rr.Code)
		}

		var task model.Task
		if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if task.Title != "Buy stationary" {
			t.Errorf("expected title 'Buy stationary' got '%s'", task.Title)
		}
	})

	t.Run("task not found", func(t *testing.T) {
		tasks = map[string]model.Task{}

		req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		GetTaskByID(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected 404 got %d", rr.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		tasks = map[string]model.Task{}

		req := httptest.NewRequest(http.MethodGet, "/tasks/abcdefg", nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		GetTaskByID(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400 got %d", rr.Code)
		}
	})

	t.Run("empty id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tasks/", nil)
		rr := httptest.NewRecorder()

		GetTaskByID(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400 got %d", rr.Code)
		}
	})
}
