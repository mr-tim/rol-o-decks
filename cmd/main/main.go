package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"indexer"
	"log"
	"net/http"
	"server/handlers"
	"store"
	"ui/ui_bundle"
)

func main() {
	restartIndexing := make(chan bool)

	serverContext := handlers.ServerContext{
		Store: store.NewSqliteStore(),
		RestartIndexing: restartIndexing,
	}

	go indexer.IndexPaths(serverContext)

	r := mux.NewRouter()
	r.Path("/api/search").
		Queries("q", "{q}").
		HandlerFunc(handlers.SearchHandler(serverContext)).
		Name("Search")

	r.HandleFunc("/api/open/{slideId}", handlers.OpenSlideHandler(serverContext))

	r.Path("/api/settings").
		Methods("GET").
		HandlerFunc(handlers.GetSettingsHandler(serverContext)).
		Name("GetSettings")

	r.Path("/api/settings/indexPaths").
		Methods("PUT").
		HandlerFunc(handlers.SetIndexPathsHandler(serverContext)).
		Name("SetIndexPaths")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, ui_bundle.Ui_indexHtml())
	})

	r.HandleFunc("/ui.min.js", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, ui_bundle.Ui_uiJs())
	})

	log.Fatal(http.ListenAndServe(":8000", r))
}
