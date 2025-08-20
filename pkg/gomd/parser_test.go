package gomd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func normalize(e []*Element) []*Element {
	out := make([]*Element, 0, len(e))
	for _, x := range e {
		if x == nil {
			continue
		}
		x.Children = normalize(x.Children)
		out = append(out, x)
	}
	return out
}

func TestSimpleParseCases(t *testing.T) {
	b := Builder{}
	p := NewParser()
	cases := []struct {
		name string
		path string
		want []*Element
	}{
		// HEADERS
		{"h1", "h1.md", []*Element{b.H1("Header test")}},
		{"h2", "h2.md", []*Element{b.H2("Header test")}},
		{"h3", "h3.md", []*Element{b.H3("Header test")}},
		{"h4", "h4.md", []*Element{b.H4("Header test")}},
		{"h5", "h5.md", []*Element{b.H5("Header test")}},
		{"h6", "h6.md", []*Element{b.H6("Header test")}},

		// TEXT
		{"textln", "text1.md", []*Element{b.Textln("hi")}},

		// BOLD
		{"bold1ln", "bold1.md", []*Element{b.Boldln("hi")}},
		{"bold2", "bold2.md", []*Element{b.Bold("hi"), b.Text(","), b.Boldln("there")}},

		// ITALIC
		{"italic1ln", "italic1.md", []*Element{b.Italicln("hi")}},
		{"italic2", "italic2.md", []*Element{b.Italic("hi"), b.Text(","), b.Italicln("there")}},

		// LINK
		{"link1ln", "link1.md", []*Element{b.Linkln("google", "https://google.com")}},
		{"link2", "link2.md", []*Element{b.Link("google", "https://google.com"), b.Linkln("amazon", "https://amazon.com")}},

		// IMAGE
		{"img", "img1.md", []*Element{b.Img("alt", "https://google.com/img")}},
		{"img2", "img2.md", []*Element{b.Img("my-alt", "https://google.com/img"), b.Img("my-alt2", "https://amazon.com/img2")}},

		// NL
		{"nl1", "nl1.md", []*Element{b.NL()}},
		{"nl2", "nl2.md", []*Element{b.Textln("hi")}},
		{"nl5", "nl5.md", []*Element{b.Textln("hi"), b.Textln("there")}},
		{"nl7", "nl7.md", []*Element{b.Textln("hi"), b.NL(), b.Textln("there")}},

		// RULE
		{"rule1", "rule1.md", []*Element{b.Rule()}},
		{"rule2", "rule2.md", []*Element{b.Textln("hi"), b.Rule(), b.Textln("there")}},

		// Code
		{"code1ln", "code1.md", []*Element{b.Codeln("hi")}},
		{"code2", "code2.md", []*Element{b.Code("hi"), b.Text(","), b.Codeln("there")}},

		// UL
		{"ul2", "ul2.md", []*Element{b.UL(b.Textln("hi"))}},
		{"ul6", "ul6.md", []*Element{b.UL(b.Text("hi "), b.Boldln("there"))}},
		{"ul7", "ul7.md", []*Element{b.UL(b.Textln("one"), b.Textln("two"), b.Textln("three"))}},
		{"ul9", "ul9.md", []*Element{b.UL(b.Textln("one"), b.Text("my link: "), b.Linkln("google", "google.com"), b.Textln("three"))}},
		{"ul10", "ul10.md", []*Element{b.UL(b.Textln("one"), b.Text("my link: "), b.Link("google", "google.com"), b.Boldln("So Cool"), b.Textln("three"))}},

		// OL
		{"ol2", "ol2.md", []*Element{b.OL(b.Textln("hi"))}},
		{"ol6", "ol6.md", []*Element{b.OL(b.Text("hi "), b.Boldln("there"))}},
		{"ol7", "ol7.md", []*Element{b.OL(b.Textln("one"), b.Textln("two"), b.Textln("three"))}},
		{"ol9", "ol9.md", []*Element{b.OL(b.Textln("one"), b.Text("my link: "), b.Linkln("google", "google.com"), b.Textln("three"))}},
		{"ol10", "ol10.md", []*Element{b.OL(b.Textln("one"), b.Text("my link: "), b.Link("google", "google.com"), b.Boldln("So Cool"), b.Textln("three"))}},
	}

	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		cmp.AllowUnexported(Element{}),
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			md, err := Read("testdata/" + tc.path)
			if err != nil {
				t.Fatal(err)
			}
			got := p.Parse(string(md))
			if diff := cmp.Diff(normalize(tc.want), got, opts...); diff != "" {
				t.Fatalf("Build mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
