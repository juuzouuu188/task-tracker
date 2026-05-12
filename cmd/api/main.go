package main

import (
	"log"
	"net/http"

	"github.com/juuzouuu188/task-tracker/internal/handler"
)

func main() {

	// 1. Create a new ServeMux instance
	mux := http.NewServeMux()

	//2. Map URL patterns to handler functions
	//mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
	//			switch r.Method{
	//			case http.MethodGet:
	//		handler.ListAllTasks(w,r)
	//	case http.MethodPost:
	//		handler.CreateTaskHandler(w,r)
	//	}
	//})

	mux.HandleFunc("GET /tasks", handler.ListAllTasks)
	mux.HandleFunc("POST /tasks", handler.CreateTaskHandler)

	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
