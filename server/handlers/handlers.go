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
	RestartIndexing chan bool
}

type handler func(w http.ResponseWriter, r *http.Request)

type searchResults struct {
	Results []store.SearchResult `json:"results"`
}

type SettingsResponse struct {
	IndexPaths []string `json:"indexPaths"`
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

func GetSettingsHandler(serverContext ServerContext) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		s := SettingsResponse{
			IndexPaths: serverContext.Store.GetIndexPaths(),
		}
		json.NewEncoder(w).Encode(s)
	}
}

func SetIndexPathsHandler(serverContext ServerContext) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		contentTypeHeader := r.Header.Get("Content-Type")
		if contentTypeHeader != "application/json" {
			badRequest(w, "JSON request required")
			return
		}
		d := json.NewDecoder(r.Body)
		var newPaths []string
		e := d.Decode(&newPaths)
		if e != nil {
			badRequest(w, "Invalid index paths specified")
		} else {
			serverContext.Store.SetIndexPaths(newPaths)
			serverContext.RestartIndexing <- true
			json.NewEncoder(w).Encode(serverContext.Store.GetIndexPaths())
		}
	}
}

func badRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(400)
	_, _ = w.Write([]byte(message))
}