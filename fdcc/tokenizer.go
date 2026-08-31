package fdcc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

// returned by Input when the reader yields a byte sequence which isn't valid UTF-8
var errInvalidUTF8 = errors.New("invalid UTF-8 encoding")

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

func (t *Token) String() string {
	return fmt.Sprintf("%s[%d]", t.Value, t.Type)
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
			return nil, t.endOfInput()
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
		if p == t.CommentChar {
			trailing := t.input.Prev() != rune(0) && t.input.Prev() != '\n'
			t.readComment(trailing)
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
}

// called when the input returns a zero rune, which means either that we've reached the end of it,
// that reading from it failed, or that it contains an actual NUL char - which isn't valid in an
// FDCC document and would otherwise silently truncate it.
func (t *Tokenizer) endOfInput() error {
	err := t.input.Err()

	switch {
	case errors.Is(err, io.EOF):
		return io.EOF
	case errors.Is(err, errInvalidUTF8):
		return newTokenizerError("invalid UTF-8 encoding", t.input.Position())
	case err != nil:
		return fmt.Errorf("error reading input: %w", err)
	default:
		return newTokenizerError("unexpected NUL char", t.input.Position())
	}
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

		if r == rune(0) {
			if errors.Is(t.input.Err(), io.EOF) {
				return nil, newTokenizerError("unterminated string literal", t.input.Position())
			}
			return nil, t.endOfInput()
		}
		if r == '\n' && !escaped {
			return nil, newTokenizerError("unterminated string literal", t.input.Position())
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

func (t *Tokenizer) readComment(trailing bool) string {
	b := strings.Builder{}
	for {
		r := t.input.Peek()
		if trailing && r == t.EscapeChar {
			break
		}
		if r == rune(0) {
			break // leave it for Next to report as end of input, a read error or a NUL char
		}
		if r == '\n' {
			if !trailing {
				t.input.Next()
			}
			break
		}
		t.input.Next()
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
	err    error
}

func NewInput(src io.Reader) *Input {
	return &Input{reader: bufio.NewReader(src)}
}

// reads the next rune from the reader, recording why it couldn't be read if it can't be
func (i *Input) readRune() rune {
	r, size, err := i.reader.ReadRune()
	if err != nil {
		if i.err == nil {
			i.err = err
		}
		return rune(0)
	}

	// ReadRune reports an invalid encoding as a single byte decoding to U+FFFD, whereas an actual
	// U+FFFD in the input is 3 bytes long
	if r == utf8.RuneError && size == 1 {
		if i.err == nil {
			i.err = errInvalidUTF8
		}
		return rune(0)
	}

	return r
}

func (i *Input) Peek() rune {
	if i.peeked != nil {
		return *i.peeked
	}
	r := i.readRune()
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
		next = i.readRune()
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

// Err returns the first error encountered reading from the underlying reader, or nil if no read has
// failed. It is io.EOF once the input has been fully consumed. Callers need this to tell a zero rune
// returned by Peek or Next apart from an actual NUL char in the input.
func (i *Input) Err() error {
	return i.err
}
