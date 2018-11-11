package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"server/handlers"
	"store"
)

func main() {
	serverContext := handlers.ServerContext{
		Store: store.InMemoryStore{},
	}
	r := mux.NewRouter()
	r.Path("/api/search").
		Queries("q", "{q}").
		HandlerFunc(handlers.SearchHandler(serverContext)).
		Name("Search")
	r.HandleFunc("/api/open/{slideId}", handlers.OpenSlideHandler(serverContext))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("ui/dist")))
	log.Fatal(http.ListenAndServe(":8000", r))
}
