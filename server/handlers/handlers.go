package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os/exec"
	"store"
)

type ServerContext struct {
	Store store.SlideStore
}

type handler func(w http.ResponseWriter, r *http.Request)

type searchResults struct {
	Results []store.SearchResult `json:"results"`
}

func SearchHandler(serverContext ServerContext) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		searchTerm := mux.Vars(r)["q"]
		results := searchResults{
			serverContext.Store.Search(searchTerm),
		}
		json.NewEncoder(w).Encode(results)
	}
}

func OpenSlideHandler(serverContext ServerContext) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		slideId := mux.Vars(r)["slideId"]
		documentPath := serverContext.Store.GetDocumentPathForSlideId(slideId)
		w.Write([]byte("opening: " + documentPath))

		c := exec.Command("open", documentPath)
		err := c.Start()
		if err != nil {
			log.Printf("Error whilst trying to open %s: %s", documentPath, err)
		}
	}
}
