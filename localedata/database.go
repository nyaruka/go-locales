package localedata

import (
	"embed"
	"fmt"

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
