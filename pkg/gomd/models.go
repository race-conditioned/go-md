package gomd

import (
	"context"
	"fmt"
)

// Document represents a complete markdown document.
type Document struct {
	Elements []*Element
}

// Element represents a single markdown element.
type Element struct {
	Kind      ElementKind
	Text      string
	LineBreak bool
	Level     int
	Href      string
	Alt       string
	ListKind  ListType
	Lang      string
	Children  []*Element
}

//go:generate stringer -type=ElementKind
type ElementKind uint8

// Kind represents the type of markdown element.
const (
	EKHeading ElementKind = iota
	EKText
	EKBold
	EKItalic
	EKCodeSpan
	EKCodeBlock
	EKNewLine
	EKRule
	EKLink
	EKImage
	EKList
	EKQuote
)

// ListType represents the type of list in markdown.
type ListType uint8

const (
	ListNone ListType = iota
	ListUnordered
	ListOrdered
)

// Builder is a simple markdown builder that accumulates markdown elements
type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

// Compounder is a struct that holds a Builder and provides methods to build markdown documents.
type Compounder struct {
	Builder Builder
}

func NewCompounder(b *Builder) *Compounder {
	return &Compounder{
		Builder: *b,
	}
}

// LEXING
type Lexer struct{}

func NewLexer() *Lexer {
	return &Lexer{}
}

//go:generate stringer -type=TokenKind
type TokenKind int

const (
	TText TokenKind = iota
	TStar
	TUnderscore
	TLBracket
	TRBracket
	TLParen
	TRParen
	TBacktick
	TBang
	TDash
	THash
	TNewline
	TOLMarker
	TEOF
)

type (
	// Pos represents a position in the source text.
	Pos struct{ Line, Col int }
	// Token represents a single token in the markdown source.
	Token struct {
		Kind   TokenKind
		Lexeme string
		Pos    Pos
	}
)

func (t Token) String() string {
	return fmt.Sprintf("%s(%q)@%d:%d", t.Kind, t.Lexeme, t.Pos.Line, t.Pos.Col)
}

// TokenParser

type TokenParser struct{}

func NewTokenParser() *TokenParser {
	return &TokenParser{}
}

// One Shot Parsing

// Parser is a 'one-step' Markdown parser that converts Markdown text into a slice of Elements.
type OnePassParser struct {
	text        string
	elements    []*Element
	leafNode    *[]*Element
	parentStack []*Element
	lineCtx     variableLineCtx
	ctx         context.Context
	err         error
}

// NewOnePassParser creates a new OnePassParser instance with initialized fields.
func NewOnePassParser() *OnePassParser {
	return &OnePassParser{
		elements:    []*Element{},
		leafNode:    &[]*Element{},
		parentStack: []*Element{},
		lineCtx: variableLineCtx{
			basePointer:      0,
			lookAheadPointer: 0,
			specialChars:     []indexChar{},
			ruleString:       "!*`[()]_",
			cache:            []byte{},
		},
	}
}

// variableLineCtx holds the context for parsing a single line of Markdown text.
type variableLineCtx struct {
	basePointer      int
	lookAheadPointer int
	specialChars     []indexChar
	ruleString       string
	cache            []byte
}

// indexChar is a helper struct to hold the index and character of special characters in the line.
// using this shortcuts seeking the next special character
type indexChar struct {
	i int
	c rune
}
