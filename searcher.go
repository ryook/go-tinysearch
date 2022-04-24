package tinysearch

import (
	"fmt"
	"math"
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
	println("aa")
	println(results)

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	total := len(results)
	if len(results) > k {
		results = results[:k]
	}

	println(total)
	println(results)

	return &TopDocs{
		totalHits: total,
		scoreDocs: results,
	}
}

func (s *Searcher) search(query []string) []*ScoreDoc {
	if s.openCursors(query) == 0 {
		return []*ScoreDoc{}
	}

	c := s.cursors[0]
	cursors := s.cursors[1:]

	docs := make([]*ScoreDoc, 0)

	for !c.Empty() {
		var nextDocId DocumemtID

		for _, cursor := range cursors {
			if cursor.NextDoc(c.DocId()); cursor.Empty() {
				return docs
			}

			if cursor.DocId() != c.DocId() {
				nextDocId = cursor.DocId()
				break
			}
		}

		if nextDocId > 0 {
			if c.NextDoc(nextDocId); c.Empty() {
				return docs
			}
		} else {
			docs = append(docs, &ScoreDoc{
				docID: c.DocId(),
				score: s.calcScore(),
			})
			c.Next()
		}
	}

	return docs
}

func (s *Searcher) openCursors(query []string) int {
	postings := s.indexReader.postingsList(query)
	if len(postings) == 0 {
		return 0
	}

	sort.Slice(postings, func(i, j int) bool {
		return postings[i].Len() < postings[j].Len()
	})

	cursors := make([]*Cursor, len(postings))
	for i, postingList := range postings {
		cursors[i] = postingList.OpenCursor()
	}
	s.cursors = cursors
	return len(cursors)
}

func (s *Searcher) calcScore() float64 {
	var score float64
	for i := 0; i < len(s.cursors); i++ {
		termFreq := s.cursors[i].Posting().TermFrequency
		docCount := s.cursors[i].PostingList.Len()
		totalDocCount := s.indexReader.totalDocCount()
		score += calcTF(termFreq) * calIDF(totalDocCount, docCount)
	}
	return score
}

func calcTF(termCount int) float64 {
	if termCount <= 0 {
		return 0
	}
	return math.Log2(float64(termCount)) + 1
}

func calIDF(N, df int) float64 {
	return math.Log2(float64(N) / float64(df))
}
