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
		x.children = normalize(x.children)
		out = append(out, x)
	}
	return out
}

func TestRoundTrip(t *testing.T) {
	b := Builder{}

	orig := []*Element{
		b.H1("Header test"),
	}
	md := b.Build(orig...)
	got := ParseMD(md, "")

	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		cmp.AllowUnexported(Element{}),
	}

	if diff := cmp.Diff(normalize(orig), got, opts...); diff != "" {
		t.Fatalf("AST mismatch after round-trip (-want +got):\n%s", diff)
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
