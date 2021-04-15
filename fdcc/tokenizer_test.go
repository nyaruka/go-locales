package fdcc_test

import (
	"io"
	"strings"
	"testing"

	"github.com/nyaruka/go-locales/fdcc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		input  string
		tokens []*fdcc.Token
	}{
		{"", []*fdcc.Token{}},
		{
			" hello \n 12345 \"str\" ",
			[]*fdcc.Token{
				{"hello", fdcc.TokenTypeIdentifier},
				{"\n", fdcc.TokenTypeNewLine},
				{"12345", fdcc.TokenTypeInteger},
				{`"str"`, fdcc.TokenTypeString},
			},
		},
		{
			"id1/\nid2/\nid3", // escaped newlines
			[]*fdcc.Token{
				{"id1", fdcc.TokenTypeIdentifier},
				{"id2", fdcc.TokenTypeIdentifier},
				{"id3", fdcc.TokenTypeIdentifier},
			},
		},
		{
			"d_t_fmt \"first/\nsecond/\nthird\"", // escaped newlines within string literal
			[]*fdcc.Token{
				{"d_t_fmt", fdcc.TokenTypeIdentifier},
				{`"firstsecondthird"`, fdcc.TokenTypeString},
			},
		},
		{
			"% comment1 \n%comment2\nhola", // comments at start of line
			[]*fdcc.Token{
				{"hola", fdcc.TokenTypeIdentifier},
			},
		},
		{
			"hola %comment\nthing", // trailing comment then newline
			[]*fdcc.Token{
				{"hola", fdcc.TokenTypeIdentifier},
				{"\n", fdcc.TokenTypeNewLine},
				{"thing", fdcc.TokenTypeIdentifier},
			},
		},
		{
			"hola %comment /\nthing", // trailing comment then escaped newline
			[]*fdcc.Token{
				{"hola", fdcc.TokenTypeIdentifier},
				{"thing", fdcc.TokenTypeIdentifier},
			},
		},
	}

	for _, tc := range tests {
		in := fdcc.NewTokenizer(fdcc.NewInput(strings.NewReader(tc.input)))
		in.CommentChar = '%'
		in.EscapeChar = '/'
		tokens := make([]*fdcc.Token, 0)
		for {
			tok, err := in.Next()
			if err == io.EOF {
				break
			}
			require.NoError(t, err, "unexpected error for input %s", tc.input)
			tokens = append(tokens, tok)
		}
		assert.Equal(t, tc.tokens, tokens, "tokens mismatch for input: %s", tc.input)
	}
}

func TestInput(t *testing.T) {
	in := fdcc.NewInput(strings.NewReader(""))
	assert.Equal(t, rune(0), in.Next())
	assert.Equal(t, fdcc.Position{1, 0}, in.Position())

	in = fdcc.NewInput(strings.NewReader("h是la\nX"))
	assert.Equal(t, 'h', in.Next())
	assert.Equal(t, fdcc.Position{1, 1}, in.Position())
	assert.Equal(t, '是', in.Next())
	assert.Equal(t, fdcc.Position{1, 2}, in.Position())
	assert.Equal(t, 'l', in.Peek())
	assert.Equal(t, fdcc.Position{1, 2}, in.Position())
	assert.Equal(t, 'l', in.Peek())
	assert.Equal(t, fdcc.Position{1, 2}, in.Position())
	assert.Equal(t, 'l', in.Next())
	assert.Equal(t, fdcc.Position{1, 3}, in.Position())
	assert.Equal(t, 'a', in.Next())
	assert.Equal(t, fdcc.Position{1, 4}, in.Position())
	assert.Equal(t, '\n', in.Next())
	assert.Equal(t, fdcc.Position{1, 5}, in.Position())
	assert.Equal(t, 'X', in.Next())
	assert.Equal(t, fdcc.Position{2, 1}, in.Position())
	assert.Equal(t, rune(0), in.Next())
	assert.Equal(t, fdcc.Position{2, 1}, in.Position())
	assert.Equal(t, rune(0), in.Next())
	assert.Equal(t, fdcc.Position{2, 1}, in.Position())
}
