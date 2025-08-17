package gomd

import (
	"fmt"
	"strconv"
	"strings"
)

// H1 is used to render a h1 markdown header.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) H1(text string) *Element {
	return &Element{name: "h1", content: text + "\n"}
}

// H2 is used to render a h2 markdown header.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) H2(text string) *Element {
	return &Element{name: "h2", content: text + "\n"}
}

// H3 is used to render a h3 markdown header.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) H3(text string) *Element {
	return &Element{name: "h3", content: text + "\n"}
}

// H4 is used to render a h4 markdown header.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) H4(text string) *Element {
	return &Element{name: "h4", content: text + "\n"}
}

// H5 is used to render a h5 markdown header.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) H5(text string) *Element {
	return &Element{name: "h5", content: text + "\n"}
}

// H6 is used to render a h6 markdown header.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) H6(text string) *Element {
	return &Element{name: "h6", content: text + "\n"}
}

// Text is used to render text.
// If you want to end the text with a newline, use the Textln method on gomd.Builder.
// Text returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) Text(text string) *Element {
	return &Element{name: "text", content: text}
}

// Textln is used to render text that is concluded with a newline.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) Textln(text string) *Element {
	return &Element{name: "textln", content: text + "\n"}
}

// Bold is used to render bold text.
// If you want to end the bold text with a newline, use the Boldln method on gomd.Builder.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) Bold(text string) *Element {
	return &Element{name: "bold", content: fmt.Sprintf("**%s**", text)}
}

// Boldln is used to render bold text.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) Boldln(text string) *Element {
	return &Element{name: "boldln", content: fmt.Sprintf("**%s**\n", text)}
}

// Italic is used to render italic text.
// If you want to end the italic text with a newline, use the Italicln method on gomd.Builder.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) Italic(text string) *Element {
	return &Element{name: "italic", content: fmt.Sprintf("_%s_", text)}
}

// Italicln is used to render italic text.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) Italicln(text string) *Element {
	return &Element{name: "italicln", content: fmt.Sprintf("_%s_\n", text)}
}

// NL is used to render a newline.
// It behaves exactly as an escaped newline point in a string, eg. "\n".
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) NL() *Element {
	return &Element{name: "nl", content: "\n"}
}

// NLs is used to render multiple newlines.
// Specify the count of newlines you require.
// It behaves exactly as an escaped newline point in a string, eg. "\n".
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) NLs(count int) *Element {
	return &Element{name: "nl", content: strings.Repeat("\n", count)}
}

// UL is used to begin rendering an unordered list.
// It receives Element pointers as children.
// This allows for custom nesting.
// Any Element (including a UL Element) can be nested in a UL.
func (b *Builder) UL(children ...*Element) *Element {
	if len(children) > 0 {
		return &Element{name: "ul", children: children}
	}
	return nil
}

// OL is used to begin rendering an ordered list.
// It receives Element pointers as children.
// This allows for custom nesting.
// Any Element (including an OL Element) can be nested in a OL.
func (b *Builder) OL(children ...*Element) *Element {
	return &Element{name: "ol", children: children}
}

// Link is used to render a hyperlink.
// If you want to end the link with a newline, use the Linkln method on gomd.Builder.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) Link(display, link string) *Element {
	return &Element{name: "link", content: fmt.Sprintf("[%s](%s) ", display, link)}
}

// Linkln is used to render a hyperlink.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) Linkln(display, link string) *Element {
	return &Element{name: "linkln", content: fmt.Sprintf("[%s](%s)\n", display, link)}
}

// Img is used to render an image.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) Img(alt, link string) *Element {
	return &Element{name: "image", content: fmt.Sprintf("![%s](%s)\n", alt, link)}
}

// Rule is used to render a horizonatl rule.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) Rule() *Element {
	return &Element{name: "rule", content: "\n---\n\n"}
}

// Code is used to render an unfenced code block.
// If you want to end the code text with a newline, use the Codeln method on gomd.Builder.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) Code(text string) *Element {
	return &Element{name: "code", content: fmt.Sprintf("`%s`", text)}
}

// Codeln is used to render an unfenced code block.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) Codeln(text string) *Element {
	return &Element{name: "code", content: fmt.Sprintf("`%s`\n", text)}
}

type ParentInfo struct {
	iter                 int
	isPrefix             bool
	olParent             bool
	firstParent          string
	firstParentOLhandled bool
}

// crawlContent is used to parse an Element and recurse through the associated children Elements.
// It converts Element pointers into compatible markdown text and handles nesting.
func (b *Builder) crawlContent(el *Element, prefix string, parentInfo *ParentInfo, level int) {
	if el == nil {
		return
	}

	nextPrefix := false
	newIteration := 0
	outStr := ""
	olParent := false

	rootKind := parentInfo.firstParent
	if rootKind == "" && (el.name == "ul" || el.name == "ol") {
		rootKind = el.name
	}

	switch el.name {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		level, _ := strconv.Atoi(strings.TrimPrefix(el.name, "h"))
		outStr = strings.Repeat("#", level) + " " + el.content
		nextPrefix = true

	case "ul":
		if !parentInfo.olParent {
			if strings.Contains(prefix, "-") {
				prefix = "  " + prefix
			} else {
				prefix += "- "
			}
		}

	case "ol":
		newIteration++
		if strings.Contains(prefix, "-") {
			prefix = "  " + prefix
		} else {
			prefix += "- "
		}
		olParent = true

	case "textln", "boldln", "italicln", "linkln", "codeln":
		nextPrefix = true
		fallthrough
	default:
		outStr = el.content
	}

	outPrefix := prefix
	if parentInfo.iter != 0 && parentInfo.isPrefix && len(prefix) >= 2 {
		spaces := strings.Repeat(" ", len(prefix)-2)
		outPrefix = fmt.Sprintf("%s%d. ", spaces, parentInfo.iter)
		parentInfo.iter++
	}

	if rootKind == "ol" && !parentInfo.firstParentOLhandled && level == 3 {
		b.output.WriteString("\n")
		parentInfo.firstParentOLhandled = true
	}

	if outStr != "" {
		if parentInfo.isPrefix {
			outStr = outPrefix + outStr
		}
		b.output.WriteString(outStr)
	}
	parentInfo.isPrefix = nextPrefix

	childInfo := &ParentInfo{
		iter:                 newIteration,
		isPrefix:             true,
		olParent:             olParent,
		firstParent:          rootKind,
		firstParentOLhandled: parentInfo.firstParentOLhandled,
	}

	for _, child := range el.children {
		b.crawlContent(child, prefix, childInfo, level+1)
	}
}

// collapseRuns returns s with any run of '\n' longer than max reduced to max.
func collapseRuns(s string, max int) string {
	if max < 1 {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	run := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			if run < max {
				b.WriteByte('\n')
			}
			run++
		} else {
			run = 0
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// Build consumes Element pointers.
// It uses a recursive crawl function to convert each Element content into an equivalent in markdown.
func (b *Builder) Build(elements ...*Element) string {
	parentInfo := &ParentInfo{iter: 0, isPrefix: true, olParent: false, firstParent: ""}
	for _, el := range elements {
		if el == nil {
			continue
		}
		b.crawlContent(el, "", parentInfo, 1)
	}

	s := b.output.String()
	b.output.Reset()

	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.TrimLeft(s, "\n")
	s = collapseRuns(s, 2)
	s = strings.TrimRight(s, " \t")
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	if strings.HasSuffix(s, "\n\n") {
		s = s[:len(s)-1]
	}

	return s
}
