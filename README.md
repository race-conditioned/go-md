# gomd

## Motivation

gomd is a markdown builder & parser in Go. It lets you create documents programmatically, and also parse/round-trip Markdown you already have.

Markdown has a loose grammar with lots of edge cases. gomd focuses on a pragmatic subset that’s stable and easy to round-trip.

This project is a WIP; early versions may have breaking changes.

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
els, err := p.ParseCtx(ctx, src, "")
if err != nil { /* handle */ }

md := gomd.Builder{}.Build(els...)

// non-context:
els2 := p.Parse(src, "")
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
 err := gomd.Write("my-company.md", md)
 if err != nil {
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
 err := gomd.Write("report.md", doc)
 if err != nil {
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

 err := gomd.Write("my-company.md", md)
 if err != nil {
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

 err := gomd.Write("my-company.md", md)
 if err != nil {
  // handle error
 }
}

```

## Features

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
