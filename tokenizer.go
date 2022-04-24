package tinysearch

import (
	"bufio"
	"bytes"
	"strings"
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

func (t *Tokenizer) TextToWordSequence(text string) []string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(t.SplitFunc)
	var result []string
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result
}
