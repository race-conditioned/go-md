package gomd

import (
	"os"
	"strings"
	"testing"
)

// Sinks to avoid dead-code elimination
var (
	sinkElems []*Element
	sinkDoc   *Document
	sinkStr   string
)

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

func dataset(t *testing.B) map[string]string {
	h3, err := Read("testdata/h3.md")
	if err != nil {
		t.Fatalf("read h3.md: %v", err)
	}
	ul10, err := Read("testdata/ul10.md")
	if err != nil {
		t.Fatalf("read h3.md: %v", err)
	}
	ol10, err := Read("testdata/ol10.md")
	if err != nil {
		t.Fatalf("read h3.md: %v", err)
	}
	ds := map[string]string{
		"h3":    string(h3),
		"ul10":  string(ul10),
		"ol10":  string(ol10),
		"mixed": "### Title\n- item\n1) my ordered item\npara _i_ **b** [x](y)\n",
		"large": strings.Repeat("## Head\n- a **bold** and _italic_\n1) link: [x](y)\n\n", 2000),
	}
	return ds
}

func BenchmarkOldParser_Parse(b *testing.B) {
	p := NewParser()
	for name, in := range dataset(b) {
		in := in // capture
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(in)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sinkElems = p.Parse(in)
				if len(sinkElems) == 0 && len(in) > 0 {
					b.Fatalf("empty result")
				}
			}
		})
	}
}

func BenchmarkPipeline_TokenizePlusParse(b *testing.B) {
	for name, in := range dataset(b) {
		in := in
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(in)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				toks, err := Tokenize(strings.NewReader(in))
				if err != nil {
					b.Fatal(err)
				}
				doc, err := ParseTokens(toks)
				if err != nil {
					b.Fatal(err)
				}
				sinkDoc = doc
				if sinkDoc == nil || len(doc.Children) == 0 && len(in) > 0 {
					b.Fatalf("empty doc")
				}
			}
		})
	}
}

func BenchmarkTokenize_Only(b *testing.B) {
	for name, in := range dataset(b) {
		in := in
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(in)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				toks, err := Tokenize(strings.NewReader(in))
				if err != nil {
					b.Fatal(err)
				}
				_ = toks
			}
		})
	}
}

func BenchmarkParseTokens_Only(b *testing.B) {
	for name, in := range dataset(b) {
		in := in
		// Pre-tokenize once outside the timed region.
		toks, err := Tokenize(strings.NewReader(in))
		if err != nil {
			b.Fatalf("pre-tokenize: %v", err)
		}
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(in)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				doc, err := ParseTokens(toks)
				if err != nil {
					b.Fatal(err)
				}
				sinkDoc = doc
			}
		})
	}
}

func BenchmarkEndToEnd_Build_OldParser(b *testing.B) {
	bldr := Builder{}
	p := NewParser()
	for name, in := range dataset(b) {
		in := in
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(in)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				elems := p.Parse(in)
				sinkStr = bldr.Build(elems...)
				if len(sinkStr) == 0 && len(in) > 0 {
					b.Fatalf("empty render")
				}
			}
		})
	}
}

func BenchmarkEndToEnd_Build_Pipeline(b *testing.B) {
	bldr := Builder{}
	for name, in := range dataset(b) {
		in := in
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(in)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				toks, err := Tokenize(strings.NewReader(in))
				if err != nil {
					b.Fatal(err)
				}
				doc, err := ParseTokens(toks)
				if err != nil {
					b.Fatal(err)
				}
				sinkStr = bldr.Build(doc.Children...)
				if len(sinkStr) == 0 && len(in) > 0 {
					b.Fatalf("empty render")
				}
			}
		})
	}
}
