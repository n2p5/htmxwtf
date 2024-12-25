package main

import (
	"fmt"
	"net/http"
	"time"
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

	http.HandleFunc("/click-delayed", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Write([]byte(`<p>This was a delayed click</p>`))
	})

	http.HandleFunc("/mouse_entered", mouseEnteredHandler())

	http.ListenAndServe(":8080", nil)
}
