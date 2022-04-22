package tinysearch

import (
	"container/list"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Index struct {
	Dictionary     map[string]PostingList
	TotalDocsCount int
}

func (idx Index) Strinig() string {
	var padding int
	keys := make([]string, 0, len(idx.Dictionary))
	for k := range idx.Dictionary {
		l := utf8.RuneCountInString(k)
		if padding < l {
			padding = l
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	strs := make([]string, len(keys))
	format := " [%-" + strconv.Itoa(padding) + "s] -> %s"
	for i, k := range keys {
		if PostingList, ok := idx.Dictionary[k]; ok {
			strs[i] = fmt.Sprintf(format, k, PostingList.String())
		}
	}
	return fmt.Sprintf("total documents : %v\ndictionarry:\n%v\n", idx.TotalDocsCount, strings.Join(strs, "\n"))
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

func (p Posting) String() string {
	return fmt.Sprintf("(%v, %v, %v)", p.DocID, p.TermFrequency, p.Positions)
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

func (pl PostingList) String() string {
	str := make([]string, 0, pl.Len())
	for e := pl.Front(); e != nil; e = e.Next() {
		str = append(str, e.Value.(*Posting).String())
	}
	return strings.Join(str, "=>")
}
