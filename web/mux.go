package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type ServeMux struct {
	http.ServeMux
}

func New() *ServeMux {
	return &ServeMux{*http.NewServeMux()}
}

type Request *http.Request
type Handler func(ResponseWriter, Request) error

func (m *ServeMux) Route(pattern string, handler Handler) {
	m.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s", r.Method, r.URL.Path)
		writer := ResponseWriter{w}
		err := handler(writer, r)
		if err == nil {
			return
		}
		writer.SendError(500, "Internal server error")
		fmt.Println("Route error:", err)
	})
}

func (m *ServeMux) Listen() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	log.Println("Starting at http://localhost:" + PORT)
	err := http.ListenAndServe(":"+PORT, m)
	if err != nil {
		log.Fatal(err)
	}
}
