package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Todo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TodoRes struct {
	Todos []Todo `json:"todos"`
}

type todoHandler struct {
	db *InMemory
}

func (t *todoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		t.create(w, r)
		return
	}

	if r.Method == http.MethodGet {
		t.read(w)
		return
	}

	if r.Method == http.MethodPatch {
		t.update(w, r)
		return
	}

	if r.Method == http.MethodDelete {
		t.delete(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (t *todoHandler) read(w http.ResponseWriter) {
	res := TodoRes{Todos: t.db.todos}
	jsonRes, err := json.Marshal(res)

	if err != nil {
		log.Fatalf("Error marshaling: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRes)
}

func (t *todoHandler) update(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	hasId := len(path) > 2

	if hasId && path[2] != "" {
		id, err := strconv.Atoi(path[2])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var updatedTodo Todo
		errDecode := json.NewDecoder(r.Body).Decode(&updatedTodo)
		if errDecode != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		isSuccessful := false
		for i, todo := range t.db.todos {
			if todo.Id == id {
				t.db.todos[i].Name = updatedTodo.Name
				isSuccessful = true
				break
			}
		}

		if isSuccessful {
			res := make(map[string]string)
			res["message"] = "Status OK"
			jsonRes, err := json.Marshal(res)
			if err != nil {
				log.Fatalf("Error marshaling: %s", err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonRes)
			return
		}

		http.Error(w, "No matched ID!", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Display path: %s", path)
}

func (t *todoHandler) delete(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	hasId := len(path) > 2

	if hasId && path[2] != "" {
		id, err := strconv.Atoi(path[2])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		removedId, err := getRemovedId(t.db.todos, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newTodos := removeAt(t.db.todos, removedId)
		t.db.todos = newTodos

		res := make(map[string]string)
		res["message"] = "Status OK"
		jsonRes, err := json.Marshal(res)
		if err != nil {
			log.Fatalf("Error marshaling: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonRes)
		return
	}

	fmt.Fprintf(w, "[DELETE] Display path: %s", path)
}

func (t *todoHandler) create(w http.ResponseWriter, r *http.Request) {
	var todo Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t.db.todos = append(t.db.todos, todo)

	res := make(map[string]string)
	res["message"] = "Status OK"
	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("Error marshaling: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRes)
}

func getRemovedId(s []Todo, id int) (int, error) {
	for i, val := range s {
		if val.Id == id {
			return i, nil
		}
	}

	return 0, fmt.Errorf("No matched ID!")
}

func removeAt[T any](s []T, i int) []T {
	lastIndex := len(s) - 1
	s[i] = s[lastIndex]
	return s[:lastIndex]
}
