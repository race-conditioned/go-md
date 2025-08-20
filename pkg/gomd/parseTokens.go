package gomd

import (
	"context"
	"strings"
)

// Parser is a Markdown parser that converts lexed tokens into a slice of Elements in a Document.
func ParseTokens(tks []Token) (*Document, error) {
	return ParseTokensCtx(context.Background(), tks)
}

// ParseTokensCtx parses lexed tokens into a Document, respecting the context for cancellation or timeout.
func ParseTokensCtx(ctx context.Context, tks []Token) (*Document, error) {
	var out []*Element

	i := 0
	bol := true // beginning of line
	var currentList *Element
	var currentListKind ListType

	checkCtx := func() error {
		if err := ctx.Err(); err != nil {
			return err
		}
		return nil
	}

	for i < len(tks) {
		if err := checkCtx(); err != nil {
			return &Document{Children: out}, err
		}
		if tks[i].Kind == TEOF {
			break
		}

		// preserve blank lines, newline at BOL means emit an empty line element.
		if tks[i].Kind == TNewline {
			if bol {
				out = append(out, &Element{Kind: KText, Text: "", LineBreak: true})
			}
			bol = true
			i++
			continue
		}

		// block dispatch at BOL
		if bol {
			// horizontal rule: line of only '-' and spaces, with >=3 dashes.
			if ok, next := isHorizontalRuleLine(tks, i); ok {
				// close any open list
				currentList = nil
				out = append(out, &Element{Kind: KRule, Text: "\n---\n", LineBreak: true})
				i = next // consume entire line (including its trailing newline if present)
				bol = true
				continue
			}

			// heading: THash+ then rest of line as text
			if tks[i].Kind == THash {
				level := 0
				for i < len(tks) && tks[i].Kind == THash {
					level++
					i++
					if err := checkCtx(); err != nil {
						return &Document{Children: out}, err
					}
				}
				if level > 6 {
					level = 6
				}
				text := collectUntilNewline(tks, i)
				text = strings.TrimLeft(text, " ")
				// advance to end-of-line, but don't consume the newline itself
				for i < len(tks) && tks[i].Kind != TNewline && tks[i].Kind != TEOF {
					i++
				}
				out = append(out, &Element{
					Kind:      KHeading,
					Level:     level,
					Text:      text,
					LineBreak: true,
				})
				bol = false // <-- INFO: prevent next newline from being treated as a blank line
				continue
			}

			// ordered list item: OL marker at BOL
			if tks[i].Kind == TOLMarker {
				if currentList == nil || currentListKind != ListOrdered {
					currentList = &Element{Kind: KList, ListKind: ListOrdered, Children: []*Element{}}
					out = append(out, currentList)
					currentListKind = ListOrdered
				}
				i++                                                     // consume marker
				elems, ni, err := parseInlineLineCtx(ctx, tks, i, true) // trim one leading space after marker
				if err != nil {
					return &Document{Children: out}, err
				}
				i = ni
				currentList.Children = append(currentList.Children, elems...)
				bol = true
				continue
			}

			// unordered list item: '-' + space at BOL
			if tks[i].Kind == TDash && i+1 < len(tks) && tks[i+1].Kind == TText && strings.HasPrefix(tks[i+1].Lexeme, " ") {
				if currentList == nil || currentListKind != ListUnordered {
					currentList = &Element{Kind: KList, ListKind: ListUnordered, Children: []*Element{}}
					out = append(out, currentList)
					currentListKind = ListUnordered
				}
				i++                                                     // consume '-'
				elems, ni, err := parseInlineLineCtx(ctx, tks, i, true) // drop a single leading space
				if err != nil {
					return &Document{Children: out}, err
				}
				i = ni
				currentList.Children = append(currentList.Children, elems...)
				bol = true
				continue
			}

			// not a list marker: close any open list
			currentList = nil

			// plain line
			elems, ni, err := parseInlineLineCtx(ctx, tks, i, false)
			if err != nil {
				return &Document{Children: out}, err
			}
			i = ni
			out = append(out, elems...)
			bol = true
			continue
		}

		// not at BOL (rare): treat as plain line until newline
		elems, ni, err := parseInlineLineCtx(ctx, tks, i, false)
		if err != nil {
			return &Document{Children: out}, err
		}
		i = ni
		out = append(out, elems...)
		bol = true
	}

	return &Document{Children: out}, nil
}

// collectUntilNewline collects tokens until a newline or EOF is encountered.
func collectUntilNewline(tks []Token, i int) string {
	var b strings.Builder
	for i < len(tks) && tks[i].Kind != TNewline && tks[i].Kind != TEOF {
		b.WriteString(tks[i].Lexeme)
		i++
	}
	return b.String()
}

// onlySpaces checks if the string contains only spaces or tabs.
func onlySpaces(s string) bool {
	for _, r := range s {
		if r != ' ' && r != '\t' {
			return false
		}
	}
	return true
}

// checks if the current line is a horizontal rule (>=3 '-' and only spaces otherwise).
// returns (ok, nextIndexAfterThisLine)
func isHorizontalRuleLine(tks []Token, i int) (bool, int) {
	j := i
	dashes := 0
	valid := true
	for j < len(tks) && tks[j].Kind != TNewline && tks[j].Kind != TEOF {
		switch tks[j].Kind {
		case TDash:
			dashes++
		case TText:
			if !onlySpaces(tks[j].Lexeme) {
				valid = false
			}
		default:
			valid = false
		}
		if !valid {
			break
		}
		j++
	}
	if !valid || dashes < 3 {
		return false, i
	}
	// consume the trailing newline if present
	if j < len(tks) && tks[j].Kind == TNewline {
		j++
	}
	return true, j
}

