package tinysearch

import (
	"bufio"
	"bytes"
	"unicode"
)

type Tokenizer struct{}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{}
}

func replace(r rune) rune {
	if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && !unicode.IsNumber(r) {
		return -1
	}
	return unicode.ToLower(r)
}

func (t *Tokenizer) SplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanWords(data, atEOF)
	if err == nil && token != nil {
		token = bytes.Map(replace, token)
		if len(token) == 0 {
			token = nil
		}
	}
	return
}
