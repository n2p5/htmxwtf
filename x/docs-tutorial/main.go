package main

import (
	"fmt"
	"net/http"
)

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

	http.HandleFunc("/mouse_entered", mouseEnteredHandler())

	http.ListenAndServe(":8080", nil)
}
