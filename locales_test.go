package locales_test

import (
	"testing"

	"github.com/nyaruka/go-locales"
	"github.com/stretchr/testify/assert"
)

func TestLocales(t *testing.T) {
	ss, err := locales.Query("es_EC", locales.LC_TIME, "abday")
	assert.NoError(t, err)
	assert.Equal(t, []string{"dom", "lun", "mar", "mié", "jue", "vie", "sáb"}, ss)

	s, err := locales.QueryString("es_EC", locales.LC_MONETARY, "currency_symbol")
	assert.NoError(t, err)
	assert.Equal(t, "$", s)

	i, err := locales.QueryInteger("es_EC", locales.LC_MONETARY, "int_frac_digits")
	assert.NoError(t, err)
	assert.Equal(t, 2, i)

	codes := locales.Codes()
	assert.Equal(t, 369, len(codes))
	assert.Equal(t, "C", codes[0])

	// LC_XLITERATE is its own category and isn't defined by any locale, so must not silently
	// resolve to another category
	_, err = locales.Query("es_EC", locales.LC_XLITERATE, "yesexpr")
	assert.EqualError(t, err, "no such category LC_XLITERATE in locale es_EC")

	ss, err = locales.Query("es_EC", locales.LC_MESSAGES, "yesexpr")
	assert.NoError(t, err)
	assert.Equal(t, []string{"^[+1sSyY]"}, ss)
}
