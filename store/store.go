package store

type SearchResult struct {
	SlideId   string            `json:"slideId"`
	Slide     int               `json:"slide"`
	Path      string            `json:"path"`
	Thumbnail string            `json:"thumbnail"`
	Match     SearchResultMatch `json:"match"`
}

type SearchResultMatch struct {
	Text   string `json:"text"`
	Start  int    `json:"start"`
	Length int    `json:"length"`
}

type SlideStore interface {
	Search(query string) []SearchResult
	GetDocumentPathForSlideId(slideId string) string
}
