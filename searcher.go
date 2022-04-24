package tinysearch

import (
	"fmt"
	"sort"
)

type TopDocs struct {
	totalHits int
	scoreDocs []*ScoreDoc
}

func (t *TopDocs) String() string {
	return fmt.Sprintf("\ntotal hits: %v\nresults: %v\n",
		t.totalHits, t.scoreDocs)
}

type ScoreDoc struct {
	docID DocumemtID
	score float64
}

func (d ScoreDoc) String() string {
	return fmt.Sprintf("docId: %v, Score: %v", d.docID, d.score)
}

type Searcher struct {
	indexReader *IndexReader
	cursors     []*Cursor
}

func NewSearcher(path string) *Searcher {
	return &Searcher{indexReader: NewIndexReader(path)}
}

func (s *Searcher) SearchTopK(query []string, k int) *TopDocs {
	results := s.search(query)

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	total := len(results)
	if len(results) > k {
		results = results[:k]
	}

	return &TopDocs{
		totalHits: total,
		scoreDocs: results,
	}
}
