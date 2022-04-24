package tinysearch

import (
	"fmt"
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
