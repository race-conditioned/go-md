package gomd

import (
	"fmt"
	"strings"
)

type listFrame struct {
	kind  ListType
	index int
}

type renderCtx struct {
	frames      []listFrame
	lineBuffer  *strings.Builder
	startOfLine bool
}

func (ctx *renderCtx) lineBreak() {
	ctx.lineBuffer.WriteString("\n")
	ctx.startOfLine = true
}

func (ctx *renderCtx) pushFrame(kind ListType) {
	ctx.frames = append(ctx.frames, listFrame{kind: kind})
}
func (ctx *renderCtx) popFrame() { ctx.frames = ctx.frames[:len(ctx.frames)-1] }

func (ctx *renderCtx) computeIndent() (*listFrame, string) {
	var indent string
	f := &ctx.frames[len(ctx.frames)-1]

	if len(ctx.frames) > 1 {
		// Count parent ULs to determine visible indentation (2 spaces per UL level)
		countUL := 0
		for i, loopFrame := range ctx.frames {
			if i == len(ctx.frames)-1 {
				break
			}
			if loopFrame.kind == ListUnordered {
				countUL++
			}
		}
		indent = strings.Repeat("  ", countUL)

		// Add newline before a UL that immediately follows an OL to break visual block
		if len(ctx.frames) == 2 && ctx.frames[0].kind == ListOrdered && f.index == 0 {
			indent = "\n" + indent
		}

		// If the parent is an OL and we're a child OL starting at 0, inherit index to continue numbering
		pf := ctx.frames[len(ctx.frames)-2]
		if pf.kind == ListOrdered && f.index == 0 {
			f.index = pf.index
		}
	}
	return f, indent
}

func (ctx *renderCtx) listPrefix() string {
	if len(ctx.frames) == 0 {
		return ""
	}

	f, indent := ctx.computeIndent()

	switch f.kind {
	case ListUnordered:
		f.index++
		return indent + "- "
	case ListOrdered:
		f.index++
		return fmt.Sprintf("%s%d. ", indent, f.index)
	default:
		return ""
	}
}

// renderText is used to parse an Element and recurse through the associated Children Elements.
// It converts Element pointers into compatible markdown text and handles nesting.
func (ctx *renderCtx) renderText(b *Builder, el *Element) {
	if el == nil {
		return
	}

	switch el.Kind {
	case KHeading:
		ctx.lineBuffer.WriteString(strings.Repeat("#", el.Level) + " " + el.Text)
	case KList:
		ctx.pushFrame(el.ListKind)
		defer ctx.popFrame()
	case KCodeBlock:
		ctx.lineBuffer.WriteString(fmt.Sprintf("```%s\n%s\n\n```", el.Lang, el.Text))
	case KQuote:
		ctx.lineBuffer.WriteString("> ")
	case KLink:
		ctx.lineBuffer.WriteString(fmt.Sprintf("[%s](%s)", el.Text, el.Href))
		if !el.LineBreak {
			ctx.lineBuffer.WriteString(" ")
		}
	case KImage:
		ctx.lineBuffer.WriteString(fmt.Sprintf("![%s](%s)", el.Alt, el.Href))
	default:
		ctx.lineBuffer.WriteString(el.Text)
	}

	if el.LineBreak {
		if ctx.lineBuffer.String() != "" {
			ctx.lineBreak()
			b.output.WriteString(ctx.listPrefix() + ctx.lineBuffer.String())
		} else {
			ctx.lineBreak()
			b.output.WriteString(ctx.lineBuffer.String())
		}
		ctx.lineBuffer.Reset()
	}

	b.cleanLastElement(el.Children)
	for _, child := range el.Children {
		ctx.startOfLine = true
		ctx.renderText(b, child)
	}
}

// collapseRuns returns s with any run of '\n' longer than max reduced to max.
func (ctx *renderCtx) collapseRuns(s string, max int) string {
	if max < 1 {
		return s
	}
	if !strings.Contains(s, "\n\n") {
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

func (ctx *renderCtx) cleanRender(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.TrimLeft(s, "\n")
	s = ctx.collapseRuns(s, 2)
	s = strings.TrimRight(s, " \t")
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	if strings.HasSuffix(s, "\n\n") {
		s = s[:len(s)-1]
	}
	return s
}
