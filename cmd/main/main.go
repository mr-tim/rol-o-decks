package main

import (
	"github.com/gorilla/mux"
	"indexer"
	"log"
	"net/http"
	"server/handlers"
	"store"
)

func main() {
	serverContext := handlers.ServerContext{
		Store: store.NewInMemoryStore(),
	}

	go indexer.IndexPaths(serverContext.Store, "data/slides1", "data/slides2")

	r := mux.NewRouter()
	r.Path("/api/search").
		Queries("q", "{q}").
		HandlerFunc(handlers.SearchHandler(serverContext)).
		Name("Search")
	r.HandleFunc("/api/open/{slideId}", handlers.OpenSlideHandler(serverContext))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("ui/dist")))
	log.Fatal(http.ListenAndServe(":8000", r))
}
