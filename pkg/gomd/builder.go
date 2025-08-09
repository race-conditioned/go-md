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
	return &Element{name: "italic", content: fmt.Sprintf("*%s*", text)}
}

// Italicln is used to render italic text.
// It returns a pointer to an Element which can be used in the Generate function.
func (b *Builder) Italicln(text string) *Element {
	return &Element{name: "italicln", content: fmt.Sprintf("*%s*\n", text)}
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
	return &Element{name: "ul", children: children}
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
	return &Element{name: "rule", content: "---\n"}
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

// crawlContent is used to parse an Element and recurse through the associated children Elements.
// It converts Element pointers into compatible markdown text and handles nesting.
func (b *Builder) crawlContent(el *Element, prefix string, iteration *struct{ num int }) {
	if el != nil {

		newIteration := &struct{ num int }{num: 0}
		outStr := ""

		switch el.name {
		case "h1", "h2", "h3", "h4", "h5":
			level, _ := strconv.Atoi(strings.Split(el.name, "")[1])
			hashes := strings.Repeat("#", level)
			outStr = hashes + " " + el.content
		case "ul":
			if strings.Contains(prefix, "-") {
				prefix = "  " + prefix
			} else {
				prefix += "- "
			}
		case "ol":
			newIteration.num++
			if strings.Contains(prefix, "-") {
				prefix = "  " + prefix
			} else {
				prefix += "- "
			}
		default:
			outStr = el.content
		}

		outPrefix := prefix
		if iteration.num != 0 {
			spaces := strings.Repeat(" ", len(prefix)-2)
			outPrefix = fmt.Sprintf("%s%d. ", spaces, iteration.num)
			iteration.num++
		}

		if len(outStr) != 0 {
			b.output.WriteString(outPrefix + outStr)
		}

		for _, child := range el.children {
			b.crawlContent(child, prefix, newIteration)
		}

		return
	}
}

// Generate consumes Element pointers.
// It uses a recursive crawl function to convert each Element content into an equivalent in markdown.
func (b *Builder) Generate(elements ...*Element) string {
	for _, el := range elements {
		iteration := &struct{ num int }{num: 0}
		b.crawlContent(el, "", iteration)
	}
	return b.output.String()
}
