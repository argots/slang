package ast

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type tokenKind int

const (
	operatorToken tokenKind = iota
	numberToken
	quoteToken
	identToken
)

type token struct {
	Kind tokenKind
	Loc
	Value string // Value is the raw string
}

type tokenizer struct {
	io.Reader
	LocMap
	Location string

	offset int
	reader *bufio.Reader
}

func (t *tokenizer) Next() (*token, error) {
	r, size, err := t.nextNonWhitespaceRune()
	if err != nil {
		return nil, err
	}
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return t.readNumber([]rune{r}, size)
	case '"', '\'', '`':
		return t.readQuote([]rune{r}, size)
	case '>', '<', '!':
		if t.isNextRuneEquals() {
			size++
			return t.readOperator([]rune{r, '='}, size)
		}
		fallthrough
	case '{', ':', ',', '[', ']', '(', ')', '+', '-', '*', '/', '&', '|', '.':
		return t.readOperator([]rune{r}, size)
	default:
		if !unicode.IsLetter(r) {
			return nil, t.error("unexpected character", r)
		}
		return t.readIdent([]rune{r}, size)
	}
}

func (t *tokenizer) readIdent(rs []rune, size int) (*token, error) {
	t.init()

	start := t.offset
	t.offset += size
	for {
		r, size, err := t.reader.ReadRune()
		switch {
		case err == io.EOF:
			return t.newToken(identToken, start, rs), nil
		case err != nil:
			return nil, err
		case r == '\'', r == '"', r == '`':
			tok, err := t.readQuote([]rune{r}, size)
			if err != nil {
				return nil, err
			}
			rs = append(rs, []rune(tok.Value)...)
			return t.newToken(identToken, start, rs), nil
		case unicode.IsSpace(r) || r < unicode.MaxASCII && unicode.IsPunct(r):
			t.require(t.reader.UnreadRune())
			return t.newToken(identToken, start, rs), nil
		}
		rs = append(rs, r)
		t.offset += size
	}
}

func (t *tokenizer) readOperator(rs []rune, size int) (*token, error) {
	t.init()

	start := t.offset
	t.offset += size
	return t.newToken(operatorToken, start, rs), nil
}

func (t *tokenizer) readNumber(rs []rune, size int) (*token, error) {
	t.init()

	start := t.offset
	t.offset += size
	dot := false
	for {
		r, size, err := t.reader.ReadRune()
		switch {
		case err == io.EOF:
			return t.newToken(numberToken, start, rs), nil
		case err != nil:
			return nil, err
		case r == '.' && !dot:
			dot = true
		case !unicode.IsDigit(r):
			t.require(t.reader.UnreadRune())
			return t.newToken(numberToken, start, rs), nil
		}
		rs = append(rs, r)
		t.offset += size
	}
}

func (t *tokenizer) readQuote(rs []rune, size int) (*token, error) {
	t.init()

	start := t.offset
	t.offset += size
	slash := false
	for {
		r, size, err := t.reader.ReadRune()
		switch {
		case err == io.EOF:
			return nil, io.ErrUnexpectedEOF
		case err != nil:
			return nil, err
		}
		rs = append(rs, r)
		t.offset += size
		if !slash && r == rs[0] {
			return t.newToken(quoteToken, start, rs), nil
		}
		slash = !slash && r == '\\'
	}
}

func (t *tokenizer) newToken(kind tokenKind, start int, rs []rune) *token {
	loc := t.Add(t.Location, uint32(start), uint32(t.offset))
	return &token{kind, loc, string(rs)}
}

func (t *tokenizer) error(reason string, r rune) error {
	return fmt.Errorf("%s %v", reason, r)
}

func (t *tokenizer) isNextRuneEquals() bool {
	t.init()
	r, _, err := t.reader.ReadRune()
	if err == nil && r != '=' {
		t.require(t.reader.UnreadRune())
	}
	return err == nil && r == '='
}

func (t *tokenizer) nextNonWhitespaceRune() (rune, int, error) {
	t.init()
	for {
		r, size, err := t.reader.ReadRune()
		if err != nil || !unicode.IsSpace(r) {
			return r, size, err
		}
		t.offset += size
	}
}

func (t *tokenizer) init() {
	if t.reader == nil {
		t.reader = bufio.NewReader(t.Reader)
	}
}

func (t *tokenizer) require(err error) {
	if err != nil {
		panic(err) // internal error
	}
}
