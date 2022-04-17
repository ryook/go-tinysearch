package tinysearch

import "container/list"

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

type PostingList struct {
	*list.List
}

func NewPostingsList(postings ...*Posting) PostingList {
	l := list.New()
	for _, posting := range postings {
		l.PushBack(posting)
	}
	return PostingList{l}
}

func (pl PostingList) add(p *Posting) {
	pl.PushBack(p)
}

func (pl PostingList) last() *Posting {
	e := pl.List.Back()
	if e == nil {
		return nil
	}
	return e.Value.(*Posting)
}

// ポスティングをリストに追加
// ポスティングリストの最後のドキュメントIDを取得して
// 一致していればポスティングを追加
// 一致してなければpositionを追加
func (pl PostingList) Add(new *Posting) {
	last := pl.last()
	if last == nil || last.DocID != new.DocID {
		pl.add(new)
		return
	}
	last.Positions = append(last.Positions, new.Positions...)
	last.TermFrequency++
}
