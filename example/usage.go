package main

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/race-conditioned/go-md/pkg/gomd"
)

func exampleUsage() {
	p := gomd.NewOnePassParser()
	md, err := gomd.Read("report.md")
	if err != nil {
		fmt.Println("Error reading markdown file:", err)
	}

	doc := p.Parse(string(md))
	for _, el := range doc.Elements {
		el.Text = strings.ToLower(el.Text)
		fmt.Println("EL!", el)
	}

	l := gomd.NewLexer()
	ctx := context.Background()
	tokens, err := l.TokenizeCtx(ctx, bytes.NewReader(md))
	if err != nil {
		fmt.Println("Error tokenizing markdown:", err)
	}
	fmt.Println("Tokens:")
	for _, tok := range tokens {
		fmt.Println(tok)
	}

	tokens2, err := l.Tokenize(bytes.NewReader(md))
	if err != nil {
		fmt.Println("Error tokenizing markdown:", err)
	}
	fmt.Println("Tokens2:")
	for _, tok := range tokens2 {
		fmt.Println(tok)
	}

	tp := gomd.NewTokenParser()
	ast, err := tp.ParseTokensCtx(ctx, tokens)
	if err != nil {
		fmt.Println("Error tokenizing markdown:", err)
	}
	for _, el := range ast.Elements {
		fmt.Println("AST Element:", el)
	}

	ast2, err := tp.ParseTokens(tokens2)
	if err != nil {
		fmt.Println("Error tokenizing markdown:", err)
	}
	for _, el := range ast2.Elements {
		fmt.Println("AST2 Element:", el)
	}

	BuildExample("build-example.md")
	CompoundExample("compound-example.md")
	BuilderMixExample("builder-mix-example.md")
	CompoundMixExample("compound-mix-example.md")
}

func BuildExample(filename string) {
	brand := "My Company"
	b := gomd.Builder{}

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

	body := []*gomd.Element{b.Text("This is the body")}

	template := []*gomd.Element{}
	template = append(template, header...)
	template = append(template, b.NL())
	template = append(template, body...)

	md := b.Build(template...)
	if err := gomd.Write(filename, md); err != nil {
		// handle error
	}
}

func CompoundExample(filename string) {
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
	if err := gomd.Write(filename, doc); err != nil {
		// handle error
	}
}

func examplefooter(comp string) []*gomd.Element {
	b := gomd.Builder{}
	return []*gomd.Element{
		b.Rule(),
		b.Textln(fmt.Sprintf("Copyright %s (c) 2025 Author. All Rights Reserved.", comp)),
	}
}

func BuilderMixExample(filename string) {
	comp := "My Company"
	b := gomd.Builder{}

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
			examplefooter(comp)..., // <- template slice spread right in
		)...,
	)

	if err := gomd.Write(filename, md); err != nil {
		// handle error
	}
}

func CompoundMixExample(filename string) {
	comp := "My Company"
	b := gomd.Builder{}
	c := gomd.Compounder{Builder: b}

	md := b.Build(
		c.Compound(
			c.Header1(fmt.Sprintf("%s Doc", comp)),
			c.Section2("Welcome", []string{
				fmt.Sprintf("This document is for %s.", comp),
			}),
			c.UL2("Departments", []string{"Ops", "Finance", "HR"}),
			examplefooter(comp), // <- builder helper dropped straight into Compounder
		)...,
	)

	if err := gomd.Write(filename, md); err != nil {
		// handle error
	}
}
