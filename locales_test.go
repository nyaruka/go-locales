package locales_test

import (
	"testing"

	"github.com/nyaruka/go-locales"
	"github.com/nyaruka/go-locales/localedata"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	tests := []struct {
		locale   string
		lc       localedata.LC
		key      string
		expected []string
	}{
		{"es_EC", locales.LC_MONETARY, "currency_symbol", []string{"$"}},
		{"es_EC", locales.LC_TIME, "abday", []string{"dom", "lun", "mar", "mié", "jue", "vie", "sáb"}}, // copies from es_BO
		{"es_EC", locales.LC_TIME, "d_fmt", []string{"%d/%m/%y"}},
	}

	for _, tc := range tests {
		ops, err := locales.Query(tc.locale, tc.lc, tc.key)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, ops)
	}
}
