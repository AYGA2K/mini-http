package main

import (
	"github.com/AYGA2K/minihttp/internal/server"
)

func main() {
	mux := server.NewServeMux()

	// GET  /           → plain hello
	mux.HandleFunc("/", func(w server.ResponseWriter, r *server.Request) {
		w.StatusCode(200)
		w.Header("X-Foo", "bar")
		w.Body([]byte("Hello server"))
		w.Send()
	})

	// GET  /user/<name> → dynamic path
	mux.HandleFunc("/user/:name", func(w server.ResponseWriter, r *server.Request) {
		name := r.Params["name"]
		w.StatusCode(200)
		w.Body([]byte("Hello, " + name))
		w.Send()
	})

	// GET  /search?q=foo → query param
	mux.HandleFunc("/search", func(w server.ResponseWriter, r *server.Request) {
		w.StatusCode(200)
		w.Body([]byte("Search for: " + r.Query["q"]))
		w.Send()
	})

	// POST /submit       → echo body
	mux.HandleFunc("/submit", func(w server.ResponseWriter, r *server.Request) {
		if r.Method != "POST" {
			w.StatusCode(405)
			w.Body([]byte("Method not allowed"))
			w.Send()
			return
		}
		w.StatusCode(200)
		w.Body([]byte("Received: " + r.Body))
		w.Send()
	})

	server.ListenAndServe(":8080", mux)
}
