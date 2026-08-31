package localedata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// a cycle in the copy chain must be reported rather than recursed into, as a stack overflow is
// fatal in go and can't be recovered from
func TestQueryCircularCopy(t *testing.T) {
	d := &Database{locales: map[string]*Locale{
		"a": {categories: map[LC]*Category{"LC_TIME": {copiesFrom: "b"}}},
		"b": {categories: map[LC]*Category{"LC_TIME": {copiesFrom: "a"}}},
		"s": {categories: map[LC]*Category{"LC_TIME": {copiesFrom: "s"}}},
	}}

	for _, code := range []string{"a", "b", "s"} {
		ops, err := d.Query(code, "LC_TIME", "abday")
		assert.EqualError(t, err, "circular copy of category LC_TIME in locale "+code)
		assert.Nil(t, ops)
	}
}
