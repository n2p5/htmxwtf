package main

import (
	"fmt"
	"net/http"
	"time"
)

func logHeaders(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("--- Headers: ---")
		for name, values := range r.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)
			}
		}
		handler(w, r)
	}
}

func mouseEnteredHandler() http.HandlerFunc {
	counter := 0
	return func(w http.ResponseWriter, r *http.Request) {
		if counter >= 10 {
			w.WriteHeader(286)
		} else {
			counter++
		}
		w.Write([]byte(fmt.Sprintf("<p>Mouse entered %d times</p>", counter)))
	}
}

func main() {

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	http.Handle("/", http.FileServer(http.Dir("./")))

	http.HandleFunc("/clicked", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<div id="parent-div"><p>woot<p></div>`))
	})

	http.HandleFunc("/click-delayed", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Write([]byte(`<p>This was a delayed click</p>`))
	})

	http.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<p>Account page</p>`))
	})

	http.HandleFunc("/mouse_entered", logHeaders(mouseEnteredHandler()))

	http.ListenAndServe(":8080", nil)
}
