package main

import (
	"log"
	"net/http"
)

func main() {
	inMemory := NewInMemory()
	mux := http.NewServeMux()
	todoHandler := &todoHandler{
		db: inMemory,
	}
	mux.Handle("/todos", todoHandler)
	mux.Handle("/todos/", todoHandler)

	log.Print("server starting on :4000")

	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
