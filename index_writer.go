package tinysearch

import (
	"bufio"
	"container/list"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type IndexWriter struct {
	indexDir string
}

func NewIndexWriter(path string) *IndexWriter {
	return &IndexWriter{path}
}

func (w *IndexWriter) Flush(index *Index) error {
	for term, postinigList := range index.Dictionary {
		if err := w.postinigList(term, postinigList); err != nil {
			fmt.Printf("failed to sabe %s postinig list: %v", term, err)
		}
	}
	return w.docCount(index.TotalDocsCount)
}

func (w *IndexWriter) postinigList(term string, list PostingList) error {
	bytes, err := json.Marshal(list)
	if err != nil {
		return err
	}

	filenamae := filepath.Join(w.indexDir, term)
	file, err := os.Create(filenamae)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(bytes)
	if err != nil {
		return err
	}
	return writer.Flush()
}

func (w *IndexWriter) docCount(count int) error {
	filename := filepath.Join(w.indexDir, "_0.dc")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(strconv.Itoa(count)))
	return err
}

func (pl PostingList) MarshalJSON() ([]byte, error) {

	postings := make([]*Posting, 0, pl.Len())
	for e := pl.Front(); e != nil; e = e.Next() {
		postings = append(postings, e.Value.(*Posting))
	}
	return json.Marshal(postings)
}

func (pl *PostingList) UnmarshalJson(b []byte) error {
	var postings []*Posting
	if err := json.Unmarshal(b, &postings); err != nil {
		return err
	}

	pl.List = list.New()
	for _, posting := range postings {
		pl.add(posting)
	}
	return nil

}
