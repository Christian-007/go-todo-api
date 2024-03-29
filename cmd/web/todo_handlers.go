package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type Todo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (t Todo) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Name, validation.Required, validation.Length(2, 200)),
	)
}

type CollectionRes[Entity any] struct {
	Results []Entity `json:"results"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type todoHandler struct {
	todoRepository *TodoRepository
}

func (t *todoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		t.create(w, r)
		return
	}

	if r.Method == http.MethodGet {
		t.read(w, r)
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

func (t *todoHandler) read(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	hasId := len(path) > 2

	if hasId {
		id := path[2]
		t.readOne(w, id)
		return
	}

	t.readAll(w)
}

func (t *todoHandler) readOne(w http.ResponseWriter, todoId string) {
	if todoId == "" {
		sendResponse(w, http.StatusBadRequest, ErrorResponse{Message: "Missing ID or invalid path"})
		return
	}

	todos := t.todoRepository.FindAll()
	selectedIndex, err := getRemovedId(todos, todoId)
	if err != nil {
		sendResponse(w, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	sendResponse(w, http.StatusAccepted, todos[selectedIndex])
}

func (t *todoHandler) readAll(w http.ResponseWriter) {
	todos := t.todoRepository.FindAll()
	res := CollectionRes[Todo]{Results: todos}
	sendResponse(w, http.StatusAccepted, res)
}

func (t *todoHandler) update(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	hasId := len(path) > 2

	if hasId && path[2] != "" {
		var updatedTodo Todo
		decodeErr := json.NewDecoder(r.Body).Decode(&updatedTodo)
		if decodeErr != nil {
			sendResponse(w, http.StatusBadRequest, ErrorResponse{Message: decodeErr.Error()})
			return
		}

		if validationErr := updatedTodo.Validate(); validationErr != nil {
			sendResponse(w, http.StatusBadRequest, ErrorResponse{Message: validationErr.Error()})
			return
		}

		id := path[2]
		var updatedTodoIndex int
		isSuccessful := false
		todos := t.todoRepository.FindAll()

		for i, todo := range todos {
			if todo.Id == id {
				updatedTodo.Id = id
				t.todoRepository.UpdateOne(i, updatedTodo)
				isSuccessful = true
				updatedTodoIndex = i
				break
			}
		}

		if isSuccessful {
			sendResponse(w, http.StatusOK, todos[updatedTodoIndex])
			return
		}

		sendResponse(w, http.StatusBadRequest, ErrorResponse{Message: "No matched ID!"})
		return
	}

	sendResponse(w, http.StatusBadRequest, ErrorResponse{Message: "Missing ID or invalid path"})
}

func (t *todoHandler) delete(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	hasId := len(path) > 2

	if hasId && path[2] != "" {
		id := path[2]
		todos := t.todoRepository.FindAll()
		removedId, err := getRemovedId(todos, id)
		if err != nil {
			sendResponse(w, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
			return
		}

		newTodos := removeAt(todos, removedId)
		t.todoRepository.UpdateAll(newTodos)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	sendResponse(w, http.StatusBadRequest, ErrorResponse{Message: "Missing ID or invalid path"})
}

func (t *todoHandler) create(w http.ResponseWriter, r *http.Request) {
	var todo Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		sendResponse(w, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	if validationErr := todo.Validate(); validationErr != nil {
		sendResponse(w, http.StatusBadRequest, ErrorResponse{Message: validationErr.Error()})
		return
	}

	todo.Id = uuid.New().String()
	t.todoRepository.CreateOne(todo)
	sendResponse(w, http.StatusOK, todo)
}

func getRemovedId(s []Todo, id string) (int, error) {
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

func sendResponse(w http.ResponseWriter, statusCode int, response any) {
	jsonRes, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonRes)
}
