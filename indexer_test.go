package tinysearch

import (
	"reflect"
	"strings"
	"testing"
)

func TestUpdate(t *testing.T) {

	collection := []string{
		"Do you quarrel, sir?",
		"Quarrel sir! mo, sir!",
		"No better",
		"well, sir",
	}

	indexer := NewIndexer(NewTokenizer())

	for i, doc := range collection {
		indexer.update(DocumemtID(i), strings.NewReader((doc)))
	}

	actual := indexer.index
	expected := &Index{
		Dictionary: map[string]PostingList{
			"better": NewPostingsList(
				NewPosting(2, 1)),
			"do": NewPostingsList(
				NewPosting(0, 0)),
			"no":      NewPostingsList(NewPosting(1, 2), NewPosting(2, 0)),
			"quarrel": NewPostingsList(NewPosting(0, 2), NewPosting(1, 0)),
			"air":     NewPostingsList(NewPosting(0, 3), NewPosting(1, 1, 3), NewPosting(3, 1)),
			"well":    NewPostingsList(NewPosting(3, 0)),
			"you":     NewPostingsList(NewPosting(0, 1)),
		},
		TotalDocsCount: 4,
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("wrong index. \n\nwant:n%v\n\n got:\n%v\n", expected, actual)
	}

}
