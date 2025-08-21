package gomd

import (
	"strings"
)

// heading is the single source of truth for H1..H6.
func (b *Builder) heading(level int, text string) *Element {
	if level < 1 {
		level = 1
	} else if level > 6 {
		level = 6
	}
	return &Element{Kind: EKHeading, Level: level, LineBreak: true, Text: text}
}

// H1 returns an Element pointer representing a level-1 heading.
func (b *Builder) H1(text string) *Element { return b.heading(1, text) }

// H2 returns an Element pointer representing a level-2 heading.
func (b *Builder) H2(text string) *Element { return b.heading(2, text) }

// H3 returns an Element pointer representing a level-3 heading.
func (b *Builder) H3(text string) *Element { return b.heading(3, text) }

// H4 returns an Element pointer representing a level-4 heading.
func (b *Builder) H4(text string) *Element { return b.heading(4, text) }

// H5 returns an Element pointer representing a level-5 heading.
func (b *Builder) H5(text string) *Element { return b.heading(5, text) }

// H6 returns an Element pointer representing a level-6 heading.
func (b *Builder) H6(text string) *Element { return b.heading(6, text) }

// Text returns an Element pointer representing markdown text.
func (b *Builder) Text(text string) *Element { return &Element{Kind: EKText, Text: text} }

// Textln returns an Element pointer representing a markdown text followed by a newline character.
func (b *Builder) Textln(text string) *Element {
	return &Element{Kind: EKText, LineBreak: true, Text: text}
}

// Bold returns an Element pointer representing bold markdown text.
func (b *Builder) Bold(text string) *Element {
	return &Element{Kind: EKBold, Text: inlineWrap("**", escapeInline(text))}
}

// Boldln returns an Element pointer representing bold markdown text followed by a newline character.
func (b *Builder) Boldln(text string) *Element {
	return &Element{Kind: EKBold, LineBreak: true, Text: inlineWrap("**", escapeInline(text))}
}

// Italic returns an Element pointer representing italic markdown text.
func (b *Builder) Italic(text string) *Element {
	return &Element{Kind: EKItalic, Text: inlineWrap("_", escapeInline(text))}
}

// Italicln returns an Element pointer representing italic markdown text followed by a newline character.
func (b *Builder) Italicln(text string) *Element {
	return &Element{Kind: EKItalic, LineBreak: true, Text: inlineWrap("_", escapeInline(text))}
}

// Code returns an Element pointer representing markdown inline code (a code span). For Fenced blocks, use CodeBlock.
func (b *Builder) Code(text string) *Element {
	return &Element{Kind: EKCodeSpan, Text: inlineWrap("`", escapeBackticks(text))}
}

// Codeln returns an Element pointer representing markdown inline code (a code span) followed by a newline character. For Fenced blocks, use CodeBlock..
func (b *Builder) Codeln(text string) *Element {
	return &Element{Kind: EKCodeSpan, LineBreak: true, Text: inlineWrap("`", escapeBackticks(text))}
}

// CodeBlock returns an Element pointer representing a markdown fenced code block.
func (b *Builder) CodeBlock(lang, code string) *Element {
	return &Element{
		Kind: EKCodeBlock,
		Lang: lang,
		Text: code,
	}
}

// NL returns an Element pointer representing a markdown nl character. The Builder.Build method ignores all newlines beyond two sequentially.
func (b *Builder) NL() *Element { return &Element{Kind: EKNewLine, LineBreak: true} }

// UL returns an Element pointer representing the bounds of an unordered list.
// Element pointers can be passed as Children.
// This allows for custom nesting.
// Any Element (including a UL Element) can be nested in a UL.
func (b *Builder) UL(Children ...*Element) *Element {
	return &Element{Kind: EKList, ListKind: ListUnordered, Children: Children}
}

// OL returns an Element pointer representing the bounds of an ordered list.
// Element pointers can be passed as Children.
// This allows for custom nesting.
// Any Element (including an OL Element) can be nested in an OL.
func (b *Builder) OL(Children ...*Element) *Element {
	return &Element{Kind: EKList, ListKind: ListOrdered, Children: Children}
}

// Link returns an Element pointer representing a markdown link.
func (b *Builder) Link(display, link string) *Element {
	// INFO: trailing space is used to allow for spacing multiple links
	return &Element{Kind: EKLink, Text: escapeLinkText(display), Href: escapeURL(link)}
}

// Linkln returns an Element pointer representing a markdown link followed by a newline character.
func (b *Builder) Linkln(display, link string) *Element {
	return &Element{Kind: EKLink, LineBreak: true, Text: escapeLinkText(display), Href: escapeURL(link)}
}

// Img returns an Element pointer representing a markdown image followed by a newline character.
func (b *Builder) Img(alt, link string) *Element {
	return &Element{Kind: EKImage, LineBreak: true, Alt: alt, Href: link}
}

// Rule returns an Element pointer representing a markdown rule, it will always pad a full newline between other Text.
func (b *Builder) Rule() *Element {
	return &Element{Kind: EKRule, LineBreak: true, Text: "\n---\n"}
}

// CodeFence renders a fenced, optionally language-tagged block.
func (b *Builder) CodeFence(lang, code string) *Element { return b.CodeBlock(lang, code) }

// Quote renders a Markdown blockquote with Children.
func (b *Builder) Quote(Children ...*Element) *Element {
	return &Element{Kind: EKQuote, Children: Children}
}

func (b *Builder) cleanLastElement(elements []*Element) {
	if len(elements) == 0 {
		return
	}

	lastEl := elements[len(elements)-1]
	if !lastEl.LineBreak && lastEl.Kind != EKList {
		lastEl.LineBreak = true
	}
}

func (b *Builder) renderText(ctx *renderCtx, buf *strings.Builder, el *Element) {
	ctx.renderText(b, buf, el)
}

// Build consumes Element pointers.
// It uses a recursive render function to convert each Element Text into an equivalent in markdown.
func (b *Builder) Build(elements ...*Element) string {
	var buf strings.Builder
	ctx := &renderCtx{
		frames:      []listFrame{},
		lineBuffer:  &strings.Builder{},
		startOfLine: false,
	}

	b.cleanLastElement(elements)

	for _, el := range elements {
		if el == nil {
			continue
		}
		b.renderText(ctx, &buf, el)
	}

	return ctx.cleanRender(buf.String())
}
