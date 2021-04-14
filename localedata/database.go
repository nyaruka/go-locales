package localedata

import (
	"embed"
	"fmt"
	"strconv"

	"github.com/nyaruka/go-locales/fdcc"

	"github.com/pkg/errors"
)

//go:embed locales/*
var static embed.FS

type LC string

type Database struct {
	locales map[string]*Locale
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

	for _, f := range files {
		name := f.Name()

		file, err := static.Open(fmt.Sprintf("locales/%s", name))
		if err != nil {
			return nil, errors.Wrapf(err, "unable to open file %s", name)
		}

		defer file.Close()

		p := fdcc.NewParser(file)

		set, err := p.Parse()
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse file %s", name)
		}

		locales[name] = newLocale(set)
	}
	return &Database{locales}, nil
}

// Query returns the operands of the given locale + category + key
func (d *Database) Query(localeName string, lc LC, key string) ([]string, error) {
	locale := d.locales[localeName]
	if locale == nil {
		return nil, fmt.Errorf("no such locale %s", localeName)
	}

	category := locale.categories[lc]
	if category == nil {
		return nil, fmt.Errorf("no such category %s in locale %s", lc, localeName)
	}

	if category.copiesFrom != "" {
		return d.Query(category.copiesFrom, lc, key)
	}

	operands, exists := category.values[key]
	if !exists {
		return nil, fmt.Errorf("no such key %s in category %s in locale %s", key, lc, localeName)
	}

	return operands, nil
}

// QueryString is a helper for keys which are a single string
func (d *Database) QueryString(localeName string, lc LC, key string) (string, error) {
	ops, err := d.Query(localeName, lc, key)
	if err != nil {
		return "", err
	}
	if len(ops) < 1 {
		return "", fmt.Errorf("key %s has no operands", key)
	}
	return ops[0], nil
}

// QueryInteger is a helper for keys which are a single integer
func (d *Database) QueryInteger(localeName string, lc LC, key string) (int, error) {
	op, err := d.QueryString(localeName, lc, key)
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(op)
	if err != nil {
		return 0, fmt.Errorf("key %s is not an integer", key)
	}
	return val, nil
}
