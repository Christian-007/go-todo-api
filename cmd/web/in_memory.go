package main

type InMemory struct {
	todos []Todo
}

func NewInMemory() *InMemory {
	return &InMemory{}
}
