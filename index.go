package tinysearch

type Index struct {
	Dictionary     map[string]PostingList
	TotalDocsCount int
}

func NewIndex() *Index {
	dict := make(map[string]PostingList)
	return &Index{
		Dictionary:     dict,
		TotalDocsCount: 0,
	}

}

type DocumemtID int64

type Posting struct {
	DocID         DocumemtID
	Positions     []int
	TermFrequency int
}

func NewPosting(docID DocumemtID, positions ...int) *Posting {
	return &Posting{docID, positions, len(positions)}
}
