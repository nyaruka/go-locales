package locales_test

import (
	"fmt"
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
		{"es_EC", locales.LC_MONETARY, "int_frac_digits", []string{"2"}},
		{"es_EC", locales.LC_TIME, "abday", []string{"dom", "lun", "mar", "mié", "jue", "vie", "sáb"}}, // copies from es_BO
		{"es_EC", locales.LC_TIME, "d_fmt", []string{"%d/%m/%y"}},
	}

	for _, tc := range tests {
		desc := fmt.Sprintf("%s > %s > %s", tc.locale, tc.lc, tc.key)
		ops, err := locales.Query(tc.locale, tc.lc, tc.key)
		assert.NoError(t, err, "unexpected error in %s", desc)
		assert.Equal(t, tc.expected, ops, "operands mismatch in %s", desc)
	}
}

func TestQueryString(t *testing.T) {
	tests := []struct {
		locale   string
		lc       localedata.LC
		key      string
		expected string
	}{
		{"es_EC", locales.LC_MONETARY, "currency_symbol", "$"},
		{"es_EC", locales.LC_MONETARY, "int_frac_digits", "2"},
		{"es_EC", locales.LC_TIME, "abday", "dom"},
		{"es_EC", locales.LC_TIME, "d_fmt", "%d/%m/%y"},
	}

	for _, tc := range tests {
		desc := fmt.Sprintf("%s > %s > %s", tc.locale, tc.lc, tc.key)
		op, err := locales.QueryString(tc.locale, tc.lc, tc.key)
		assert.NoError(t, err, "unexpected error in %s", desc)
		assert.Equal(t, tc.expected, op, "operand mismatch in %s", desc)
	}
}

func TestQueryInteger(t *testing.T) {
	tests := []struct {
		locale   string
		lc       localedata.LC
		key      string
		expected int
	}{
		{"es_EC", locales.LC_MONETARY, "int_frac_digits", 2},
		{"zh_CN", locales.LC_ADDRESS, "country_isbn", 7},
	}

	for _, tc := range tests {
		desc := fmt.Sprintf("%s > %s > %s", tc.locale, tc.lc, tc.key)
		op, err := locales.QueryInteger(tc.locale, tc.lc, tc.key)
		assert.NoError(t, err, "unexpected error in %s", desc)
		assert.Equal(t, tc.expected, op)
		assert.Equal(t, tc.expected, op, "operand mismatch in %s", desc)
	}
}
