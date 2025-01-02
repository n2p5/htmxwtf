package main

import (
	"fmt"
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
	todos := []Todo{}
	fmt.Println(todos)

	err := h.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("todo")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			// TODO rip all this out into a separate function so that
			// the abstraction is cleaner in a separation of concerns
			// I want to be able to interact with the DB in a way that
			// is more abstracted from the HTTP handlers
			err := item.Value(func(v []byte) error {
				// TODO unmarshal v into a Todo struct
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

type Todo struct {
	ID          string
	Description string
	Done        bool
}
