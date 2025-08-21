package gomd

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRoundTrip(t *testing.T) {
	b := NewBuilder()
	p := NewOnePassParser()
	l := NewLexer()
	tp := NewTokenParser()

	cases := []struct {
		name string
		path string
	}{
		// HEADERS
		{"h1", "h1.md"},
		{"h2", "h2.md"},
		{"h3", "h3.md"},
		{"h4", "h4.md"},
		{"h5", "h5.md"},
		{"h6", "h6.md"},

		// TEXT
		{"text1", "text1.md"},

		// BOLD
		{"bold1", "bold1.md"},
		{"bold1ln", "bold1.md"},
		{"bold2", "bold2.md"},
		{"bold2ln", "bold2.md"},

		// ITALIC
		{"italic1", "italic1.md"},
		{"italic1ln", "italic1.md"},
		{"italic2", "italic2.md"},
		{"italic2ln", "italic2.md"},

		// LINK
		{"link1", "link1.md"},
		{"link1ln", "link1.md"},
		{"link2", "link2.md"},
		{"link2ln", "link2.md"},

		// IMAGE
		{"img", "img1.md"},
		{"img2", "img2.md"},

		// NL
		{"nl1", "nl1.md"},
		{"nl2", "nl2.md"},
		{"nl3", "nl2.md"},
		{"nl4", "nl2.md"},
		{"nl5", "nl5.md"},
		{"nl6", "nl5.md"},
		{"nl7", "nl7.md"},
		{"nl8", "nl1.md"},
		{"nl9", "nl9.md"},

		// Rule
		{"rule1", "rule1.md"},
		{"rule2", "rule2.md"},

		// Code
		{"code1", "code1.md"},
		{"code1ln", "code1.md"},
		{"code2", "code2.md"},
		{"code2ln", "code2.md"},

		// UL
		{"ul1", "nl1.md"},
		{"ul2", "ul2.md"},
		{"ul3", "ul3.md"},
		{"ul4", "ul4.md"},
		{"ul5", "ul3.md"},
		{"ul6", "ul6.md"},
		{"ul7", "ul7.md"},
		{"ul8", "ul8.md"},
		{"ul9", "ul9.md"},
		{"ul10", "ul10.md"},

		// OL
		{"ol1", "nl1.md"},
		{"ol2", "ol2.md"},
		{"ol3", "ol3.md"},
		{"ol4", "ol4.md"},
		{"ol5", "ol3.md"},
		{"ol6", "ol6.md"},
		{"ol7", "ol7.md"},
		{"ol8", "ol8.md"},
		{"ol9", "ol9.md"},
		{"ol10", "ol10.md"},
	}

	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		cmp.AllowUnexported(Element{}),
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			origBytes, err := Read("testdata/" + tc.path)
			if err != nil {
				t.Fatal(err)
			}
			orig := string(origBytes)

			// Route A: OneShotParser -> elements -> Build
			docA := p.Parse(orig)
			gotA := b.Build(docA.Elements...)

			// Route B: Tokenize -> ParseTokens -> elements -> Build
			toks, err := l.Tokenize(strings.NewReader(orig))
			if err != nil {
				t.Fatalf("Tokenize error: %v", err)
			}
			docB, err := tp.ParseTokens(toks)
			if err != nil {
				t.Fatalf("ParseTokens error: %v", err)
			}
			gotB := b.Build(docB.Elements...)

			// 1) Both routes should render identically
			if diff := cmp.Diff(gotA, gotB, opts...); diff != "" {
				t.Fatalf("Parse vs ParseTokens render mismatch (-Parse +ParseTokens):\n%s", diff)
			}

			// 2) And both should round-trip to the original source
			if diff := cmp.Diff(gotA, orig, opts...); diff != "" {
				t.Fatalf("Round trip via Parse mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(gotB, orig, opts...); diff != "" {
				t.Fatalf("Round trip via ParseTokens mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
