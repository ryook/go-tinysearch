package tinysearch

import (
	"database/sql"
	"io"
	"os"
	"path/filepath"
)

type Engine struct {
	tokenzer      *Tokenizer
	indexer       *Indexer
	documentStore *DocumentStore
	indexDir      string
}

func NewSearchEngine(db *sql.DB) *Engine {
	tokenizer := NewTokenizer()
	indexer := NewIndexer(tokenizer)
	documentStore := NewDocumentStore(db)

	path, ok := os.LookupEnv("INDEX_DIR_PATH")
	if !ok {
		current, _ := os.Getwd()
		path = filepath.Join(current, "_index_data")
	}

	return &Engine{
		tokenzer:      tokenizer,
		indexer:       indexer,
		documentStore: documentStore,
		indexDir:      path,
	}

}

func (e *Engine) AddDocument(title string, reader io.Reader) error {
	id, err := e.documentStore.save(title)

	if err != nil {
		return err
	}
	e.indexer.update(id, reader)
	return nil
}
