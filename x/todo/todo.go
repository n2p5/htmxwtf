package main

import (
	"slices"
	"time"

	"github.com/nrednav/cuid2"
)

type Todo struct {
	ID          string
	CreatedAt   int64
	Description string
	Done        bool
}

func NewTodo(description string) Todo {
	return Todo{
		ID:          NewID(),
		CreatedAt:   time.Now().Unix(),
		Description: description,
	}
}

func NewID() string {
	return "todo_" + cuid2.Generate()
}

func sort(todos []Todo) {
	slices.SortFunc(todos, func(a, b Todo) int {
		if a.CreatedAt < b.CreatedAt {
			return -1
		}
		if a.CreatedAt > b.CreatedAt {
			return 1
		}
		return 0
	})
}

func reverseSort(todos []Todo) {
	slices.SortFunc(todos, func(a, b Todo) int {
		if a.CreatedAt > b.CreatedAt {
			return -1
		}
		if a.CreatedAt < b.CreatedAt {
			return 1
		}
		return 0
	})
}
