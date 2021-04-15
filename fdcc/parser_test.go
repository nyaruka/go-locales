package fdcc_test

import (
	"bytes"
	"io/ioutil"
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
	assert.Equal(t, "title", set.Categories[0].Body[0].Identifier)
	assert.Equal(t, []string{"Spanish locale for Bolivia"}, set.Categories[0].Body[0].Operands)
	assert.Equal(t, "abday", set.Categories[6].Body[0].Identifier)
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
