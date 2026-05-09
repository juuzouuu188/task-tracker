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
	body := bytes.NewBufferString(`{"title": "Buy groceries"}`)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application.json")

	rr := httptest.NewRecorder()

	CreateTaskHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201 got %d", rr.Code)
	}

	var task model.Task
	json.NewDecoder(rr.Body).Decode(&task)

	if task.Title != "Buy groceries" {
		t.Errorf("expected title 'Buy groceries' got '%s'", task.Title)
	}

	if task.Status != model.StatusTodo {
		t.Errorf("expected status 'todo' got '%s'", task.Status)
	}

	t.Run("Empty Title", func(t *testing.T) {
		body := bytes.NewBufferString(`{"title": ""}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set("Content-Type", "application.json")

		rr := httptest.NewRecorder()
		CreateTaskHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400 got %d", rr.Code)
		}
	})

}

func TestListAllTasks(t *testing.T) {
	tasks = map[string]model.Task{} // reset
	// seed the map with a task
	tasks["1"] = model.Task{ID: "1", Title: "Buy groceries", Status: model.StatusTodo}
	tasks["2"] = model.Task{ID: "2", Title: "Buy stationary", Status: model.StatusTodo}

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	req.Header.Set("Content-Type", "application.json")

	rr := httptest.NewRecorder()

	ListAllTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 got %d", rr.Code)
	}

	var result []model.Task
	json.NewDecoder(rr.Body).Decode(&result)

	if len(result) != 2 {
		t.Errorf("expected 2 tasks got %d", len(result))
	}

	//Empty List Should still be a valid Response
	t.Run("Empty ListAllTasks", func(t *testing.T) {
		tasks = map[string]model.Task{} // reset
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		req.Header.Set("Content-Type", "application.json")

		rr := httptest.NewRecorder()

		ListAllTasks(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200 got %d", rr.Code)
		}

		var result []model.Task
		json.NewDecoder(rr.Body).Decode(&result)

		if len(result) != 0 {
			t.Errorf("expected 0 tasks got %d", len(result))
		}

	})
}
