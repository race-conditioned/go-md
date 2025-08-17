package gomd

// section is a helper function which builds a Section with the given Header Builder function.
// If the title is empty, it will not render a header for the section.
func (c *Compounder) section(h func(string) *Element, title string, paras ...string) []*Element {
	out := []*Element{}
	if title != "" {
		out = append(out, h(title))
		out = append(out, c.Builder.NL())
	}

	for _, p := range paras {
		out = append(out, c.Builder.Textln(p), c.Builder.NL())
	}

	return out
}

// Section1 is used to render a markdown section with a H1 header and paragraphs of text.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) Section1(title string, paras []string) []*Element {
	return c.section(c.Builder.H1, title, paras...)
}

// Section2 is used to render a markdown section with a H2 header and paragraphs of text.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) Section2(title string, paras []string) []*Element {
	return c.section(c.Builder.H2, title, paras...)
}

// Section3 is used to render a markdown section with a H3 header and paragraphs of text.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) Section3(title string, paras []string) []*Element {
	return c.section(c.Builder.H3, title, paras...)
}

// Section4 is used to render a markdown section with a H4 header and paragraphs of text.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) Section4(title string, paras []string) []*Element {
	return c.section(c.Builder.H4, title, paras...)
}

// Section5 is used to render a markdown section with a H5 header and paragraphs of text.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) Section5(title string, paras []string) []*Element {
	return c.section(c.Builder.H5, title, paras...)
}

// Section6 is used to render a markdown section with a H6 header and paragraphs of text.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) Section6(title string, paras []string) []*Element {
	return c.section(c.Builder.H6, title, paras...)
}

// Header1 is used to render a markdown H1 header.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) Header1(text string) []*Element {
	return []*Element{c.Builder.H1(text), c.Builder.NL()}
}

// Header2 is used to render a markdown H2 header.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) Header2(text string) []*Element {
	return []*Element{c.Builder.H2(text), c.Builder.NL()}
}

// Header3 is used to render a markdown H3 header.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) Header3(text string) []*Element {
	return []*Element{c.Builder.H3(text), c.Builder.NL()}
}

// Header4 is used to render a markdown H4 header.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) Header4(text string) []*Element {
	return []*Element{c.Builder.H4(text), c.Builder.NL()}
}

// Header5 is used to render a markdown H5 header.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) Header5(text string) []*Element {
	return []*Element{c.Builder.H5(text), c.Builder.NL()}
}

// Header6 is used to render a markdown H6 header.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) Header6(text string) []*Element {
	return []*Element{c.Builder.H6(text), c.Builder.NL()}
}

// appendChildrenLines is a helper function that appends lines of text as children elements.
func (c *Compounder) appendChildrenLines(texts []string) []*Element {
	children := []*Element{}
	for _, text := range texts {
		children = append(children, c.Builder.Textln(text))
	}
	return children
}

// ul is a helper function which builds an unordered list with the given Header Builder function.
// If the title is empty, it will not render a header for the section.
func (c *Compounder) ul(h func(string) *Element, title string, texts []string) []*Element {
	out := []*Element{}
	if title != "" {
		out = append(out, h(title))
		out = append(out, c.Builder.NL())
	}
	return append(out, c.Builder.UL(c.appendChildrenLines(texts)...), c.Builder.NL())
}

// UL1 is used to render a markdown unordered list with a H1 header.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) UL1(title string, texts []string) []*Element {
	return c.ul(c.Builder.H1, title, texts)
}

// UL2 is used to render a markdown unordered list with a H2 header.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) UL2(title string, texts []string) []*Element {
	return c.ul(c.Builder.H2, title, texts)
}

// UL3 is used to render a markdown unordered list with a H3 header.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) UL3(title string, texts []string) []*Element {
	return c.ul(c.Builder.H3, title, texts)
}

// UL4 is used to render a markdown unordered list with a H4 header.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) UL4(title string, texts []string) []*Element {
	return c.ul(c.Builder.H4, title, texts)
}

// UL5 is used to render a markdown unordered list with a H5 header.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) UL5(title string, texts []string) []*Element {
	return c.ul(c.Builder.H5, title, texts)
}

// UL6 is used to render a markdown unordered list with a H6 header.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) UL6(title string, texts []string) []*Element {
	return c.ul(c.Builder.H6, title, texts)
}

// ol is a helper function which builds an ordered list with the given Header Builder function.
// If the title is empty, it will not render a header for the section.
func (c *Compounder) ol(h func(string) *Element, title string, texts []string) []*Element {
	out := []*Element{}
	if title != "" {
		out = append(out, h(title))
		out = append(out, c.Builder.NL())
	}
	return append(out, c.Builder.OL(c.appendChildrenLines(texts)...), c.Builder.NL())
}

// OL1 is used to render a markdown ordered list with a H1 header.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) OL1(title string, texts []string) []*Element {
	return c.ol(c.Builder.H1, title, texts)
}

// OL2 is used to render a markdown ordered list with a H2 header.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) OL2(title string, texts []string) []*Element {
	return c.ol(c.Builder.H2, title, texts)
}

// OL3 is used to render a markdown ordered list with a H3 header.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) OL3(title string, texts []string) []*Element {
	return c.ol(c.Builder.H3, title, texts)
}

// OL4 is used to render a markdown ordered list with a H4 header.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) OL4(title string, texts []string) []*Element {
	return c.ol(c.Builder.H4, title, texts)
}

// OL5 is used to render a markdown ordered list with a H5 header.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) OL5(title string, texts []string) []*Element {
	return c.ol(c.Builder.H5, title, texts)
}

// OL6 is used to render a markdown ordered list with a H6 header.
// If the title is empty, it will not render a header for the section.
// It returns a slice of pointers to an Element which can be used in the Compound function.
func (c *Compounder) OL6(title string, texts []string) []*Element {
	return c.ol(c.Builder.H6, title, texts)
}

// Compound is used to join compounder methods.
// It returns a slice of pointers to an Element which can be used in the Build function.
func (c *Compounder) Compound(groups ...[]*Element) []*Element {
	out := []*Element{}
	for _, g := range groups {
		out = append(out, g...)
	}
	return out
}
