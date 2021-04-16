package localedata_test

import (
	"fmt"
	"testing"

	"github.com/nyaruka/go-locales"
	"github.com/nyaruka/go-locales/localedata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	database, err := localedata.LoadDatabase()
	require.NoError(t, err)

	tests := []struct {
		locale   string
		lc       localedata.LC
		keyword  string
		expected []string
		err      string
	}{
		{"es_EC", locales.LC_MONETARY, "currency_symbol", []string{"$"}, ""},
		{"es_EC", locales.LC_MONETARY, "int_frac_digits", []string{"2"}, ""},
		{"es_EC", locales.LC_TIME, "abday", []string{"dom", "lun", "mar", "mié", "jue", "vie", "sáb"}, ""}, // copies from es_BO
		{"es_EC", locales.LC_TIME, "d_fmt", []string{"%d/%m/%y"}, ""},
		{"uk_UA", locales.LC_TIME, "abday", []string{"нд", "пн", "вт", "ср", "чт", "пт", "сб"}, ""},

		{"xx_XX", locales.LC_TIME, "d_fmt", nil, "no such locale xx_XX"},
		{"es_EC", localedata.LC("LC_GOATS"), "d_fmt", nil, "no such category LC_GOATS in locale es_EC"},
		{"es_EC", locales.LC_TIME, "num_goats", nil, "no such keyword num_goats in category LC_TIME in locale es_BO"},
	}

	for _, tc := range tests {
		desc := fmt.Sprintf("%s > %s > %s", tc.locale, tc.lc, tc.keyword)
		ops, err := database.Query(tc.locale, tc.lc, tc.keyword)
		if tc.err == "" {
			assert.NoError(t, err, "unexpected error in %s", desc)
			assert.Equal(t, tc.expected, ops, "operands mismatch in %s", desc)
		} else {
			assert.EqualError(t, err, tc.err, "error mismatch in %s", desc)
			assert.Nil(t, ops)
		}
	}
}

func TestQueryString(t *testing.T) {
	database, err := localedata.LoadDatabase()
	require.NoError(t, err)

	tests := []struct {
		locale   string
		lc       localedata.LC
		keyword  string
		expected string
		err      string
	}{
		{"es_EC", locales.LC_MONETARY, "currency_symbol", "$", ""},
		{"es_EC", locales.LC_MONETARY, "int_frac_digits", "2", ""},
		{"es_EC", locales.LC_TIME, "abday", "dom", ""},
		{"es_EC", locales.LC_TIME, "d_fmt", "%d/%m/%y", ""},
		{"zh_CN", locales.LC_NAME, "name_gen", "", ""}, // empty string
	}

	for _, tc := range tests {
		desc := fmt.Sprintf("%s > %s > %s", tc.locale, tc.lc, tc.keyword)
		op, err := database.QueryString(tc.locale, tc.lc, tc.keyword)
		if tc.err == "" {
			assert.NoError(t, err, "unexpected error in %s", desc)
			assert.Equal(t, tc.expected, op, "operand mismatch in %s", desc)
		} else {
			assert.EqualError(t, err, tc.err, "error mismatch in %s", desc)
			assert.Nil(t, op)
		}
	}
}

func TestQueryInteger(t *testing.T) {
	database, err := localedata.LoadDatabase()
	require.NoError(t, err)

	tests := []struct {
		locale   string
		lc       localedata.LC
		keyword  string
		expected int
		err      string
	}{
		{"es_EC", locales.LC_MONETARY, "int_frac_digits", 2, ""},
		{"zh_CN", locales.LC_ADDRESS, "country_isbn", 7, ""},
		{"es_EC", locales.LC_TIME, "d_fmt", 0, "keyword d_fmt is not an integer"},
	}

	for _, tc := range tests {
		desc := fmt.Sprintf("%s > %s > %s", tc.locale, tc.lc, tc.keyword)
		op, err := database.QueryInteger(tc.locale, tc.lc, tc.keyword)
		if tc.err == "" {
			assert.NoError(t, err, "unexpected error in %s", desc)
			assert.Equal(t, tc.expected, op, "operand mismatch in %s", desc)
		} else {
			assert.EqualError(t, err, tc.err, "error mismatch in %s", desc)
			assert.Equal(t, 0, op)
		}
	}
}

func TestCodes(t *testing.T) {
	database, err := localedata.LoadDatabase()
	require.NoError(t, err)

	codes := database.Codes()
	assert.Equal(t, 355, len(codes))
	assert.Equal(t, "POSIX", codes[0])
	assert.Equal(t, "aa_DJ", codes[1])
	assert.Equal(t, "aa_ER", codes[2])
}
