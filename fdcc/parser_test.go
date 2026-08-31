package fdcc_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/nyaruka/go-locales/fdcc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	file, err := ioutil.ReadFile("../localedata/locales/es_BO")
	require.NoError(t, err)

	p := fdcc.NewParser(bytes.NewReader(file))
	set, err := p.Parse()
	assert.NoError(t, err)
	assert.Equal(t, 12, len(set.Categories))
	assert.Equal(t, "LC_IDENTIFICATION", set.Categories[0].Name)
	assert.Equal(t, "title", set.Categories[0].Body[0].Keyword)
	assert.Equal(t, []string{"Spanish locale for Bolivia"}, set.Categories[0].Body[0].Operands)
	assert.Equal(t, "abday", set.Categories[6].Body[0].Keyword)
	assert.Equal(t, []string{"dom", "lun", "mar", "mié", "jue", "vie", "sáb"}, set.Categories[6].Body[0].Operands)
}

func TestUnescapeUnicode(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"none", "none"},
		{"s<U00E1>bado", "sábado"},
		{"<U0627><U0644><U0623><U062D><U062F>", "الأحد"},
	}

	for _, tc := range tests {
		u := fdcc.UnescapeUnicode(tc.input)
		assert.Equal(t, tc.output, u, "output mismatch for input %s", tc.input)
	}
}

// a reader that fails part way through a stream
type failingReader struct {
	data []byte
	read int
}

func (r *failingReader) Read(p []byte) (int, error) {
	if r.read >= len(r.data) {
		return 0, errors.New("read failed")
	}
	p[0] = r.data[r.read]
	r.read++
	return 1, nil
}

const testDoc = "LC_TIME\nabday \"a\"\nEND LC_TIME\nLC_NAME\nname_fmt \"%d\"\nEND LC_NAME\n"

func TestParserErrors(t *testing.T) {
	// a NUL char is not valid and must not silently truncate the document
	withNUL := "LC_TIME\nabday \"a\"\nEND LC_TIME\n\x00LC_NAME\nname_fmt \"%d\"\nEND LC_NAME\n"
	_, err := fdcc.NewParser(strings.NewReader(withNUL)).Parse()
	assert.EqualError(t, err, "unexpected NUL char at 3:12")

	// a reader which fails mid stream must not look like a successful parse
	_, err = fdcc.NewParser(&failingReader{data: []byte(testDoc[:30])}).Parse()
	assert.EqualError(t, err, "error reading input: read failed")

	// an invalid escape in an operand list must abort rather than being skipped, which would let
	// the rest of the line be re-read as a new keyword
	_, err = fdcc.NewParser(strings.NewReader("LC_TIME\nabday \"a\";\\x;\"nope\"\nEND LC_TIME\n")).Parse()
	assert.EqualError(t, err, "unexpected char x after escape char at 2:11")

	_, err = fdcc.NewParser(strings.NewReader("LC_TIME\nabday \"a\nEND LC_TIME\n")).Parse()
	assert.EqualError(t, err, "unterminated string literal at 2:9")

	// invalid UTF-8 must be rejected rather than silently becoming U+FFFD
	_, err = fdcc.NewParser(strings.NewReader("LC_TIME\nabday \"\xff\xfe\"\nEND LC_TIME\n")).Parse()
	assert.EqualError(t, err, "invalid UTF-8 encoding at 2:7")

	// END must name the category it closes
	_, err = fdcc.NewParser(strings.NewReader("LC_TIME\nabday \"a\"\nEND LC_NAME\n")).Parse()
	assert.EqualError(t, err, "expected END LC_TIME but found END LC_NAME")

	_, err = fdcc.NewParser(strings.NewReader("LC_TIME\nabday \"a\"\nEND")).Parse()
	assert.EqualError(t, err, "unexpected EOF reading end of category LC_TIME")

	// and a valid document still parses
	set, err := fdcc.NewParser(strings.NewReader(testDoc)).Parse()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(set.Categories))
}
