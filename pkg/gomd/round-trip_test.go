package gomd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRoundTrip(t *testing.T) {
	b := Builder{}

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
			orig, err := Read("testdata/" + tc.path)
			if err != nil {
				t.Fatal(err)
			}
			want := b.Build(ParseMD(string(orig), "")...)
			if diff := cmp.Diff(want, string(orig), opts...); diff != "" {
				t.Fatalf("Round trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// // Property: Build -> Parse equals original AST (after nil/empty normalization).
// func FuzzRoundTripElements(f *testing.F) {
// 	f.Add(uint64(1)) // seeds for reproducibility
// 	f.Add(uint64(2))
// 	f.Fuzz(func(t *testing.T, seed uint64) {
// 		r := rand.New(rand.NewSource(seed))
// 		b := &Builder{}
// 		orig := genTree(r, 0)     // random elements
// 		md := b.Build(orig...) // under test
// 		got := ParseMD(md, "")    // round-trip
//
// 		opts := []cmp.Option{
// 			cmpopts.EquateEmpty(),
// 			cmp.AllowUnexported(Element{}),
// 		}
// 		if diff := cmp.Diff(normalize(orig), got, opts...); diff != "" {
// 			t.Fatalf("round-trip diff (-orig +parsed):\n%s\nmd:\n%s", diff, md)
// 		}
// 	})
// }
//
// func genTree(r *rand.Rand, depth int) []*Element {
// 	if depth > 3 { // limit depth
// 		return nil
// 	}
// 	n := r.Intn(3) + 1
// 	out := make([]*Element, 0, n)
// 	for i := 0; i < n; i++ {
// 		switch r.Intn(6) {
// 		case 0:
// 			out = append(out, (&Builder{}).H1(randText(r)))
// 		case 1:
// 			out = append(out, (&Builder{}).H2(randText(r)))
// 		case 2:
// 			// UL/OL with random children
// 			if r.Intn(2) == 0 {
// 				out = append(out, (&Builder{}).UL(genTree(r, depth+1)...))
// 			} else {
// 				out = append(out, (&Builder{}).OL(genTree(r, depth+1)...))
// 			}
// 		case 3:
// 			out = append(out, (&Builder{}).Text(randText(r)))
// 		case 4:
// 			out = append(out, (&Builder{}).Bold(randText(r)))
// 		case 5:
// 			out = append(out, (&Builder{}).Italic(randText(r)))
// 		}
// 	}
// 	return out
// }
//
// func randText(r *rand.Rand) string {
// 	words := []string{"hi", "hello", "great", "list", "item", "link", "msft", "amazon"}
// 	k := r.Intn(3) + 1
// 	var b strings.Builder
// 	for i := 0; i < k; i++ {
// 		if i > 0 {
// 			b.WriteByte(' ')
// 		}
// 		b.WriteString(words[r.Intn(len(words))])
// 	}
// 	return b.String()
// }
