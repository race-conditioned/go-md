package gomd

import "strings"

var inlineReplacer = strings.NewReplacer(
	"*", "\\*", "_", "\\_",
	"[", "\\[", "]", "\\]",
	"(", "\\(", ")", "\\)",
	"~", "\\~", "`", "\\`",
)

// escapeInline escapes characters that have special meaning in inline markdown syntax.
func escapeInline(s string) string { return inlineReplacer.Replace(s) }

// escapeBackticks escapes backticks in a string, which is useful for inline code blocks.
func escapeBackticks(s string) string { return strings.ReplaceAll(s, "`", "\\`") }

// escapeLinkText escapes characters in link text that have special meaning in markdown.
func escapeLinkText(s string) string { return escapeInline(s) }

// escapeURL escapes spaces and parentheses in URLs, which is useful for markdown links.
func escapeURL(u string) string {
	r := strings.NewReplacer(" ", "%20", "(", "%28", ")", "%29")
	return r.Replace(u)
}
