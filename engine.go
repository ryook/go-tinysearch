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

func (e *Engine) Flush() error {
	writer := NewIndexWriter(e.indexDir)
	return writer.Flush(e.indexer.index)
}

func (e *Engine) Search(query string, k int) ([]*SearchResult, error) {
	terms := e.tokenzer.TextToWordSequence(query)

	docs := NewSearcher(e.indexDir).SearchTopK(terms, k)

	results := make([]*SearchResult, 0, k)
	for _, result := range docs.scoreDocs {
		title, err := e.documentStore.fetchTitle(result.docID)
		if err != nil {
			return nil, err
		}
		results = append(results, &SearchResult{
			result.docID, result.score, title,
		})
	}
	return results, nil
}

type SearchResult struct {
	DocID DocumemtID
	Score float64
	Title string
}
