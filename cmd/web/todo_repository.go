package main

type TodoRepository struct {
	todos []Todo
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		todos: []Todo{},
	}
}

func (tr *TodoRepository) FindAll() []Todo {
	return tr.todos
}

func (tr *TodoRepository) UpdateAll(updatedTodos []Todo) {
	tr.todos = updatedTodos
}

func (tr *TodoRepository) UpdateOne(index int, updatedTodo Todo) {
	tr.todos[index] = updatedTodo
}

func (tr *TodoRepository) CreateOne(newTodo Todo) {
	tr.todos = append(tr.todos, newTodo)
}
