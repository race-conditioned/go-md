package gomd

import "strings"

type Kind uint8

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

type ListType uint8

const (
	ListNone ListType = iota
	ListUnordered
	ListOrdered
)

type MarkDownFormat string

// // TODO: support formats
// const (
// 	Slack   MarkDownFormat = "slack"
// 	Discord MarkDownFormat = "discord"
// )
//
// type BuilderConfig struct {
// 	format MarkDownFormat
// }

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

// INFO: example refactor model
// type Element struct {
// 	Kind     Kind
// 	Text     string   // for text, code span, heading
// 	Level    int      // for headings 1..6
// 	Href     string   // link/image
// 	Alt      string   // image alt
// 	Lang     string   // code block language
// 	ListKind ListType // list kind
// 	Children []*Element
// }

type Builder struct {
	//	config BuilderConfig
	output *strings.Builder
}

type Compounder struct {
	Builder Builder
}
