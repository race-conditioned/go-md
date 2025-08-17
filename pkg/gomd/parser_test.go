package gomd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// TODO: handle concatenation
func TestSimpleParseCases(t *testing.T) {
	b := Builder{}
	cases := []struct {
		name string
		path string
		got  []*Element
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
		// {"italic1", "italic1.md", []*Element{b.Italic("hi")}},
		// {"italic1ln", "italic1.md", []*Element{b.Italicln("hi")}},
		// {"italic2", "italic2.md", []*Element{b.Italic("hi"), b.Text(","), b.Italicln("there")}},
		// {"italic2ln", "italic2.md", []*Element{b.Italic("hi"), b.Text(","), b.Italicln("there")}},
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
			want := ParseMD(string(md), "")
			if diff := cmp.Diff(want, tc.got, opts...); diff != "" {
				t.Fatalf("Build mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
