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
	r.Get("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))).ServeHTTP)

	http.ListenAndServe(":3000", r)

}

func (h handlers) getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.store.All()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	TodoPage(todos).Render(r.Context(), w)
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

// NOTES: this is just a really rough sketch on how I wan to handle
// submitted data. By using form data only, it simplifies the process
// of handling the data. I can use the form data to create a new todo.
// I need to think about string safety and how to handle that.

func (h handlers) createTodo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	description := r.FormValue("description")

	if description == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.store.New(NewTodo(description))

	w.WriteHeader(http.StatusCreated)

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
