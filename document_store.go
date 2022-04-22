package tinysearch

import (
	"database/sql"
	"log"
)

type DocumentStore struct {
	db *sql.DB
}

func NewDocumentStore(db *sql.DB) *DocumentStore {
	return &DocumentStore{db: db}
}

func (ds *DocumentStore) save(title string) (DocumemtID, error) {
	query := "INSERT INTO documents (document_title) VALUES (?)"
	result, err := ds.db.Exec(query, title)
	if err != nil {
		log.Fatal(err)
	}
	id, err := result.LastInsertId()
	return DocumemtID(id), err
}
