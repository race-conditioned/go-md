package gomd

import "strings"

type MarkDownFormat string

const (
	Slack   MarkDownFormat = "slack"
	Discord MarkDownFormat = "discord"
)

type BuilderConfig struct {
	format MarkDownFormat
}

type Element struct {
	name     string
	content  string
	children []*Element
}

type Builder struct {
	config BuilderConfig
	output strings.Builder
}

type Compounder struct {
	Builder Builder
}