// parse a single logical line into inline Elements.
// If trimLeadingSpace is true, drop exactly one leading space in the first TText.
func parseInlineLineCtx(ctx context.Context, tks []Token, i int, trimLeadingSpace bool) ([]*Element, int, error) {
	var out []*Element
	var buf strings.Builder
	flushText := func(linebreak bool) {
		if buf.Len() == 0 {
			return
		}
		out = append(out, &Element{Kind: KText, Text: buf.String(), LineBreak: linebreak})
		buf.Reset()
	}

	first := true
	lastWasLink := false

	checkCtx := func() error {
		if err := ctx.Err(); err != nil {
			return err
		}
		return nil
	}

loop:
	for i < len(tks) {
		if err := checkCtx(); err != nil {
			return out, i, err
		}

		t := tks[i]
		switch t.Kind {
		case TNewline, TEOF:
			flushText(false)
			break loop

		case TBacktick:
			// `code`
			if i+2 < len(tks) && tks[i+1].Kind == TText && tks[i+2].Kind == TBacktick {
				flushText(false)
				code := tks[i+1].Lexeme
				out = append(out, &Element{Kind: KCodeSpan, Text: inlineWrap("`", escapeBackticks(code))})
				i += 3
				first = false
				lastWasLink = false
				continue
			}
			buf.WriteString(t.Lexeme)
			i++
			first = false
			lastWasLink = false

		case TStar:
			// **bold** only; single '*' treated as literal
			if i+4 < len(tks) && tks[i+1].Kind == TStar && tks[i+2].Kind == TText && tks[i+3].Kind == TStar && tks[i+4].Kind == TStar {
				flushText(false)
				inner := tks[i+2].Lexeme
				out = append(out, &Element{Kind: KBold, Text: inlineWrap("**", escapeInline(inner))})
				i += 5
				first = false
				lastWasLink = false
				continue
			}
			// literal '*'
			buf.WriteString(t.Lexeme)
			i++
			first = false
			lastWasLink = false

		case TUnderscore:
			// _italic_
			if i+2 < len(tks) && tks[i+1].Kind == TText && tks[i+2].Kind == TUnderscore {
				flushText(false)
				inner := tks[i+1].Lexeme
				out = append(out, &Element{Kind: KItalic, Text: inlineWrap("_", escapeInline(inner))})
				i += 3
				first = false
				lastWasLink = false
				continue
			}
			buf.WriteString(t.Lexeme)
			i++
			first = false
			lastWasLink = false

		case TBang:
			// ![alt](src)
			if i+6 < len(tks) &&
				tks[i+1].Kind == TLBracket &&
				tks[i+2].Kind == TText &&
				tks[i+3].Kind == TRBracket &&
				tks[i+4].Kind == TLParen &&
				tks[i+5].Kind == TText &&
				tks[i+6].Kind == TRParen {
				flushText(false)
				alt := tks[i+2].Lexeme
				src := tks[i+5].Lexeme
				out = append(out, &Element{Kind: KImage, Alt: alt, Href: escapeURL(src)})
				i += 7
				first = false
				lastWasLink = false
				continue
			}
			buf.WriteString(t.Lexeme)
			i++
			first = false
			lastWasLink = false

		case TLBracket:
			// [text](href)
			if i+5 < len(tks) &&
				tks[i+1].Kind == TText &&
				tks[i+2].Kind == TRBracket &&
				tks[i+3].Kind == TLParen &&
				tks[i+4].Kind == TText &&
				tks[i+5].Kind == TRParen {
				flushText(false)
				text := escapeLinkText(tks[i+1].Lexeme)
				href := escapeURL(tks[i+4].Lexeme)
				out = append(out, &Element{Kind: KLink, Text: text, Href: href})
				i += 6
				first = false
				lastWasLink = true
				continue
			}
			buf.WriteString(t.Lexeme)
			i++
			first = false
			lastWasLink = false

		default:
			// TText or punctuation
			lit := t.Lexeme
			if first && trimLeadingSpace && t.Kind == TText && strings.HasPrefix(lit, " ") {
				lit = lit[1:] // drop one leading space after "-" or "1)"/"2."
			}
			// avoid double spacing: builder already appends one space after links (when !LineBreak).
			if lastWasLink && t.Kind == TText && lit == " " {
				i++
				first = false
				lastWasLink = false
				continue
			}
			buf.WriteString(lit)
			i++
			first = false
			lastWasLink = false
		}
	}

	// empty line
	if len(out) == 0 && buf.Len() == 0 {
		out = append(out, &Element{Kind: KText, Text: "", LineBreak: true})
		return out, i + btoi(i < len(tks) && tks[i].Kind == TNewline), nil
	}

	if buf.Len() > 0 {
		out = append(out, &Element{Kind: KText, Text: buf.String()})
	}
	// mark only the last as a line break
	out[len(out)-1].LineBreak = true
	return out, i + btoi(i < len(tks) && tks[i].Kind == TNewline), nil
}
