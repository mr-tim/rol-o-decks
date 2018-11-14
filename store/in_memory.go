package store

import (
	"strings"
	"time"
)

type InMemoryStore struct {
	docsByPath map[string]Document
}

func NewInMemoryStore() InMemoryStore {
	return InMemoryStore{
		docsByPath: make(map[string]Document, 0),
	}
}

func (s InMemoryStore) Search(query string) []SearchResult {
	var results []SearchResult
	if len(query) > 0 {
		for path, doc := range s.docsByPath {
			for _, slide := range doc.Slides {
				if strings.Contains(strings.ToLower(slide.TextContent), strings.ToLower(query)) {
					results = append(results, SearchResult{
						SlideId:   "abcd1234",
						Slide:     slide.SlideNumber,
						Path:      path,
						Thumbnail: slide.ThumbnailBase64,
						Match: SearchResultMatch{
							Text:   slide.TextContent,
							Start:  strings.Index(strings.ToLower(slide.TextContent), strings.ToLower(query)),
							Length: len(query),
						},
					})
				}
			}
		}
	}
	return results
}

func (s InMemoryStore) GetDocumentPathForSlideId(slideId string) string {
	return "/path/to/slide.pptx"
}

func (s InMemoryStore) IsFileModified(path string, modifiedTime time.Time, fileSize int64) bool {
	return true
}

func (s InMemoryStore) Save(document Document) {
	s.docsByPath[document.Path] = document
}
