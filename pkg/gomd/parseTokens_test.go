package gomd

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func mustParse(t *testing.T, src string) []*Element {
	t.Helper()
	toks, err := Tokenize(strings.NewReader(src))
	if err != nil {
		t.Fatalf("Tokenize error: %v", err)
	}
	doc, err := ParseTokens(toks)
	if err != nil {
		t.Fatalf("ParseTokens error: %v", err)
	}
	return doc.Children
}

func assertElems(t *testing.T, got, want []*Element) {
	t.Helper()
	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		cmp.AllowUnexported(Element{}),
	}
	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Fatalf("AST mismatch (-want +got):\n%s", diff)
	}
}

func TestParseTokens_Heading(t *testing.T) {
	got := mustParse(t, "# Hello\n")
	want := []*Element{
		{Kind: KHeading, Level: 1, Text: "Hello", LineBreak: true},
	}
	assertElems(t, got, want)
}

func TestParseTokens_Paragraph_Simple(t *testing.T) {
	got := mustParse(t, "hello world\n")
	want := []*Element{
		{Kind: KText, Text: "hello world", LineBreak: true},
	}
	assertElems(t, got, want)
}

func TestParseTokens_Inline_Bold_Italic_Code_Link_Image(t *testing.T) {
	src := "**b** _i_ `c` [x](y) ![alt](img)\n"
	got := mustParse(t, src)
	// Bold uses "**...**" per Builder’s conventions.
	// We do NOT emit a separate " " after a link; the builder appends one automatically.
	want := []*Element{
		{Kind: KBold, Text: "**b**"},
		{Kind: KText, Text: " "},
		{Kind: KItalic, Text: "_i_"},
		{Kind: KText, Text: " "},
		{Kind: KCodeSpan, Text: "`c`"},
		{Kind: KText, Text: " "},
		{Kind: KLink, Text: "x", Href: "y"},
		{Kind: KImage, Alt: "alt", Href: "img", LineBreak: true}, // last inline in line gets LineBreak
	}
	assertElems(t, got, want)
}

func TestParseTokens_UnorderedList(t *testing.T) {
	got := mustParse(t, "- one\n- two\n")
	want := []*Element{
		{
			Kind: KList, ListKind: ListUnordered,
			Children: []*Element{
				{Kind: KText, Text: "one", LineBreak: true},
				{Kind: KText, Text: "two", LineBreak: true},
			},
		},
	}
	assertElems(t, got, want)
}

func TestParseTokens_OrderedList(t *testing.T) {
	got := mustParse(t, "1) one\n2. two\n")
	want := []*Element{
		{
			Kind: KList, ListKind: ListOrdered,
			Children: []*Element{
				{Kind: KText, Text: "one", LineBreak: true},
				{Kind: KText, Text: "two", LineBreak: true},
			},
		},
	}
	assertElems(t, got, want)
}

func TestParseTokens_Mixed_Blocks(t *testing.T) {
	src := "### Title\n- a\n1) b\npara\n"
	got := mustParse(t, src)
	want := []*Element{
		{Kind: KHeading, Level: 3, Text: "Title", LineBreak: true},
		{Kind: KList, ListKind: ListUnordered, Children: []*Element{
			{Kind: KText, Text: "a", LineBreak: true},
		}},
		{Kind: KList, ListKind: ListOrdered, Children: []*Element{
			{Kind: KText, Text: "b", LineBreak: true},
		}},
		{Kind: KText, Text: "para", LineBreak: true},
	}
	assertElems(t, got, want)
}
