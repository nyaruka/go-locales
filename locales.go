package locales

import "github.com/nyaruka/go-locales/localedata"

const (
	LC_IDENTIFICATION localedata.LC = "LC_IDENTIFICATION" // Versions and status of categories
	LC_CTYPE          localedata.LC = "LC_CTYPE"          // Character classification, case conversion and code transformation.
	LC_COLLATE        localedata.LC = "LC_COLLATE"        // Collation order.
	LC_TIME           localedata.LC = "LC_TIME"           // Date and time formats.
	LC_NUMERIC        localedata.LC = "LC_NUMERIC"        // Numeric, non-monetary formatting.
	LC_MONETARY       localedata.LC = "LC_MONETARY"       // Monetary formatting.
	LC_MESSAGES       localedata.LC = "LC_MESSAGES"       // Formats of informative and diagnostic messages and interactive responses.
	LC_XLITERATE      localedata.LC = "LC_MESSAGES"       // Character transliteration.
	LC_NAME           localedata.LC = "LC_NAME"           // Format of writing personal names.
	LC_ADDRESS        localedata.LC = "LC_ADDRESS"        // Format of postal addresses.
	LC_TELEPHONE      localedata.LC = "LC_TELEPHONE"      // Format for telephone numbers, and other telephone information.
)

var database *localedata.Database

func init() {
	var err error
	database, err = localedata.LoadDatabase()
	if err != nil {
		panic(err)
	}
}

// Query returns the operands of the given locale + category + keyword
func Query(localeName string, lc localedata.LC, keyword string) ([]string, error) {
	return database.Query(localeName, lc, keyword)
}

// QueryString is a helper for keywords which are a single string
func QueryString(localeName string, lc localedata.LC, keyword string) (string, error) {
	return database.QueryString(localeName, lc, keyword)
}

// QueryInteger is a helper for keywords which are a single integer
func QueryInteger(localeName string, lc localedata.LC, keyword string) (int, error) {
	return database.QueryInteger(localeName, lc, keyword)
}

// Codes returns the list of all locale codes, sorted alphabetically
func Codes() []string {
	return database.Codes()
}
