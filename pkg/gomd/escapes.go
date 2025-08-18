package gomd

import "strings"

var inlineReplacer = strings.NewReplacer(
	"*", "\\*", "_", "\\_",
	"[", "\\[", "]", "\\]",
	"(", "\\(", ")", "\\)",
	"~", "\\~", "`", "\\`",
)

func escapeInline(s string) string    { return inlineReplacer.Replace(s) }
func escapeBackticks(s string) string { return strings.ReplaceAll(s, "`", "\\`") }
func escapeLinkText(s string) string  { return escapeInline(s) }
func escapeURL(u string) string {
	r := strings.NewReplacer(" ", "%20", "(", "%28", ")", "%29")
	return r.Replace(u)
}
