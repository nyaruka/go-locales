package fdcc

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// See http://std.dkuug.dk/jtc1/sc22/wg20/docs/n897-14652w25.pdf

// Parser is a parser for FDCC formatted documents
type Parser struct {
	tokenizer *Tokenizer
}

// NewParser creates a new parser from the given reader
func NewParser(src io.Reader) *Parser {
	return &Parser{
		tokenizer: NewTokenizer(NewInput(src)),
	}
}

func (p *Parser) Parse() (*Set, error) {
	set := &Set{}

	for {
		t, err := p.tokenizer.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if t.Type == TokenTypeNewLine {
			// ignore empty lines
		} else if strings.HasPrefix(t.Value, "LC_") {
			category, err := p.readCategory(t.Value)
			if err != nil {
				return nil, err
			}
			set.Categories = append(set.Categories, category)
		} else {
			if err := p.readPreCategoryStatement(t.Value); err != nil {
				return nil, err
			}
		}
	}

	return set, nil
}

func (p *Parser) readPreCategoryStatement(token string) error {
	operands := p.readOperands()

	switch token {
	case "comment_char":
		if len(operands) != 1 || len(operands[0]) != 1 {
			return fmt.Errorf("invalid comment_char operands: %v", operands)
		}
		p.tokenizer.CommentChar = []rune(operands[0])[0]
	case "escape_char":
		if len(operands) != 1 || len(operands[0]) != 1 {
			return fmt.Errorf("invalid escape_char operands: %v", operands)
		}
		p.tokenizer.EscapeChar = []rune(operands[0])[0]
	case "repertoiremap", "charmap":
		// TODO ?
	default:
		return fmt.Errorf("invalid precategory statement: %s", token)
	}
	return nil
}

func (p *Parser) readCategory(token string) (*Category, error) {
	category := &Category{Name: token}

	for {
		t, err := p.tokenizer.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("unexpected EOF reading category %s", category.Name)
		}
		if err != nil {
			return nil, err
		}

		if t.Value == "END" {
			p.tokenizer.Next() // read category name again after END
			break
		} else if t.Type == TokenTypeNewLine {
			// ignore empty lines
		} else {
			operands := p.readOperands()
			category.Body = append(category.Body, &Line{Identifier: t.Value, Operands: operands})
		}
	}

	return category, nil
}

func (p *Parser) readOperands() []string {
	operands := make([]string, 0, 5)
	for {
		t, _ := p.tokenizer.Next()
		if t == nil || t.Type == TokenTypeNewLine {
			break
		}

		operand := t.Value
		if t.Type == TokenTypeString {
			operand = operand[1 : len(operand)-1] // strip quotes
		}

		operands = append(operands, UnescapeUnicode(operand))
	}
	return operands
}

var regexUnicode = regexp.MustCompile(`<U[[:xdigit:]]{4}>`)

// UnescapeUnicode escapes unicode sequences of the form <U0000>
func UnescapeUnicode(s string) string {
	unescaped := regexUnicode.ReplaceAllStringFunc(s, func(m string) string {
		hex := m[2:6]
		code, _ := strconv.ParseInt(hex, 16, 32)
		return string(rune(code))
	})
	return unescaped
}
