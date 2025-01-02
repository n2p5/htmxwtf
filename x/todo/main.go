package main

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	middleware "github.com/go-chi/chi/v5/middleware"
)

type handlers struct {
	store Store
}

func main() {

	h := handlers{store: NewBadgerStore()}
	defer h.store.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/todos", func(r chi.Router) {
		r.Get("/", h.getTodos)
		r.Post("/", h.createTodo)
		r.Route("/{todoID}", func(r chi.Router) {
			r.Get("/", h.getTodo)
			r.Put("/", h.updateTodo)
			r.Delete("/", h.deleteTodo)
		})
	})
	http.ListenAndServe(":3000", r)
}

func (h handlers) getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.store.All()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(todos)
}

func (h handlers) getTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoID")
	todo, err := h.store.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(todo)
}

func (h handlers) createTodo(w http.ResponseWriter, r *http.Request) {
	h.store.New(NewTodo("test"))
}

func (h handlers) updateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoID")
	todo := Todo{
		ID:          id,
		Description: "test",
		Done:        true,
	}
	err := h.store.Update(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h handlers) deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoID")
	err := h.store.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
