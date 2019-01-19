package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"indexer"
	"log"
	"net/http"
	"os"
	"os/user"
	"server/handlers"
	"store"
	"ui"
)

type Config struct {
	Locations []string `json:"locations"`
	Database struct {
		Uri string `json:"uri"`
	} `json:"database"`
}

func main() {
	config := LoadConfig()

	serverContext := handlers.ServerContext{
		Store: store.NewSqliteStore(config.Database.Uri),
	}

	go indexer.IndexPaths(serverContext.Store, config.Locations...)

	r := mux.NewRouter()
	r.Path("/api/search").
		Queries("q", "{q}").
		HandlerFunc(handlers.SearchHandler(serverContext)).
		Name("Search")
	r.HandleFunc("/api/open/{slideId}", handlers.OpenSlideHandler(serverContext))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, ui.Ui_indexHtml())
	})
	r.HandleFunc("/ui.min.js", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, ui.Ui_uiMinJs())
	})
	log.Fatal(http.ListenAndServe(":8000", r))
}

func LoadConfig() Config {
	var config Config

	u, err := user.Current()
	if err != nil {
		panic("Could not determine home directory!")
	}


	f, err := os.Open(u.HomeDir + "/.rolodecks/config.json")
	defer f.Close()
	if err != nil {
		panic("Failed to load config!")
	}

	decoder := json.NewDecoder(f)
	decoder.Decode(&config)
	return config
}
