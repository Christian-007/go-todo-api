package main

import (
	"log"
	"net/http"
)

func main() {
	todoRepository := NewTodoRepository()
	mux := http.NewServeMux()
	todoHandler := &todoHandler{
		todoRepository: todoRepository,
	}
	mux.Handle("/todos", todoHandler)
	mux.Handle("/todos/", todoHandler)

	log.Print("server starting on :4000")

	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
