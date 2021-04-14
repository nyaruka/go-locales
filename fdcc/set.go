package fdcc

// Set is the complete FDCC document
type Set struct {
	Categories []*Category
}

// Category is a grouping of directives
type Category struct {
	Name string
	Body []*Line
}

// CopiesFrom returns the name of the set to copy this category from if the first line is the copy keyword
func (c *Category) CopiesFrom() string {
	if len(c.Body) == 1 && c.Body[0].Identifier == "copy" && len(c.Body[0].Operands) >= 1 {
		return c.Body[0].Operands[0]
	}
	return ""
}

// Line is a directive within a category
type Line struct {
	Identifier string
	Operands   []string
}
