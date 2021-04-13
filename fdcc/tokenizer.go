package fdcc

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type TokenizerError struct {
	Message  string
	Position Position
}

func newTokenizerError(m string, p Position) TokenizerError {
	return TokenizerError{m, p}
}

func (e TokenizerError) Error() string {
	return fmt.Sprintf("%s at %d:%d", e.Message, e.Position.Line, e.Position.Col)
}

type TokenType int

const (
	TokenTypeIdentifier TokenType = iota
	TokenTypeString
	TokenTypeInteger
	TokenTypeEscape
	TokenTypeNewLine
	TokenTypeChar
)

type Token struct {
	Value string
	Type  TokenType
}

type Tokenizer struct {
	input       *Input
	CommentChar rune
	EscapeChar  rune
}

func NewTokenizer(i *Input) *Tokenizer {
	return &Tokenizer{
		input:       i,
		CommentChar: '#',  // default from spec
		EscapeChar:  '\\', // default from spec
	}
}

func (t *Tokenizer) Next() (*Token, error) {
	escaped := false

	for {
		p := t.input.Peek()

		if p == rune(0) {
			break
		}
		if p == t.EscapeChar && !escaped {
			t.input.Next() // discard it
			escaped = true
			continue
		}

		if escaped {
			escaped = false
			if p == '\n' {
				t.input.Next()
				continue
			} else {
				return nil, newTokenizerError(fmt.Sprintf("unexpected char %s after escape char", string(p)), t.input.Position())
			}
		}

		escaped = false

		if p == '\n' {
			t.input.Next()
			return &Token{Value: "\n", Type: TokenTypeNewLine}, nil
		}
		if p == ';' {
			t.input.Next() // discard it
			continue
		}
		if p == t.CommentChar && (t.input.Prev() == rune(0) || t.input.Prev() == '\n') {
			t.readLineRemainder()
			continue
		}
		if unicode.IsSpace(p) {
			t.input.Next()
			continue
		}
		if unicode.IsLetter(p) {
			return t.readIdentifier()
		}
		if unicode.IsDigit(p) {
			return t.readInteger()
		}
		if p == '"' {
			return t.readString()
		}

		// if it's nothing else, read and return as single char token
		t.input.Next()
		return &Token{Value: string(p), Type: TokenTypeChar}, nil
	}

	return nil, io.EOF
}

func (t *Tokenizer) readIdentifier() (*Token, error) {
	b := strings.Builder{}
	for {
		p := t.input.Peek()

		if unicode.IsLetter(p) || unicode.IsDigit(p) || p == '_' {
			r := t.input.Next()
			b.WriteRune(r)
		} else {
			break
		}
	}
	return &Token{Value: b.String(), Type: TokenTypeIdentifier}, nil
}

func (t *Tokenizer) readString() (*Token, error) {
	b := strings.Builder{}
	b.WriteRune(t.input.Next()) // opening "

	escaped := false

	for {
		r := t.input.Next()

		if r == t.EscapeChar && !escaped {
			escaped = true
			continue
		}

		if r == rune(0) || (r == '\n' && !escaped) {
			return nil, fmt.Errorf("unterminated string literal")
		}

		if r != '\n' {
			b.WriteRune(r)
		}

		escaped = false

		if r == '"' {
			break
		}
	}
	return &Token{Value: b.String(), Type: TokenTypeString}, nil
}

func (t *Tokenizer) readInteger() (*Token, error) {
	b := strings.Builder{}
	for {
		p := t.input.Peek()

		if unicode.IsDigit(p) {
			r := t.input.Next()
			b.WriteRune(r)
		} else {
			break
		}
	}
	return &Token{Value: b.String(), Type: TokenTypeInteger}, nil
}

func (t *Tokenizer) readLineRemainder() string {
	b := strings.Builder{}
	for {
		r := t.input.Next()
		if r == '\n' || r == rune(0) {
			break
		}
		b.WriteRune(r)
	}
	return b.String()
}

type Position struct {
	Line int
	Col  int
}

type Input struct {
	reader *bufio.Reader
	peeked *rune
	prev   *rune
	pos    Position
}

func NewInput(src io.Reader) *Input {
	return &Input{reader: bufio.NewReader(src)}
}

func (i *Input) Peek() rune {
	if i.peeked != nil {
		return *i.peeked
	}
	r, _, _ := i.reader.ReadRune()
	i.peeked = &r
	return r
}

func (i *Input) Next() rune {
	var next rune
	if i.peeked != nil {
		r := *i.peeked
		i.peeked = nil
		next = r
	} else {
		r, _, _ := i.reader.ReadRune()
		next = r
	}

	if i.prev == nil || *i.prev == '\n' {
		i.pos.Line++
		i.pos.Col = 0
	}
	if next != rune(0) {
		i.pos.Col++
	}
	i.prev = &next
	return next
}

func (i *Input) Prev() rune {
	if i.prev != nil {
		return *i.prev
	}
	return rune(0)
}

func (i *Input) Position() Position {
	return i.pos
}
