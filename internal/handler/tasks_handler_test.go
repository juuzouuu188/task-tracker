package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juuzouuu188/task-tracker/internal/model"
)

// -----------------------------
// Helpers
// -----------------------------

func newJSONRequest(method, path, body string) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func decodeBody(t *testing.T, rr *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.NewDecoder(rr.Body).Decode(v); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

// -----------------------------
// CREATE TASK
// -----------------------------

func TestCreateTask(t *testing.T) {
	tasks = map[string]model.Task{}

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantTitle  string
	}{
		{"valid task", `{"title":"Buy groceries"}`, http.StatusCreated, "Buy groceries"},
		{"empty title", `{}`, http.StatusBadRequest, ""},
		{"whitespace title", `{"title":""}`, http.StatusBadRequest, ""},
		{"invalid json", `{invalid-json}`, http.StatusBadRequest, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := newJSONRequest(http.MethodPost, "/", tc.body)
			rr := httptest.NewRecorder()

			CreateTaskHandler(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("expected %d got %d", tc.wantStatus, rr.Code)
			}

			if tc.wantStatus == http.StatusCreated {
				var task model.Task
				decodeBody(t, rr, &task)

				if task.Title != tc.wantTitle {
					t.Errorf("expected %s got %s", tc.wantTitle, task.Title)
				}
			}
		})
	}
}

// -----------------------------
// LIST TASKS
// -----------------------------

func TestListAllTasks(t *testing.T) {
	tests := []struct {
		name       string
		setup      func()
		wantLen    int
		wantTitles []string
	}{
		{
			name: "multiple tasks",
			setup: func() {
				tasks = map[string]model.Task{
					"1": {ID: "1", Title: "Buy groceries"},
					"2": {ID: "2", Title: "Buy stationary"},
				}
			},
			wantLen:    2,
			wantTitles: []string{"Buy groceries", "Buy stationary"},
		},
		{
			name: "empty list",
			setup: func() {
				tasks = map[string]model.Task{}
			},
			wantLen:    0,
			wantTitles: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			req := newJSONRequest(http.MethodGet, "/tasks", "")
			rr := httptest.NewRecorder()

			ListAllTasks(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("expected 200 got %d", rr.Code)
			}

			var result []model.Task
			decodeBody(t, rr, &result)

			if len(result) != tc.wantLen {
				t.Errorf("expected %d got %d", tc.wantLen, len(result))
			}

			found := map[string]bool{}
			for _, task := range result {
				found[task.Title] = true
			}

			for _, title := range tc.wantTitles {
				if !found[title] {
					t.Errorf("missing %s", title)
				}
			}
		})
	}
}

// -----------------------------
// GET TASK BY ID
// -----------------------------

func TestGetTask(t *testing.T) {
	tests := []struct {
		name       string
		setup      func()
		path       string
		wantStatus int
		wantTitle  string
	}{
		{
			name: "existing task",
			setup: func() {
				tasks = map[string]model.Task{
					"1": {ID: "1", Title: "Buy groceries"},
					"2": {ID: "2", Title: "Buy stationary"},
				}
			},
			path:       "/tasks/2",
			wantStatus: http.StatusOK,
			wantTitle:  "Buy stationary",
		},
		{
			name: "task not found",
			setup: func() {
				tasks = map[string]model.Task{}
			},
			path:       "/tasks/999",
			wantStatus: http.StatusNotFound,
			wantTitle:  "",
		},
		{
			name: "invalid id format",
			setup: func() {
				tasks = map[string]model.Task{}
			},
			path:       "/tasks/abc",
			wantStatus: http.StatusBadRequest,
			wantTitle:  "",
		},
		{
			name:       "empty id",
			setup:      func() {},
			path:       "/tasks/",
			wantStatus: http.StatusBadRequest,
			wantTitle:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			req := newJSONRequest(http.MethodGet, tc.path, "")
			rr := httptest.NewRecorder()

			GetTaskByID(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("expected %d got %d", tc.wantStatus, rr.Code)
			}

			if tc.wantStatus == http.StatusOK {
				var task model.Task
				decodeBody(t, rr, &task)

				if task.Title != tc.wantTitle {
					t.Errorf("expected %s got %s", tc.wantTitle, task.Title)
				}
			}
		})
	}
}
