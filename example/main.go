package main

import (
	"fmt"

	"github.com/race-conditioned/go-md/pkg/gomd"
)

func header(title string) []*gomd.Element {
	b := gomd.NewBuilder()
	return []*gomd.Element{
		b.H1(title),
		b.NL(),
	}
}

func footer(comp string) []*gomd.Element {
	b := gomd.NewBuilder()
	return []*gomd.Element{
		b.Rule(),
		b.Textln(fmt.Sprintf("Copyright %s (c) 2025 Author. All Rights Reserved.", comp)),
	}
}

func installSection() []*gomd.Element {
	b := gomd.NewBuilder()
	return []*gomd.Element{
		b.H2("Install"),
		b.NL(),
		b.CodeFence("bash", "go get github.com/race-conditioned/go-md/pkg/gomd"),
		b.NL(),
	}
}

func whyThisExistsSection() []*gomd.Element {
	b := gomd.NewBuilder()
	return []*gomd.Element{
		b.H2("Why this exists"),
		b.NL(),
		b.Textln("A lot of Markdown libraries are either heavyweight, strictly spec-driven, or hard to round-trip. gomd aims to be:"),
		b.UL(
			b.Textln("Lightweight: small surface area, simple data model."),
			b.Textln("Fast: a one-pass parser for the common path, with snapshot benches below."),
			b.Textln("Practical: stable subset that round-trips well for programmatic generation and edits."),
			b.Textln("Cancellable: both lexer and parsers respect context cancel/timeout."),
			b.Textln("Tooling-friendly: an optional tokenize → parse pipeline with positions for editors/linters."),
		),
		b.NL(),
	}
}

func featureSetSection() []*gomd.Element {
	b := gomd.NewBuilder()
	return []*gomd.Element{
		b.H2("Feature set"),
		b.NL(),
		b.UL(
			b.Bold("Builder API"), b.Textln(" — headings (H1–H6), text, bold, italic, code spans, images, links, rules, lists (UL/OL), block quotes, fenced code blocks."),
			b.Bold("Compounder API"), b.Textln(" — ergonomic helpers for common sections and titled lists (e.g., Section2, UL3, OL2) that compose cleanly."),
			b.Bold("Render quality"), b.Textln(" — newline collapsing, whitespace trimming, predictable list prefixes/indentation."),
			b.Bold("Two parse routes"), b.Textln(" — (1) fast one-pass parser; (2) tokenize → parse pipeline with token positions."),
			b.Bold("Context support"), b.Textln(" — TokenizeCtx / ParseTokensCtx / ParseCtx honor cancellation and timeouts."),
			b.Bold("Round-trip"), b.Textln(" — builder ⇄ parser tests ensure stable text output for the supported subset."),
			b.Bold("Fuzz & benches"), b.Textln(" — fuzz tests for lexer round-trip; benchmarks for parsers and end-to-end build."),
			b.Bold("File I/O"), b.Textln(" — tiny helpers: Read(file), Write(file, text)."),
			b.Bold("Thread-friendly builder"), b.Textln(" — Builder is now pure/stateless (no internal buffers)."),
		),
		b.NL(),
	}
}

func benchmarksSection() []*gomd.Element {
	b := gomd.NewBuilder()
	return []*gomd.Element{
		b.H2("Benchmarks (snapshot)"),
		b.NL(),
		b.UL(
			// Old parser
			b.Textln("Old parser (ParseCtx)"),
			b.UL(
				b.Textln("h3: ~100 ns/op, 168 B, 3 allocs"),
				b.Textln("mixed: ~1.55 µs/op, 2.4 KB, 38 allocs"),
				b.Textln("large: ~3.06 ms/op, 4.06 MB, 60k allocs"),
			),
			// Pipeline
			b.Textln("Pipeline (Tokenize+ParseTokensCtx)"),
			b.UL(
				b.Textln("h3: ~1.37 µs/op, 4.9 KB, 12 allocs"),
				b.Textln("mixed: ~5.15 µs/op, 15 KB, 59 allocs"),
				b.Textln("large: ~9.57 ms/op, 27.8 MB, 78k allocs"),
			),
		),
		b.NL(),
		b.Italicln("Takeaway: old parser is ~3x faster and ~6–7x lower memory on large docs; gap is even bigger on tiny docs."),
		b.NL(),
		b.Italicln("Note: numbers vary by Go version/CPU; these are for relative shape, not absolute truth."),
	}
}

