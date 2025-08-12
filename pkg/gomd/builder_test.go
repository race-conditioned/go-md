package gomd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TODO: handle concatenation
func TestSimpleBuilderCases(t *testing.T) {
	b := Builder{}

	cases := []struct {
		name string
		path string
		got  string
	}{
		// HEADERS
		{"h1", "h1.md", b.Generate(b.H1("Header test"))},
		{"h2", "h2.md", b.Generate(b.H2("Header test"))},
		{"h3", "h3.md", b.Generate(b.H3("Header test"))},
		{"h4", "h4.md", b.Generate(b.H4("Header test"))},
		{"h5", "h5.md", b.Generate(b.H5("Header test"))},
		{"h6", "h6.md", b.Generate(b.H6("Header test"))},

		// TEXT
		{"text1", "text1.md", b.Generate(b.Text("hi"))},
		{"textln", "text1.md", b.Generate(b.Textln("hi"))},

		// BOLD
		{"bold1", "bold1.md", b.Generate(b.Bold("hi"))},
		{"bold1ln", "bold1.md", b.Generate(b.Boldln("hi"))},
		{"bold2", "bold2.md", b.Generate(b.Bold("hi"), b.Text(","), b.Boldln("there"))},
		{"bold2ln", "bold2.md", b.Generate(b.Bold("hi"), b.Text(","), b.Boldln("there"))},

		// ITALIC
		{"italic1", "italic1.md", b.Generate(b.Italic("hi"))},
		{"italic1ln", "italic1.md", b.Generate(b.Italicln("hi"))},
		{"italic2", "italic2.md", b.Generate(b.Italic("hi"), b.Text(","), b.Italicln("there"))},
		{"italic2ln", "italic2.md", b.Generate(b.Italic("hi"), b.Text(","), b.Italicln("there"))},

		// LINK
		{"link1", "link1.md", b.Generate(b.Link("google", "https://google.com"))},
		{"link1ln", "link1.md", b.Generate(b.Linkln("google", "https://google.com"))},
		{"link2", "link2.md", b.Generate(b.Link("google", "https://google.com"), b.Linkln("amazon", "https://amazon.com"))},
		{"link2ln", "link2.md", b.Generate(b.Link("google", "https://google.com"), b.Linkln("amazon", "https://amazon.com"))},

		// IMAGE
		{"img", "img1.md", b.Generate(b.Img("alt", "https://google.com/img"))},
		{"img2", "img2.md", b.Generate(b.Img("my-alt", "https://google.com/img"), b.Img("my-alt2", "https://amazon.com/img2"))},

		// NL
		{"nl1", "nl1.md", b.Generate(b.NL())},
		{"nl2", "nl2.md", b.Generate(b.NL(), b.Text("hi"), b.NL())},
		{"nl3", "nl2.md", b.Generate(b.NL(), b.Text("hi"))},
		{"nl4", "nl2.md", b.Generate(b.NL(), b.Textln("hi"))},
		{"nl5", "nl5.md", b.Generate(b.Text("hi"), b.NL(), b.Text("there"))},
		{"nl6", "nl5.md", b.Generate(b.Textln("hi"), b.Text("there"))},
		{"nl7", "nl7.md", b.Generate(b.Textln("hi"), b.NL(), b.Text("there"))},
		{"nl8", "nl1.md", b.Generate(b.NLs(8))},
		{"nl9", "nl9.md", b.Generate(b.Textln("hi"), b.NLs(8), b.Text("there"))},

		// Rule
		//{"rule1", "rule1.md", b.Generate(b.Rule())}, -- WARN: quite an edge case to just have a rule
		{"rule2", "rule2.md", b.Generate(b.Textln("hi"), b.Rule(), b.Textln("there"))},

		// Code
		{"code1", "code1.md", b.Generate(b.Code("hi"))},
		{"code1ln", "code1.md", b.Generate(b.Codeln("hi"))},
		{"code2", "code2.md", b.Generate(b.Code("hi"), b.Text(","), b.Code("there"))},
		{"code2ln", "code2.md", b.Generate(b.Code("hi"), b.Text(","), b.Codeln("there"))},

		// UL
		{"ul1", "nl1.md", b.Generate(b.UL())},
		{"ul2", "ul2.md", b.Generate(b.UL(b.Text("hi")))},
		{"ul3", "ul3.md", b.Generate(b.UL(b.Text("hi"), b.Text(" there")))},
		{"ul4", "ul4.md", b.Generate(b.UL(b.Text("hi"), b.Text("there")))},
		{"ul5", "ul3.md", b.Generate(b.UL(b.Text("hi"), b.Textln(" there")))},
		{"ul6", "ul6.md", b.Generate(b.UL(b.Text("hi "), b.Bold("there")))},
		{"ul7", "ul7.md", b.Generate(b.UL(b.Textln("one"), b.Textln("two"), b.Textln("three")))},
		{"ul8", "ul8.md", b.Generate(b.UL(b.Textln("one"), b.Text("two "), b.Textln("items"), b.Textln("three")))},
		{"ul9", "ul9.md", b.Generate(b.UL(b.Textln("one"), b.Text("my link: "), b.Linkln("google", "google.com"), b.Textln("three")))},
		{"ul10", "ul10.md", b.Generate(b.UL(b.Textln("one"), b.Text("my link: "), b.Link("google", "google.com"), b.Boldln("So Cool"), b.Textln("three")))},

		// OL
		{"ol1", "nl1.md", b.Generate(b.OL())},
		{"ol2", "ol2.md", b.Generate(b.OL(b.Text("hi")))},
		{"ol3", "ol3.md", b.Generate(b.OL(b.Text("hi"), b.Text(" there")))},
		{"ol4", "ol4.md", b.Generate(b.OL(b.Text("hi"), b.Text("there")))},
		{"ol5", "ol3.md", b.Generate(b.OL(b.Text("hi"), b.Textln(" there")))},
		{"ol6", "ol6.md", b.Generate(b.OL(b.Text("hi "), b.Bold("there")))},
		{"ol7", "ol7.md", b.Generate(b.OL(b.Textln("one"), b.Textln("two"), b.Textln("three")))},
		{"ol8", "ol8.md", b.Generate(b.OL(b.Textln("one"), b.Text("two "), b.Textln("items"), b.Textln("three")))},
		{"ol9", "ol9.md", b.Generate(b.OL(b.Textln("one"), b.Text("my link: "), b.Linkln("google", "google.com"), b.Textln("three")))},
		{"ol10", "ol10.md", b.Generate(b.OL(b.Textln("one"), b.Text("my link: "), b.Link("google", "google.com"), b.Boldln("So Cool"), b.Textln("three")))},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			want, err := LoadMD("testdata/" + tc.path)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(string(want), tc.got); diff != "" {
				t.Fatalf("Generate mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
