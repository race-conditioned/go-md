package main

import (
	"fmt"

	"github.com/race-conditioned/go-md/pkg/gomd"
)

func header(title string) []*gomd.Element {
	b := gomd.Builder{}
	return []*gomd.Element{
		b.H1(title),
		b.NL(),
	}
}

func footer(comp string) []*gomd.Element {
	b := gomd.Builder{}
	return []*gomd.Element{
		b.Rule(),
		b.Textln(fmt.Sprintf("Copyright %s (c) 2025 Author. All Rights Reserved.", comp)),
	}
}

func benchmarksSection() []*gomd.Element {
	b := gomd.Builder{}
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
	}
}

func parsingUsageSection() []*gomd.Element {
	b := gomd.Builder{}
	parts := []*gomd.Element{
		b.H2("Parsing: fast vs pipeline"),
		b.NL(),
		b.Textln("gomd supports two parse paths:"),
		// bullet 1
		b.UL(
			b.Bold("Fast one-pass parser"),
			b.Text(" → "),
			b.Code("ParseCtx"),
			b.Text(" (or "),
			b.Code("Parse"),
			b.Text(") — best when you just need "),
			b.Codeln("[]*Element"),

			// bullet 2
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
			"p := gomd.NewParser()\n"+
				"b := gomd.Builder{}\n\n"+
				"ctx := context.Background()\n"+
				"els, err := p.ParseCtx(ctx, src, \"\")\n"+
				"if err != nil { /* handle */ }\n\n"+
				"md := gomd.Builder{}.Build(els...)\n\n"+
				"// non-context:\n"+
				"els2 := p.Parse(src, \"\")\n"+
				"md2  := b.Build(els2...)",
		),
		b.NL(),

		b.H3("Pipeline (tokens + token parser)"),
		b.NL(),
		b.CodeFence("go",
			"ctx := context.Background()\n\n"+
				"toks, err := gomd.TokenizeCtx(ctx, strings.NewReader(src))\n"+
				"if err != nil { /* handle */ }\n\n"+
				"doc, err := gomd.ParseTokensCtx(ctx, toks)\n"+
				"if err != nil { /* handle */ }\n\n"+
				"md := gomd.Builder{}.Build(doc.Children...)",
		),
		b.NL(),
		b.Italicln("Rule of thumb: prefer the fast parser for speed and simpler apps; use the pipeline when you need tokens for tooling."),
	}
	return parts
}

func whyTokensSection() []*gomd.Element {
	b := gomd.Builder{}
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

func contributingSection() []*gomd.Element {
	b := gomd.Builder{}
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

var builderExample = `package main

import (
	"fmt"
	"github.com/race-conditioned/go-md/pkg/gomd"
)

func main() {
	brand := "My Company"
	b := gomd.Builder{}

	header := []*gomd.Element{
		b.H1(fmt.Sprintf("My %s Document", brand)),
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
	err := gomd.Write("my-company.md", md)
	if err != nil {
		// handle error
	}
}`

var compounderExample = `package main

import (
	"github.com/race-conditioned/go-md/pkg/gomd"
)

func main() {
	b := gomd.Builder{}
	c := gomd.Compounder{Builder: b}

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
	err := gomd.Write("report.md", doc)
	if err != nil {
		// handle error
	}
}`

var builderMixAndMatchExample = `package main

import (
	"fmt"
	"github.com/race-conditioned/go-md/pkg/gomd"
)

func footer(comp string) []*gomd.Element {
	b := gomd.Builder{}
	return []*gomd.Element{
		b.Rule(),
		b.Textln(fmt.Sprintf("Copyright %s (c) 2025 Author. All Rights Reserved.", comp)),
	}
}

func main() {
	comp := "My Company"
	b := gomd.Builder{}

	// Compose ad-hoc + template slices
	md := b.Build(
		append(
			[]*gomd.Element{
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
			footer(comp)..., // <- template slice spread right in
		)...,
	)

	err := gomd.Write("my-company.md", md)
	if err != nil {
		// handle error
	}
}`

var compounderMixAndMatchExample = `package main

import (
	"fmt"
	"github.com/race-conditioned/go-md/pkg/gomd"
)

func footer(comp string) []*gomd.Element {
	b := gomd.Builder{}
	return []*gomd.Element{
		b.Rule(),
		b.Textln(fmt.Sprintf("Copyright %s (c) 2025 Author. All Rights Reserved.", comp)),
	}
}

func main() {
	comp := "My Company"
	b := gomd.Builder{}
	c := gomd.Compounder{Builder: b}

	md := b.Build(
		c.Compound(
			c.Header1(fmt.Sprintf("My %s Doc", comp)),
			c.Section2("Welcome", []string{
				fmt.Sprintf("This document is for %s.", comp),
			}),
			c.UL2("Departments", []string{"Ops", "Finance", "HR"}),
			footer(comp), // <- builder helper dropped straight into Compounder
		)...,
	)

	err := gomd.Write("my-company.md", md)
	if err != nil {
		// handle error
	}
}`

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

			// Two parsing routes (high-level, plain text)
			c.Section2("Two parsing routes", []string{
				"gomd supports two parse paths:",
				"1) Fast one-pass parser → ParseCtx (or Parse for back-compat). Best when you just need []*Element.",
				"2) Pipeline → TokenizeCtx → ParseTokensCtx. Heavier, but exposes tokens for tooling.",
			}),

			// Practical usage with code fences, and token rationale
			parsingUsageSection(),
			whyTokensSection(),

			// Usage (Builder + Compounder) + real code fences
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

			// Builder mix & match templates, rendered with CodeFence
			c.Section2("Builder Mix & match templates", []string{
				"You can compose reusable templates (headers, footers, TOCs) and spread them directly into a Build(...) call.",
				"This enables programmatic generation & updates at scale (think 1,000+ docs).",
			}),
			[]*gomd.Element{
				b.CodeFence("go", builderMixAndMatchExample),
				b.NL(),
			},

			// Compounder mix & match templates, rendered with CodeFence
			c.Section2("Compounder Mix & match templates", []string{
				"You can compose reusable templates (headers, footers, TOCs) and pass them to Compound(...) along with other sections.",
				"Compounder will flatten groups for you.",
			}),
			[]*gomd.Element{
				b.CodeFence("go", compounderMixAndMatchExample),
				b.NL(),
			},

			// Features checklist
			c.UL2("Features", []string{
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
			},
		)...,
	)

	const filename = "README.md"
	if err := gomd.Write(filename, md); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(fmt.Sprintf("wrote markdown to %s", filename))
}
