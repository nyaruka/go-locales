package localedata

import (
	"embed"
	"fmt"
	"sort"
	"strconv"

	"github.com/nyaruka/go-locales/fdcc"

	"github.com/pkg/errors"
)

//go:embed locales/*
var static embed.FS

type LC string

type Database struct {
	locales map[string]*Locale
	codes   []string
}

type Locale struct {
	categories map[LC]*Category
}

func newLocale(s *fdcc.Set) *Locale {
	categories := make(map[LC]*Category, len(s.Categories))
	for _, c := range s.Categories {
		categories[LC(c.Name)] = newCategory(c)
	}

	return &Locale{categories: categories}
}

type Category struct {
	copiesFrom string
	values     map[string][]string
}

func newCategory(c *fdcc.Category) *Category {
	copiesFrom := c.CopiesFrom()
	if copiesFrom != "" {
		return &Category{copiesFrom: copiesFrom}
	}

	values := make(map[string][]string, len(c.Body))
	for _, l := range c.Body {
		values[l.Identifier] = l.Operands
	}

	return &Category{values: values}
}

func LoadDatabase() (*Database, error) {
	files, err := static.ReadDir("locales")
	if err != nil {
		return nil, err
	}

	locales := make(map[string]*Locale, len(files))
	codes := make([]string, 0, len(files))

	for _, f := range files {
		code := f.Name()

		file, err := static.Open(fmt.Sprintf("locales/%s", code))
		if err != nil {
			return nil, errors.Wrapf(err, "unable to open file %s", code)
		}

		defer file.Close()

		p := fdcc.NewParser(file)

		set, err := p.Parse()
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse file %s", code)
		}

		locales[code] = newLocale(set)
		codes = append(codes, code)
	}

	sort.Strings(codes)

	return &Database{locales, codes}, nil
}

// Query returns the operands of the given locale + category + key
func (d *Database) Query(code string, lc LC, key string) ([]string, error) {
	locale := d.locales[code]
	if locale == nil {
		return nil, fmt.Errorf("no such locale %s", code)
	}

	category := locale.categories[lc]
	if category == nil {
		return nil, fmt.Errorf("no such category %s in locale %s", lc, code)
	}

	if category.copiesFrom != "" {
		return d.Query(category.copiesFrom, lc, key)
	}

	operands, exists := category.values[key]
	if !exists {
		return nil, fmt.Errorf("no such key %s in category %s in locale %s", key, lc, code)
	}

	return operands, nil
}

// QueryString is a helper for keys which are a single string
func (d *Database) QueryString(code string, lc LC, key string) (string, error) {
	ops, err := d.Query(code, lc, key)
	if err != nil {
		return "", err
	}
	if len(ops) < 1 {
		return "", fmt.Errorf("key %s has no operands", key)
	}
	return ops[0], nil
}

// QueryInteger is a helper for keys which are a single integer
func (d *Database) QueryInteger(code string, lc LC, key string) (int, error) {
	op, err := d.QueryString(code, lc, key)
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(op)
	if err != nil {
		return 0, fmt.Errorf("key %s is not an integer", key)
	}
	return val, nil
}

// Codes returns the list of all locale codes (mostly BCP47 tho includes other special values such as POSIX, i18n etc), sorted alphabetically
func (d *Database) Codes() []string {
	return d.codes
}