func parsingUsageSection() []*gomd.Element {
	b := gomd.NewBuilder()
	parts := []*gomd.Element{
		b.H2("Parsing: fast vs pipeline"),
		b.NL(),
		b.Textln("gomd supports two parse paths:"),
		// bullets
		b.UL(
			b.Bold("Fast one-pass parser"),
			b.Text(" → "),
			b.Code("ParseCtx"),
			b.Text(" (or "),
			b.Code("Parse"),
			b.Text(") — best when you just need "),
			b.Codeln("[]*Element"),

			b.Bold("Pipeline"),
			b.Text(" → "),
			b.Code("TokenizeCtx"),
			b.Text(" → "),
			b.Code("ParseTokensCtx"),
			b.Textln(". Heavier, but exposes tokens for tooling (highlighting, linting, conversions, etc.)."),
		),
		b.NL(),

		b.H3("Fast parser (with and without context)"),
		b.NL(),
		b.CodeFence("go",
			"p := gomd.NewOnePassParser()\n"+
				"b := gomd.NewBuilder()\n\n"+
				"ctx := context.Background()\n"+
				"els, err := p.ParseCtx(ctx, src)\n"+
				"if err != nil { /* handle */ }\n\n"+
				"md := b.Build(els...)\n\n"+
				"// non-context:\n"+
				"els2 := p.Parse(src)\n"+
				"md2  := b.Build(els2...)",
		),
		b.NL(),

		b.H3("Pipeline (tokens + token parser)"),
		b.NL(),
		b.CodeFence("go",
			"l := gomd.NewLexer()\n"+
				"tp := gomd.NewTokenParser()\n"+
				"b := gomd.NewBuilder()\n"+
				"ctx := context.Background()\n\n"+
				"toks, err := l.TokenizeCtx(ctx, bytes.NewReader(md))\n"+
				"if err != nil { /* handle */ }\n\n"+
				"doc, err := tp.ParseTokensCtx(ctx, toks)\n"+
				"if err != nil { /* handle */ }\n\n"+
				"md := b.Build(doc.Children...)",
		),
		b.NL(),
		b.Italicln("Rule of thumb: prefer the fast parser for speed and simpler apps; use the pipeline when you need tokens for tooling."),
	}
	return parts
}

func whyTokensSection() []*gomd.Element {
	b := gomd.NewBuilder()
	return []*gomd.Element{
		b.H2("Why tokens?"),
		b.NL(),
		b.Textln("Tokens unlock tooling that a one-pass parser can’t easily support:"),
		b.UL(
			b.Textln("Syntax highlighting / editor integrations."),
			b.Textln("Linters and formatters (detect spacing, missing markers, etc.)."),
			b.Textln("Precise error spans and clickable source locations."),
			b.Textln("Multiple back-ends: emit HTML/AST or transform docs."),
			b.Textln("Incremental parsing (reuse tokens between edits)."),
		),
		b.NL(),
		b.Italicln("If you don’t need any of that, stick to the fast parser."),
	}
}

func compatibilitySection() []*gomd.Element {
	b := gomd.NewBuilder()
	return []*gomd.Element{
		b.H2("Compatibility & limitations"),
		b.NL(),
		b.UL(
			b.Textln("Not full CommonMark — this is a pragmatic subset tuned for round-tripping."),
			b.Textln("Ordered-list markers: one-pass and pipeline aim for parity; multi-digit and \")\"/\".\" styles supported in the pipeline; one-pass focuses on the common case."),
			b.Textln("Deep list nesting: partial support (tracked in tests/roadmap)."),
			b.Textln("Horizontal rules: recognized as lines of dashes/spaces with ≥3 dashes."),
			b.Textln("Escaping: inline emphasis/code/link text/url escaping is pragmatic; edge cases may differ from strict CommonMark."),
		),
		b.NL(),
		b.Italicln("If you hit an edge case, please open an issue with a minimal repro."),
	}
}

