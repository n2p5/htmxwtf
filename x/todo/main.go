package main

import (
	"log"
	"net/http"

	badger "github.com/dgraph-io/badger/v4"
	chi "github.com/go-chi/chi/v5"
	middleware "github.com/go-chi/chi/v5/middleware"
)

type handlers struct {
	db *badger.DB
}

func main() {

	// Initialize Badger
	db, err := badger.Open(badger.DefaultOptions("./badgerdb"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	h := handlers{db: db}

	r.Route("/todos", func(r chi.Router) {
		r.Get("/", h.getTodos)
		r.Route("/{todoID}", func(r chi.Router) {
			r.Get("/", h.getTodo)
			r.Post("/", h.createTodo)
			r.Put("/", h.updateTodo)
			r.Delete("/", h.deleteTodo)
		})
	})
	http.ListenAndServe(":3000", r)
}

func (h handlers) getTodos(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getTodos"))
}

func (h handlers) getTodo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getTodo"))
}

func (h handlers) createTodo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("createTodo"))
}

func (h handlers) updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("updateTodo"))
}

func (h handlers) deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("deleteTodo"))
}
