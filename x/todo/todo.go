package main

import (
	"github.com/nrednav/cuid2"
)

type Todo struct {
	ID          string
	Description string
	Done        bool
}

func NewTodo(description string) Todo {
	return Todo{
		ID:          NewID(),
		Description: description,
	}
}

func NewID() string {
	return "todo_" + cuid2.Generate()
}
