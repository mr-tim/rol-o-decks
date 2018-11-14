package store

type Document struct {
	Path   string
	Slides []Slide
}

type Slide struct {
	SlideNumber     int
	ThumbnailBase64 string
	TextContent     string
}

type SlideNumberSorter []Slide

func (a SlideNumberSorter) Len() int           { return len(a) }
func (a SlideNumberSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SlideNumberSorter) Less(i, j int) bool { return a[i].SlideNumber < a[j].SlideNumber }