func contributingSection() []*gomd.Element {
	b := gomd.NewBuilder()
	return []*gomd.Element{
		b.H2("Contributing"),
		b.NL(),
		b.Textln("PRs welcome! A few tips:"),
		b.UL(
			b.Text("Run tests: "), b.Codeln("go test ./pkg/gomd/... | ./pkg/bin/colorize"),
			b.Text("Run benches: "), b.Codeln("go test -bench=. -benchmem -run '^$' ./pkg/gomd/..."),
			b.Text("Regenerate token names (if TokenKind changes): "), b.Codeln("go generate ./pkg/gomd/..."),
			b.Textln("Add ctx-cancel/timeout tests for long loops (see *_cancel_test.go)."),
			b.Textln("Keep round-trip tests green (builder ⇄ parser ⇄ builder)."),
		),
		b.NL(),
		b.Textln("Open an issue to discuss bigger changes (block elements, CommonMark edges, etc.)."),
	}
}

var licenseMIT = `MIT License

Copyright (c) 2025 Author

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the “Software”), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.`

func licenseSection() []*gomd.Element {
	b := gomd.Builder{}
	return []*gomd.Element{
		b.H2("License"),
		b.NL(),
		b.Textln("Licensed under the MIT License (full text below)."),
		b.NL(),
		b.CodeFence("text", licenseMIT),
		b.NL(),
	}
}

var builderExample = `package main

import (
	"fmt"
	"github.com/race-conditioned/go-md/pkg/gomd"
)

func main() {
	brand := "My Company"
	b := gomd.NewBuilder()

	header := []*gomd.Element{
		b.H1(fmt.Sprintf("%s Document", brand)),
		b.NL(),
		b.Textln("great!"),
		b.NL(),
		b.UL(
			b.Textln("first"),
			b.Textln("second"),
			b.OL(
				b.Bold("first"),
				b.Textln(" element"),
			),
		),
	}

	body := []*gomd.Element{ b.Text("This is the body") }

	template := []*gomd.Element{}
	template = append(template, header...)
	template = append(template, b.NL())
	template = append(template, body...)

	md := b.Build(template...)
	if err := gomd.Write("my-company.md", md); err != nil {
		// handle error
	}
}`

var compounderExample = `package main

import (
	"github.com/race-conditioned/go-md/pkg/gomd"
)

func main() {
	b := gomd.NewBuilder()
	c := gomd.NewCompounder(b)

	doc := b.Build(
		c.Compound(
			c.Header1("Quarterly Report"),
			c.Section2("Summary", []string{
				"Revenue up 12%.",
				"Conversion improved.",
			}),
			c.UL2("Highlights", []string{"Ops", "Finance", "Eng"}),
		)...,
	)
	if err := gomd.Write("report.md", doc); err != nil {
		// handle error
	}
}`

var builderMixAndMatchExample = `package main

import (
	"fmt"
	"github.com/race-conditioned/go-md/pkg/gomd"
)

func footer(comp string) []*gomd.Element {
	b := gomd.NewBuilder()
	return []*gomd.Element{
		b.Rule(),
		b.Textln(fmt.Sprintf("Copyright %s (c) 2025 Author. All Rights Reserved.", comp)),
	}
}

func main() {
	comp := "My Company"
	b := gomd.NewBuilder()

	// Compose ad-hoc + template slices
	md := b.Build(
		append(
			[]*gomd.Element{
				b.H1(fmt.Sprintf("%s Doc", comp)),
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
			footer(comp)..., // <- template slice spread right in
		)...,
	)

	if err := gomd.Write("my-company.md", md); err != nil {
		// handle error
	}
}`

var compounderMixAndMatchExample = `package main

import (
	"fmt"
	"github.com/race-conditioned/go-md/pkg/gomd"
)

func footer(comp string) []*gomd.Element {
	b := gomd.NewBuilder
	return []*gomd.Element{
		b.Rule(),
		b.Textln(fmt.Sprintf("Copyright %s (c) 2025 Author. All Rights Reserved.", comp)),
	}
}

func main() {
	comp := "My Company"
	b := gomd.NewBuilder()
	c := gomd.NewCompounder(b)

	md := b.Build(
		c.Compound(
			c.Header1(fmt.Sprintf("%s Doc", comp)),
			c.Section2("Welcome", []string{
				fmt.Sprintf("This document is for %s.", comp),
			}),
			c.UL2("Departments", []string{"Ops", "Finance", "HR"}),
			footer(comp), // <- builder helper dropped straight into Compounder
		)...,
	)

	if err := gomd.Write("my-company.md", md); err != nil {
		// handle error
	}
}`

