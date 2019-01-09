package store

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteMemoryStore struct {
	db *sql.DB
}

func NewSqliteStore(uri string) *SqliteMemoryStore {
	db := ensureDatabasePresent(uri)

	return &SqliteMemoryStore{db}
}

func ensureDatabasePresent(uri string) *sql.DB {
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(1)

	s, err := db.Prepare(`CREATE TABLE IF NOT EXISTS document (
	id INTEGER NOT NULL,
		path VARCHAR,
		created DATETIME,
		last_modified DATETIME,
		PRIMARY KEY (id)
	)`)
	checkErrors(err, "Failed to create documents table DML")
	_, err = s.Exec()
	checkErrors(err, "Failed to create documents table")

	s, err = db.Prepare(`CREATE TABLE IF NOT EXISTS slide (
		id INTEGER NOT NULL, 
		thumbnail_png BLOB, 
		slide INTEGER, 
		document_id INTEGER, 
		PRIMARY KEY (id), 
		FOREIGN KEY(document_id) REFERENCES document (id) on delete cascade
	)`)
	checkErrors(err, "Failed to create slides table DML")
	_, err = s.Exec()
	checkErrors(err, "Failed to create slides table")

	s, err = db.Prepare(`CREATE VIRTUAL TABLE IF NOT EXISTS slide_content 
  		using fts4(slide_id INTEGER, content TEXT)`)
	checkErrors(err, "Failed to create slide content table DML")
	_, err = s.Exec()
	checkErrors(err, "Failed to create slide content table")

	return db
}

func checkErrors(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", message, err))
	}
}

func (s *SqliteMemoryStore) Search(query string) []SearchResult {
	log.Printf("Searching for '%s'", query)
	rows, err := s.db.Query(`
				select
				    d.path, s.id as slide_id, s.thumbnail_png, s.slide as slide_no,
					sc.content as search_content
				from
					slide_content sc, document d, slide s
				where
				      sc.content match ?
				      and sc.slide_id = s.id
				      and s.document_id = d.id`, "'" + query + "*'")
	checkErrors(err, "Failed to search for documents")
	defer rows.Close()

	searchResults := make([]SearchResult, 0)

	for rows.Next() {
		var path string
		var slideId int
		var thumbnailPng string
		var slideNo int
		var searchContent string

		err = rows.Scan(&path, &slideId, &thumbnailPng, &slideNo, &searchContent)
		checkErrors(err, "Failed to load column values")
		log.Printf("Found document: %s, slide %s", path, slideNo)

		startIndex := strings.Index(strings.ToLower(searchContent), strings.ToLower(query))

		searchResults = append(searchResults, SearchResult{
			fmt.Sprintf("%d", slideId),
			slideNo,
			path,
			thumbnailPng,
			SearchResultMatch{
				searchContent,
				startIndex,
				len(query),
			},
		})
	}

	return searchResults
}

func (s *SqliteMemoryStore) GetDocumentPathForSlideId(slideId string) string {
	rows, err := s.db.Query(
		`select d.path from document d, slide s 
		where s.id = ? and d.id = s.document_id`, slideId)
	checkErrors(err, "Failed to lookup document path")
	defer rows.Close()

	for rows.Next() {
		var path string
		err := rows.Scan(&path)
		checkErrors(err, "Failed to retrieve document path")
		return path
	}
	panic(fmt.Sprintf("Failed to find slide with id: %s", slideId))
}

func (s *SqliteMemoryStore) IsFileModified(path string, modifiedTime time.Time, fileSize int64) bool {
	rows, err := s.db.Query(`select last_modified from document where path = ?`, path)
	checkErrors(err, "Failed to load last modified time")
	defer rows.Close()
	if rows.Next() {
		var dbLastModified time.Time
		err := rows.Scan(&dbLastModified)
		checkErrors(err, "Failed to load last modified time from database")
		return modifiedTime.After(dbLastModified)
	} else {
		return true
	}
}

func (s *SqliteMemoryStore) Save(document Document) {
	log.Printf("Saving document: %s", document.Path)

	tx, err := s.db.Begin()
	checkErrors(err, "Failed to begin transaction")

	_, err = tx.Exec(`delete from slide_content where slide_id in (
  		select s.id from slide s, document d
  		where s.document_id = d.id
  		and d.path = ?
	)`, document.Path)
	checkErrors(err, "Failed to delete slide content")

	_, err = tx.Exec(`delete from document where path = ?`, document.Path)
	checkErrors(err, "Failed to delete documents")

	res, err := tx.Exec(`insert into document (path, created, last_modified)
		values (?, ?, ?)`, document.Path, time.Now(), time.Now())
	checkErrors(err, "Failed to insert document")

	documentId, err := res.LastInsertId()
	checkErrors(err, "Failed to get last insert id")

	for _, slide := range document.Slides {
		res, err = tx.Exec(`insert into slide (thumbnail_png, slide, document_id) 
			values (?, ?, ?)`, slide.ThumbnailBase64, slide.SlideNumber, documentId)
		checkErrors(err, "Failed to insert slide")
		slideId, err := res.LastInsertId()
		checkErrors(err, "Failed to get slide id")

		_, err = tx.Exec(`insert into slide_content (slide_id, content) values (?, ?)`,
			slideId, slide.TextContent)
		checkErrors(err, "Failed to index slide text content")
	}

	err = tx.Commit()
	checkErrors(err, "Failed to commit transaction")
}