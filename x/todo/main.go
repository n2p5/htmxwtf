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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/todos", http.StatusFound)
	})
	r.Route("/todos", func(r chi.Router) {
		r.Get("/", h.getTodos)
		r.Post("/", h.createTodo)
		r.Route("/{todoID}", func(r chi.Router) {
			r.Get("/", h.getTodo)
			r.Put("/description", h.updateTodoDescription)
			r.Put("/toggle", h.updateTodoToggle)
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
	sort(todos)
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

	todo := NewTodo(description)
	h.store.New(todo)

	w.WriteHeader(http.StatusCreated)
	TodoItem(todo).Render(r.Context(), w)
}

func (h handlers) updateTodoDescription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoID")
	todo, err := h.store.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	todo.Description = r.Form["description"][0]

	if err := h.store.Update(todo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	TodoItem(todo).Render(r.Context(), w)
}

func (h handlers) updateTodoToggle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoID")
	todo, err := h.store.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	val := r.Form["done"]

	if len(val) > 0 {
		todo.Done = true
	} else {
		todo.Done = false
	}

	if err := h.store.Update(todo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	TodoItem(todo).Render(r.Context(), w)
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
