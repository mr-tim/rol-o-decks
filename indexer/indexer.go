package indexer

import (
	"archive/zip"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/xmlpath.v2"
	"log"
	"os"
	"path/filepath"
	"sort"
	"store"
	"strconv"
	"strings"
)

func IndexPaths(s store.SlideStore, paths ...string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	filesToIndex := make(chan string, 100)

	for i := 0; i < 4; i++ {
		go indexingWorker(filesToIndex, s)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if isPathIndexable(event.Name) {
					log.Println("event:", event)
					//TODO: deal with deletions and renames too
					if event.Op&fsnotify.Write == fsnotify.Write {
						log.Println("modified file:", event.Name)
						queueForIndexing(s, filesToIndex, event.Name)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	for _, path := range paths {
		log.Printf("Starting index %s", path)
		err := watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}

		//scan all the files in that directory
		filepath.Walk(path, checkAndIndexFileCallback(s, filesToIndex))
	}

	<-done
}

func indexingWorker(filesToIndex <-chan string, slideStore store.SlideStore) {
	select {
	case fileToIndex, moreFiles := <-filesToIndex:
		if !moreFiles {
			return
		}
		doIndex(slideStore, fileToIndex)
	}
}

func doIndex(slideStore store.SlideStore, fileToIndex string) {
	log.Printf("Indexing %s", fileToIndex)
	r, err := zip.OpenReader(fileToIndex)
	if err != nil {
		log.Printf("Error whilst opening pptx file: %s", err)
	}
	defer r.Close()

	doc := store.Document{
		Path: fileToIndex,
	}

	slideCount := 0
	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "ppt/slides/") &&
			!strings.HasPrefix(f.Name, "ppt/slides/_rels") {
			slideNumberStr := f.Name[len("ppt/slides/slide") : len(f.Name)-len(".xml")]
			slideNumber, _ := strconv.Atoi(slideNumberStr)
			log.Printf("Processing file: %s - slide %d", f.Name, slideNumber)
			slideCount++
			// grab the text content
			doc.Slides = append(doc.Slides, store.Slide{
				SlideNumber:     slideNumber,
				TextContent:     extractSlideContent(f),
				ThumbnailBase64: "",
			})
		}
	}

	sort.Sort(store.SlideNumberSorter(doc.Slides))

	slideStore.Save(doc)

	log.Printf("File %s contains %d slides", fileToIndex, slideCount)
}

func extractSlideContent(f *zip.File) string {
	p := xmlpath.MustCompile("//t")
	zr, _ := f.Open()
	defer zr.Close()
	root, _ := xmlpath.Parse(zr)
	i := p.Iter(root)
	content := make([]string, 0)
	for i.Next() {
		n := i.Node()
		content = append(content, n.String())
	}
	textContent := strings.Join(content, "\n")
	return textContent
}

func checkAndIndexFileCallback(s store.SlideStore, filesToIndex chan<- string) func(string, os.FileInfo, error) error {
	return func(path string, f os.FileInfo, err error) error {
		if isPathIndexable(path) {
			if s.IsFileModified(path, f.ModTime(), f.Size()) {
				//index the file again
				queueForIndexing(s, filesToIndex, path)
			}
		}
		return nil
	}
}

func isPathIndexable(path string) bool {
	return strings.HasSuffix(path, ".pptx")
}

func queueForIndexing(s store.SlideStore, filesToIndex chan<- string, path string) {
	log.Printf("Queuing file %s for indexing", path)
	filesToIndex <- path
}