func usageSection(b gomd.Builder, c gomd.Compounder) []*gomd.Element {
	return c.Compound(
		c.Section2("Usage", []string{
			"There are two main ways to build markdown. The Compounder is ergonomic for simple docs; the Builder gives you full control.",
			"This README itself is produced with gomd.Compounder.",
		}),
		[]*gomd.Element{
			b.H3("Builder Example"),
			b.CodeFence("go", builderExample),
			b.NL(),
			b.H3("Compounder Example"),
			b.CodeFence("go", compounderExample),
			b.NL(),
		},
	)
}

func mixAndMatchSections(b gomd.Builder, c gomd.Compounder) []*gomd.Element {
	return c.Compound(
		c.Section2("Builder Mix & match templates", []string{
			"You can compose reusable templates (headers, footers, TOCs) and spread them directly into a Build(...) call.",
			"This enables programmatic generation & updates at scale (think 1,000+ docs).",
		}),
		[]*gomd.Element{
			b.CodeFence("go", builderMixAndMatchExample),
			b.NL(),
		},
		c.Section2("Compounder Mix & match templates", []string{
			"You can compose reusable templates (headers, footers, TOCs) and pass them to Compound(...) along with other sections.",
			"Compounder will flatten groups for you.",
		}),
		[]*gomd.Element{
			b.CodeFence("go", compounderMixAndMatchExample),
			b.NL(),
		},
	)
}

func main() {
	b := gomd.Builder{}
	c := gomd.Compounder{Builder: b}

	md := b.Build(
		c.Compound(
			// Header (template)
			header("gomd"),

			// Motivation
			c.Section2("Motivation", []string{
				"gomd is a markdown builder & parser in Go. It lets you create documents programmatically, and also parse/round-trip Markdown you already have.",
				"Markdown has a loose grammar with lots of edge cases. gomd focuses on a pragmatic subset that’s stable and easy to round-trip.",
				"This project is a WIP; early versions may have breaking changes.",
			}),
			whyThisExistsSection(),

			installSection(),

			// Two parsing routes (high-level, plain text)
			c.Section2("Two parsing routes", []string{
				"gomd supports two parse paths:",
				"1) Fast one-pass parser → ParseCtx (or Parse for back-compat). Best when you just need []*Element.",
				"2) Pipeline → TokenizeCtx → ParseTokensCtx. Heavier, but exposes tokens for tooling.",
			}),
			parsingUsageSection(),
			whyTokensSection(),

			featureSetSection(),
			compatibilitySection(),
			usageSection(b, c),
			mixAndMatchSections(b, c),

			// Features checklist (kept as a quick glance summary)
			c.UL2("At a glance", []string{
				"markdown builder ✅",
				"markdown compounder ✅",
				"Read and Write markdown ✅",
				"Basic Markdown syntax supported ✅",
				"Builder and composer tested ✅",
				"Round trip support ✅",
				"Dual parse routes (fast + pipeline) ✅",
				"Context cancellation in both routes ✅",
				"Deep nesting – partial support",
				"Serve to a viewer – planned",
				"Conversion to HTML – planned",
				"CommonMark compatibility – aspirational (longer-term)",
			}),

			benchmarksSection(),
			contributingSection(),

			[]*gomd.Element{
				b.H2("Running tests & benches"),
				b.CodeFence("bash",
					"# run tests (with colorized output via your awk script)\n"+
						"go test ./pkg/gomd/... | ./pkg/bin/colorize\n\n"+
						"# run benches\n"+
						"go test -bench=. -benchmem -run '^$' ./pkg/gomd/..."),
				b.NL(),
			},

			licenseSection(),
		)...,
	)

	const filename = "README.md"
	if err := gomd.Write(filename, md); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(fmt.Sprintf("wrote markdown to %s", filename))
}
