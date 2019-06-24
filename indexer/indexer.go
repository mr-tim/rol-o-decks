package indexer

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/xmlpath.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"server/handlers"
	"sort"
	"store"
	"strconv"
	"strings"
)

func IndexPaths(ctx handlers.ServerContext) {
	done := make(chan bool)
	go doIndexPaths(ctx.Store, done)

	go func() {
		for {
			select {
			case _, ok := <-ctx.RestartIndexing:
				done <- true
				if !ok {
					return
				}
				go doIndexPaths(ctx.Store, done)
			}
		}
	}()
}

func doIndexPaths(s store.SlideStore, done chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	filesToIndex := make(chan string, 100)

	const workerCount = 1
	for i := 0; i < workerCount; i++ {
		go indexingWorker(filesToIndex, s)
	}

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

	paths := s.GetIndexPaths()

	var validPaths []string
	//check that that paths exist
	for _, p := range paths {
		info, err := os.Lstat(p)
		if err == nil && info.IsDir() {
			validPaths = append(validPaths, p)
		} else if err != nil {
			log.Printf("Failed to index %s because of error: %s\n", p, err)
		} else {
			log.Printf("Failed to index %s - it is not a directory\n", p)
		}
	}

	if len(paths) != len(validPaths) {
		s.SetIndexPaths(validPaths)
	}

	for _, path := range validPaths {
		log.Printf("Starting index %s", path)
		err := watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}

		//scan all the files in that directory
		filepath.Walk(path, checkAndIndexFileCallback(s, filesToIndex))
	}

	<-done
	log.Println("Stopping indexers...")
	close(filesToIndex)
}

func indexingWorker(filesToIndex <-chan string, slideStore store.SlideStore) {
	for {
		select {
		case fileToIndex, moreFiles := <-filesToIndex:
			if !moreFiles {
				log.Println("Worker finished")
				return
			}
			doIndex(slideStore, fileToIndex)
		}
	}
}

func doIndex(slideStore store.SlideStore, fileToIndex string) {
	log.Printf("Indexing %s", fileToIndex)
	r, err := zip.OpenReader(fileToIndex)
	if err != nil {
		log.Printf("Error whilst opening pptx file - skipping indexing: %s", err)
		return
	}
	defer r.Close()

	doc := store.Document{
		Path: fileToIndex,
	}

	slideCount := 0

	slideRels := make(map[int][]string, 0)

	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "ppt/slides/") {
			if !strings.HasPrefix(f.Name, "ppt/slides/_rels") {
				slideNumberStr := f.Name[len("ppt/slides/slide") : len(f.Name)-len(".xml")]
				slideNumber, _ := strconv.Atoi(slideNumberStr)
				slideCount++
				// grab the text content
				doc.Slides = append(doc.Slides, store.Slide{
					SlideNumber:     slideNumber,
					TextContent:     extractSlideContent(f),
					ThumbnailBase64: generateThumbnail(fileToIndex, slideNumber),
				})
			} else {
				slideNumberStr := f.Name[len("ppt/slides/_rels/slide") : len(f.Name)-len(".xml.rels")]
				slideNumber, _ := strconv.Atoi(slideNumberStr)
				slideRels[slideNumber] = readRefs(f)
				log.Printf("Slide %d refers to: %#v", slideNumber, slideRels[slideNumber])
			}
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

func readRefs(f *zip.File) []string {
	p := xmlpath.MustCompile("//Relationship/@Target")
	zr, _ := f.Open()
	defer zr.Close()

	root, _ := xmlpath.Parse(zr)
	i := p.Iter(root)
	refs := make([]string, 0)
	for i.Next() {
		n := i.Node()
		refs = append(refs, n.String())
	}
	return refs
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

func generateThumbnail(document string, slideNumber int) string {
	baseName := document[0 : len(document)-5]
	baseNameNoDir := baseName[strings.LastIndex(baseName, "/"):]
	tf, err := ioutil.TempFile("/tmp", fmt.Sprintf("%s_slide_%d", baseNameNoDir, slideNumber))
	if err != nil {
		//TODO: handle error
		panic(err)
	}
	// TODO: delete temp file after usage
	slideFileName := tf.Name() + ".pptx"

	log.Printf("Writing thumbnail doc to %s", slideFileName)
	createSingleSlideDocument(slideFileName, document, slideNumber)
	// generate the thumbnail
	cmd := exec.Command("qlmanage", "-t", "-s", "267", "-o", "/tmp", slideFileName)
	e := cmd.Run()
	if e != nil {
		// TODO: handle error
		panic(e)
	}

	f, err := os.Open("/tmp/" + slideFileName[strings.LastIndex(slideFileName, "/"):] + ".png")
	if err != nil {
		// TODO: handle error
		panic(err)
	}
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		// TODO: handle error
		panic(err)
	}

	buf := new(bytes.Buffer)
	enc := base64.NewEncoder(base64.StdEncoding, buf)
	enc.Write(bs)
	enc.Close()

	return buf.String()
}

func createSingleSlideDocument(slideFileName string, document string, slideNumber int) {
	file, err := os.OpenFile(slideFileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		//TODO: handle error
		panic(err)
	}
	defer file.Close()
	w := zip.NewWriter(file)
	defer w.Close()
	r, err := zip.OpenReader(document)
	if err != nil {
		//TODO: handle error
		panic(err)
	}
	defer r.Close()
	for _, zippedFile := range r.File {
		var fileWriter io.Writer
		var err error
		var toFilename string

		if zippedFile.Name == fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber) {
			toFilename = "ppt/slides/slide1.xml"
		} else if zippedFile.Name == fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", slideNumber) {
			toFilename = "ppt/slides/_rels/slide1.xml.rels"
		} else {
			toFilename = zippedFile.Name
		}
		fileWriter, err = w.Create(toFilename)
		if err != nil {
			//TODO: handle error
			panic(err)
		}
		copyZippedContent(&fileWriter, zippedFile)
	}
	w.Flush()
}
func copyZippedContent(writer *io.Writer, file *zip.File) {
	r, err := file.Open()
	if err != nil {
		//TODO: handle error
		panic(err)
	}
	defer r.Close()
	io.Copy(*writer, r)
}
