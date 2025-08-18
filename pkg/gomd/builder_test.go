package gomd

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSimpleBuilderCases(t *testing.T) {
	b := Builder{}

	cases := []struct {
		name string
		path string
		got  string
	}{
		// HEADERS
		{"h1", "h1.md", b.Build(b.H1("Header test"))},
		{"h2", "h2.md", b.Build(b.H2("Header test"))},
		{"h3", "h3.md", b.Build(b.H3("Header test"))},
		{"h4", "h4.md", b.Build(b.H4("Header test"))},
		{"h5", "h5.md", b.Build(b.H5("Header test"))},
		{"h6", "h6.md", b.Build(b.H6("Header test"))},

		// TEXT
		{"text1", "text1.md", b.Build(b.Text("hi"))},
		{"textln", "text1.md", b.Build(b.Textln("hi"))},

		// BOLD
		{"bold1", "bold1.md", b.Build(b.Bold("hi"))},
		{"bold1ln", "bold1.md", b.Build(b.Boldln("hi"))},
		{"bold2", "bold2.md", b.Build(b.Bold("hi"), b.Text(","), b.Boldln("there"))},
		{"bold2ln", "bold2.md", b.Build(b.Bold("hi"), b.Text(","), b.Boldln("there"))},

		// ITALIC
		{"italic1", "italic1.md", b.Build(b.Italic("hi"))},
		{"italic1ln", "italic1.md", b.Build(b.Italicln("hi"))},
		{"italic2", "italic2.md", b.Build(b.Italic("hi"), b.Text(","), b.Italicln("there"))},
		{"italic2ln", "italic2.md", b.Build(b.Italic("hi"), b.Text(","), b.Italicln("there"))},

		// LINK
		{"link1", "link1.md", b.Build(b.Link("google", "https://google.com"))},
		{"link1ln", "link1.md", b.Build(b.Linkln("google", "https://google.com"))},
		{"link2", "link2.md", b.Build(b.Link("google", "https://google.com"), b.Linkln("amazon", "https://amazon.com"))},
		{"link2ln", "link2.md", b.Build(b.Link("google", "https://google.com"), b.Linkln("amazon", "https://amazon.com"))},

		// IMAGE
		{"img", "img1.md", b.Build(b.Img("alt", "https://google.com/img"))},
		{"img2", "img2.md", b.Build(b.Img("my-alt", "https://google.com/img"), b.Img("my-alt2", "https://amazon.com/img2"))},

		// NL
		{"nl1", "nl1.md", b.Build(b.NL())},
		{"nl2", "nl2.md", b.Build(b.NL(), b.Text("hi"), b.NL())},
		{"nl3", "nl2.md", b.Build(b.NL(), b.Text("hi"))},
		{"nl4", "nl2.md", b.Build(b.NL(), b.Textln("hi"))},
		{"nl5", "nl5.md", b.Build(b.Text("hi"), b.NL(), b.Text("there"))},
		{"nl6", "nl5.md", b.Build(b.Textln("hi"), b.Text("there"))},
		{"nl7", "nl7.md", b.Build(b.Textln("hi"), b.NL(), b.Text("there"))},
		{"nl9", "nl9.md", b.Build(b.Textln("hi"), b.NL(), b.Text("there"))},

		// Rule
		{"rule1", "rule1.md", b.Build(b.Rule())},
		{"rule2", "rule2.md", b.Build(b.Textln("hi"), b.Rule(), b.Textln("there"))},

		// Code
		{"code1", "code1.md", b.Build(b.Code("hi"))},
		{"code1ln", "code1.md", b.Build(b.Codeln("hi"))},
		{"code2", "code2.md", b.Build(b.Code("hi"), b.Text(","), b.Code("there"))},
		{"code2ln", "code2.md", b.Build(b.Code("hi"), b.Text(","), b.Codeln("there"))},

		// UL
		{"ul1", "nl1.md", b.Build(b.UL())},
		{"ul2", "ul2.md", b.Build(b.UL(b.Text("hi")))},
		{"ul3", "ul3.md", b.Build(b.UL(b.Text("hi"), b.Text(" there")))},
		{"ul4", "ul4.md", b.Build(b.UL(b.Text("hi"), b.Text("there")))},
		{"ul5", "ul3.md", b.Build(b.UL(b.Text("hi"), b.Textln(" there")))},
		{"ul6", "ul6.md", b.Build(b.UL(b.Text("hi "), b.Bold("there")))},
		{"ul7", "ul7.md", b.Build(b.UL(b.Textln("one"), b.Textln("two"), b.Textln("three")))},
		{"ul8", "ul8.md", b.Build(b.UL(b.Textln("one"), b.Text("two "), b.Textln("items"), b.Textln("three")))},
		{"ul9", "ul9.md", b.Build(b.UL(b.Textln("one"), b.Text("my link: "), b.Linkln("google", "google.com"), b.Textln("three")))},
		{"ul10", "ul10.md", b.Build(b.UL(b.Textln("one"), b.Text("my link: "), b.Link("google", "google.com"), b.Boldln("So Cool"), b.Textln("three")))},

		// // OL
		{"ol1", "nl1.md", b.Build(b.OL())},
		{"ol2", "ol2.md", b.Build(b.OL(b.Text("hi")))},
		{"ol3", "ol3.md", b.Build(b.OL(b.Text("hi"), b.Text(" there")))},
		{"ol4", "ol4.md", b.Build(b.OL(b.Text("hi"), b.Text("there")))},
		{"ol5", "ol3.md", b.Build(b.OL(b.Text("hi"), b.Textln(" there")))},
		{"ol6", "ol6.md", b.Build(b.OL(b.Text("hi "), b.Bold("there")))},
		{"ol7", "ol7.md", b.Build(b.OL(b.Textln("one"), b.Textln("two"), b.Textln("three")))},
		{"ol8", "ol8.md", b.Build(b.OL(b.Textln("one"), b.Text("two "), b.Textln("items"), b.Textln("three")))},
		{"ol9", "ol9.md", b.Build(b.OL(b.Textln("one"), b.Text("my link: "), b.Linkln("google", "google.com"), b.Textln("three")))},
		{"ol10", "ol10.md", b.Build(b.OL(b.Textln("one"), b.Text("my link: "), b.Link("google", "google.com"), b.Boldln("So Cool"), b.Textln("three")))},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			want, err := Read("testdata/" + tc.path)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(string(want), tc.got); diff != "" {
				t.Fatalf("Build mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCompoundCases(t *testing.T) {
	b := Builder{}
	c := Compounder{Builder: b}

	comp := "langfire"
	footer := footer(comp)

	cases := []struct {
		name string
		path string
		got  string
	}{
		{
			"complex", "complex.md", b.Build(
				append(
					[]*Element{
						b.H1(fmt.Sprintf("My %s Doc", comp)),
						b.NL(),
						b.Textln(fmt.Sprintf("Welcome to %s document", comp)),
						b.NL(),
						b.H5("Contents"),
						b.NL(),
						b.UL(
							b.Textln("Operations"),
							b.Textln("Finances"),
							b.Textln("HR"),
							b.Textln("Engineering"),
						),
					},
					footer...,
				)...,
			),
		},
		{
			"compound", "compound.md", b.Build(
				c.Compound(
					c.Section1("Title", []string{"para 1", "para 2", "para 3"}),
					c.Section2("Title 2", []string{"bara 1", "bara 2", "bara 3"}),
					c.OL1("", []string{"item 1", "item 2", "item 3"}),
					c.UL3("My Orders", []string{"order 1", "order 2", "order 3"}),
					footer,
				)...,
			),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			want, err := Read("testdata/compound/" + tc.path)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(string(want), tc.got); diff != "" {
				t.Fatalf("Build mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDeepNesting(t *testing.T) {
	b := Builder{}
	cases := []struct {
		name string
		path string
		got  string
	}{
		{
			"deep nesting", "nest.md", b.Build(
				b.UL(
					b.Textln("one"),
					b.Text("and "),
					b.Boldln("two"),
					b.Textln("and"),
					b.OL(
						b.Textln("first"),
						b.Bold("and"),
						b.Textln(" second"),
						b.UL(
							b.Textln("one more"),
							b.Text("and "),
							b.Boldln("two more"),
							b.Textln("and"),
							b.OL(
								b.Textln("first again"),
								b.Bold("and"),
								b.Textln(" second again"),
							),
						),
					),
				),
			),
		},
		{
			"deepnest2", "nest2.md", b.Build(
				b.OL(
					b.Textln("one"),
					b.Text("and "),
					b.Boldln("two"),
					b.Textln("and"),
					b.UL(
						b.Textln("first"),
						b.Bold("and"),
						b.Textln(" second"),
						b.OL(
							b.Textln("one more"),
							b.Text("and "),
							b.Boldln("two more"),
							b.Textln("and"),
							b.UL(
								b.Textln("first again"),
								b.Bold("and"),
								b.Textln(" second again"),
							),
						),
					),
				),
			),
		},
		// WARN: deep nesting is just not in scope, it's nice to have some cases
		{
			"deepnes3", "nest3.md", b.Build(
				b.OL(
					b.Textln("one"),
					b.Text("and "),
					b.Boldln("two"),
					b.Textln("and"),
					b.OL(
						b.Textln("first"),
						b.Bold("and"),
						b.Textln(" second"),
						b.OL(
							b.Textln("one more"),
							b.Text("and "),
							b.Boldln("two more"),
							b.Textln("and"),
							b.OL(
								b.Textln("first again"),
								b.Bold("and"),
								b.Textln(" second again"),
							),
						),
					),
				),
			),
		},
	}

	_ = Write("check.md", cases[2].got)

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			want, err := Read("testdata/nesting/" + tc.path)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(string(want), tc.got); diff != "" {
				t.Fatalf("Build mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TODO: fix edge cases (rule, colon in header)
func footer(comp string) []*Element {
	b := Builder{}
	return []*Element{
		b.Rule(),
		b.Textln(fmt.Sprintf("Copyright %s (c) 2025 Author. All Rights Reserved.", comp)),
	}
}
