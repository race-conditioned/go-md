# gomd

## Motivation

gomd is a markdown builder & parser in Go. It lets you create documents programmatically, and also parse/round-trip Markdown you already have.

Markdown has a loose grammar with lots of edge cases. gomd focuses on a pragmatic subset that’s stable and easy to round-trip.

This project is a WIP; early versions may have breaking changes.

## Why this exists

A lot of Markdown libraries are either heavyweight, strictly spec-driven, or hard to round-trip. gomd aims to be:

- 🪶 Lightweight: small surface area, simple data model.
- ⚡ Fast: a one-pass parser for the common path, with snapshot benches below.
- 🧠 Practical: stable subset that round-trips well for programmatic generation and edits.
- 🛑 Cancellable: both lexer and parsers respect context cancel/timeout.
- 🔧 Tooling-friendly: an optional tokenize → parse pipeline with positions for editors/linters.

## Install

```bash
go get github.com/race-conditioned/go-md/pkg/gomd

```

## Two parsing routes

gomd supports two parse paths:

1. Fast one-pass parser → ParseCtx (or Parse for back-compat). Best when you just need []\*Element.

2. Pipeline → TokenizeCtx → ParseTokensCtx. Heavier, but exposes tokens for tooling.

## Parsing: fast vs pipeline

gomd supports two parse paths:

- **Fast one-pass parser** → `ParseCtx` (or `Parse`) — best when you just need `[]*Element`
- **Pipeline** → `TokenizeCtx` → `ParseTokensCtx`. Heavier, but exposes tokens for tooling (highlighting, linting, conversions, etc.).

### Fast parser (with and without context)

```go
p := gomd.NewParser()
b := gomd.Builder{}

ctx := context.Background()
els, err := p.ParseCtx(ctx, src)
if err != nil { /* handle */ }

md := gomd.Builder{}.Build(els...)

// non-context:
els2 := p.Parse(src)
md2  := b.Build(els2...)

```

### Pipeline (tokens + token parser)

```go
ctx := context.Background()

toks, err := gomd.TokenizeCtx(ctx, strings.NewReader(src))
if err != nil { /* handle */ }

doc, err := gomd.ParseTokensCtx(ctx, toks)
if err != nil { /* handle */ }

md := gomd.Builder{}.Build(doc.Children...)

```

_Rule of thumb: prefer the fast parser for speed and simpler apps; use the pipeline when you need tokens for tooling._

## Why tokens?

Tokens unlock tooling that a one-pass parser can’t easily support:

- Syntax highlighting / editor integrations.
- Linters and formatters (detect spacing, missing markers, etc.).
- Precise error spans and clickable source locations.
- Multiple back-ends: emit HTML/AST or transform docs.
- Incremental parsing (reuse tokens between edits).

_If you don’t need any of that, stick to the fast parser._

## Feature set

- **Builder API** — headings (H1–H6), text, bold, italic, code spans, images, links, rules, lists (UL/OL), block quotes, fenced code blocks.
- **Compounder API** — ergonomic helpers for common sections and titled lists (e.g., Section2, UL3, OL2) that compose cleanly.
- **Render quality** — newline collapsing, whitespace trimming, predictable list prefixes/indentation.
- **Two parse routes** — (1) fast one-pass parser; (2) tokenize → parse pipeline with token positions.
- **Context support** — TokenizeCtx / ParseTokensCtx / ParseCtx honor cancellation and timeouts.
- **Round-trip** — builder ⇄ parser tests ensure stable text output for the supported subset.
- **Fuzz & benches** — fuzz tests for lexer round-trip; benchmarks for parsers and end-to-end build.
- **File I/O** — tiny helpers: Read(file), Write(file, text).
- **Thread-friendly builder** — Builder is now pure/stateless (no internal buffers).

## Compatibility & limitations

- Not full CommonMark — this is a pragmatic subset tuned for round-tripping.
- Ordered-list markers: one-pass and pipeline aim for parity; multi-digit and ")"/"." styles supported in the pipeline; one-pass focuses on the common case.
- Deep list nesting: partial support (tracked in tests/roadmap).
- Horizontal rules: recognized as lines of dashes/spaces with ≥3 dashes.
- Escaping: inline emphasis/code/link text/url escaping is pragmatic; edge cases may differ from strict CommonMark.

_If you hit an edge case, please open an issue with a minimal repro._

## Usage

There are two main ways to build markdown. The Compounder is ergonomic for simple docs; the Builder gives you full control.

This README itself is produced with gomd.Compounder.

### Builder Example

```go
package main

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
 if err := gomd.Write("my-company.md", md); err != nil {
  // handle error
 }
}

```

### Compounder Example

```go
package main

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
 if err := gomd.Write("report.md", doc); err != nil {
  // handle error
 }
}

```

## Builder Mix & match templates

You can compose reusable templates (headers, footers, TOCs) and spread them directly into a Build(...) call.

This enables programmatic generation & updates at scale (think 1,000+ docs).

```go
package main

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

 if err := gomd.Write("my-company.md", md); err != nil {
  // handle error
 }
}

```

## Compounder Mix & match templates

You can compose reusable templates (headers, footers, TOCs) and pass them to Compound(...) along with other sections.

Compounder will flatten groups for you.

```go
package main

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

 if err := gomd.Write("my-company.md", md); err != nil {
  // handle error
 }
}

```

## At a glance

- markdown builder ✅
- markdown compounder ✅
- Read and Write markdown ✅
- Basic Markdown syntax supported ✅
- Builder and composer tested ✅
- Round trip support ✅
- Dual parse routes (fast + pipeline) ✅
- Context cancellation in both routes ✅
- Deep nesting – partial support
- Serve to a viewer – planned
- Conversion to HTML – planned
- CommonMark compatibility – aspirational (longer-term)

## Benchmarks (snapshot)

- Old parser (ParseCtx)
  - h3: ~100 ns/op, 168 B, 3 allocs
  - mixed: ~1.55 µs/op, 2.4 KB, 38 allocs
  - large: ~3.06 ms/op, 4.06 MB, 60k allocs
- Pipeline (Tokenize+ParseTokensCtx)
  - h3: ~1.37 µs/op, 4.9 KB, 12 allocs
  - mixed: ~5.15 µs/op, 15 KB, 59 allocs
  - large: ~9.57 ms/op, 27.8 MB, 78k allocs

_Takeaway: old parser is \~3x faster and \~6–7x lower memory on large docs; gap is even bigger on tiny docs._

_Note: numbers vary by Go version/CPU; these are for relative shape, not absolute truth._

## Contributing

PRs welcome! A few tips:

- Run tests: `go test ./pkg/gomd/... | ./pkg/bin/colorize`
- Run benches: `go test -bench=. -benchmem -run '^$' ./pkg/gomd/...`
- Regenerate token names (if TokenKind changes): `go generate ./pkg/gomd/...`
- Add ctx-cancel/timeout tests for long loops (see \*\_cancel_test.go).
- Keep round-trip tests green (builder ⇄ parser ⇄ builder).

Open an issue to discuss bigger changes (block elements, CommonMark edges, etc.).

## Running tests & benches

```bash
# run tests (with colorized output via your awk script)
go test ./pkg/gomd/... | ./pkg/bin/colorize

# run benches
go test -bench=. -benchmem -run '^$' ./pkg/gomd/...

```

## License

Licensed under the MIT License (full text below).

```text
MIT License

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
THE SOFTWARE.

```
