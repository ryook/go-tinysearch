package tinysearch

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type IndexReader struct {
	indexDir      string
	postingsCache map[string]*PostingList
	docCountCache int
}

func NewIndexReader(path string) *IndexReader {
	cache := make(map[string]*PostingList)
	return &IndexReader{path, cache, -1}
}

func (r *IndexReader) postingsList(terms []string) []*PostingList {
	postingLists := make([]*PostingList, 0, len(terms))
	for _, term := range terms {
		if postings := r.postings(term); postings != nil {
			postingLists = append(postingLists, postings)
		}
	}
	return postingLists
}

func (r *IndexReader) postings(term string) *PostingList {
	if postingsList, ok := r.postingsCache[term]; ok {
		return postingsList
	}
	filename := filepath.Join(r.indexDir, term)
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}

	var postingsList PostingList
	err = json.Unmarshal(bytes, &postingsList)

	if err != nil {
		return nil
	}

	r.postingsCache[term] = &postingsList
	return &postingsList
}

func (r *IndexReader) totalDocCount() int {

	if r.docCountCache > 0 {
		return r.docCountCache
	}

	filename := filepath.Join(r.indexDir, "_0.dc")
	file, err := os.Open(filename)
	if err != nil {
		return 0
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return 0
	}
	count, err := strconv.Atoi(striing(bytes))
	if err != nil {
		return 0
	}
	r.docCountCache = count
	return count
}
