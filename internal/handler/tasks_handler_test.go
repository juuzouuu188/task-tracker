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
}
