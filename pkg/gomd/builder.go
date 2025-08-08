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

type Builder struct {
	config  BuilderConfig
	content []string
}

func (b *Builder) H1(text string) {
	b.content = append(b.content, "# "+text)
}

func (b *Builder) Text(text string) {
	b.content = append(b.content, text)
}

func (b *Builder) Generate() string {
	return strings.Join(b.content, "\n")
}
