package gomd

//go:generate stringer -type=Kind
type Kind uint8

// Kind represents the type of markdown element.
const (
	KHeading Kind = iota
	KText
	KBold
	KItalic
	KCodeSpan
	KCodeBlock
	KNewLine
	KRule
	KLink
	KImage
	KList
	KQuote
)

// ListType represents the type of list in markdown.
type ListType uint8

const (
	ListNone ListType = iota
	ListUnordered
	ListOrdered
)

// Document represents a complete markdown document.
type Document struct {
	Children []*Element
}

// Element represents a single markdown element.
type Element struct {
	Kind      Kind
	Text      string
	LineBreak bool
	Level     int
	Href      string
	Alt       string
	ListKind  ListType
	Lang      string
	Children  []*Element
}

// Builder is a simple markdown builder that accumulates markdown elements
type Builder struct{}

// Compounder is a struct that holds a Builder and provides methods to build markdown documents.
type Compounder struct {
	Builder Builder
}

// LEXING

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
