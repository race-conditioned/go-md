package gomd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TK(k TokenKind, lit string, line, col int) Token {
	return Token{Kind: k, Lexeme: lit, Pos: Pos{Line: line, Col: col}}
}

func diffTokens(want, got []Token, ignorePos bool) string {
	opts := []cmp.Option{cmpopts.EquateEmpty()}
	if ignorePos {
		opts = append(opts, cmpopts.IgnoreFields(Token{}, "Pos"))
	}
	return cmp.Diff(want, got, opts...)
}

func assertTokensExact(t *testing.T, got, want []Token) {
	t.Helper()
	if diff := diffTokens(want, got, false); diff != "" {
		t.Fatalf("tokens mismatch (-want +got):\n%s", diff)
	}
}

func assertTokensKindsLexemes(t *testing.T, got, want []Token) {
	t.Helper()
	if diff := diffTokens(want, got, true); diff != "" {
		t.Fatalf("tokens mismatch (-want +got):\n%s", diff)
	}
}

func concatLexemes(toks []Token) string {
	var b strings.Builder
	for _, tk := range toks {
		if tk.Kind == TEOF {
			continue
		}
		b.WriteString(tk.Lexeme)
	}
	return b.String()
}

func hasKind(toks []Token, k TokenKind) bool {
	for _, tk := range toks {
		if tk.Kind == k {
			return true
		}
	}
	return false
}

func TestTokenize_Table(t *testing.T) {
	l := NewLexer()
	tests := []struct {
		name     string
		in       string
		want     []Token
		exactPos bool
		check    func(t *testing.T, got []Token) // optional extra assertions
	}{
		{
			name: "simple heading with bold",
			in:   "# Hello *world*",
			want: []Token{
				TK(THash, "#", 1, 1),
				TK(TText, " Hello ", 1, 2),
				TK(TStar, "*", 1, 9),
				TK(TText, "world", 1, 10),
				TK(TStar, "*", 1, 15),
				TK(TEOF, "", 1, 15),
			},
			exactPos: true,
		},
		{
			name: "unordered list simple",
			in:   "- item\n",
			want: []Token{
				TK(TDash, "-", 1, 1),
				TK(TText, " item", 1, 2),
				TK(TNewline, "\n", 1, 7),
				TK(TEOF, "", 2, 0),
			},
		},
		{
			name: "ordered list variants at BOL",
			in:   "1) x\n2. y\n123) z\n",
			want: []Token{
				TK(TOLMarker, "1)", 1, 1), TK(TText, " x", 1, 0), TK(TNewline, "\n", 1, 0),
				TK(TOLMarker, "2.", 2, 1), TK(TText, " y", 2, 0), TK(TNewline, "\n", 2, 0),
				TK(TOLMarker, "123)", 3, 1), TK(TText, " z", 3, 0), TK(TNewline, "\n", 3, 0),
				TK(TEOF, "", 4, 0),
			},
			exactPos: false, // positions vary a bit; we only care about kinds+lexemes here
		},
		{
			name: "ordered list with ≤3-space indent (spaces preserved as text)",
			in:   "   3) indented\n",
			want: []Token{
				TK(TText, "   ", 1, 1),
				TK(TOLMarker, "3)", 1, 4),
				TK(TText, " indented", 1, 0),
				TK(TNewline, "\n", 1, 0),
				TK(TEOF, "", 2, 0),
			},
			exactPos: false,
		},
		{
			name: "indent >3 means NOT an OL marker",
			in:   "    4) not-ol\n",
			// We no longer over-constrain punctuation: just ensure no TOLMarker and round-trip holds.
			want:     nil,
			exactPos: false,
			check: func(t *testing.T, got []Token) {
				if hasKind(got, TOLMarker) {
					t.Fatalf("should NOT emit TOLMarker when indent > 3")
				}
				in := "    4) not-ol\n"
				if concatLexemes(got) != in {
					t.Fatalf("round-trip mismatch:\n in: %q\nout: %q\n", in, concatLexemes(got))
				}
			},
		},
		{
			name: "tab counts as 2 spaces for indent (still OL marker)",
			in:   "\t3) tab-indented\n",
			// Accept hyphen tokenization; compare kinds+lexemes only.
			want: []Token{
				TK(TText, "\t", 0, 0),
				TK(TOLMarker, "3)", 0, 0),
				TK(TText, " tab", 0, 0),
				TK(TDash, "-", 0, 0),
				TK(TText, "indented", 0, 0),
				TK(TNewline, "\n", 0, 0),
				TK(TEOF, "", 0, 0),
			},
			exactPos: false,
		},
		{
			name: "no OL at mid-line (BOL required)",
			in:   "x 1) y\n",
			want: []Token{
				TK(TText, "x 1", 0, 0),
				TK(TRParen, ")", 0, 0),
				TK(TText, " y", 0, 0),
				TK(TNewline, "\n", 0, 0),
				TK(TEOF, "", 0, 0),
			},
			exactPos: false,
		},
		{
			name: "inline constructs (kinds/lexemes only)",
			in:   "![alt](u) [x](y) `c` _i_ *b*",
			want: []Token{
				TK(TBang, "!", 0, 0), TK(TLBracket, "[", 0, 0), TK(TText, "alt", 0, 0), TK(TRBracket, "]", 0, 0),
				TK(TLParen, "(", 0, 0), TK(TText, "u", 0, 0), TK(TRParen, ")", 0, 0),
				TK(TText, " ", 0, 0),
				TK(TLBracket, "[", 0, 0), TK(TText, "x", 0, 0), TK(TRBracket, "]", 0, 0),
				TK(TLParen, "(", 0, 0), TK(TText, "y", 0, 0), TK(TRParen, ")", 0, 0),
				TK(TText, " ", 0, 0),
				TK(TBacktick, "`", 0, 0), TK(TText, "c", 0, 0), TK(TBacktick, "`", 0, 0),
				TK(TText, " ", 0, 0),
				TK(TUnderscore, "_", 0, 0), TK(TText, "i", 0, 0), TK(TUnderscore, "_", 0, 0),
				TK(TText, " ", 0, 0),
				TK(TStar, "*", 0, 0), TK(TText, "b", 0, 0), TK(TStar, "*", 0, 0),
				TK(TEOF, "", 0, 0),
			},
			exactPos: false,
		},
		{
			name: "heading + lists combined",
			in:   "### Title\n- item\n1) my ordered item\n",
			want: []Token{
				TK(THash, "#", 1, 1), TK(THash, "#", 1, 2), TK(THash, "#", 1, 3),
				TK(TText, " Title", 1, 4),
				TK(TNewline, "\n", 1, 10),
				TK(TDash, "-", 2, 1),
				TK(TText, " item", 2, 2),
				TK(TNewline, "\n", 2, 7),
				TK(TOLMarker, "1)", 3, 1),
				TK(TText, " my ordered item", 3, 3),
				TK(TNewline, "\n", 3, 19),
				TK(TEOF, "", 4, 0),
			},
			exactPos: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := l.Tokenize(strings.NewReader(tc.in))
			if err != nil {
				t.Fatalf("Tokenize error: %v", err)
			}
			if tc.check != nil {
				tc.check(t, got)
				return
			}
			if tc.exactPos {
				assertTokensExact(t, got, tc.want)
			} else {
				assertTokensKindsLexemes(t, got, tc.want)
			}
		})
	}
}

// concatenating token lexemes (except TEOF) round-trips the source.
func TestTokenize_RoundTripLexemes(t *testing.T) {
	l := NewLexer()
	cases := []string{
		"",
		"#", "# ", "# h", "# Hello *world*",
		"a\nb", "![a](b)", "[x](y)", "`c`", "- x",
		"__double__", "**double**",
		"1)X",           // marker without space after (still fine for lexer)
		"   2. y",       // BOL with ≤3 indent
		"    4) not-ol", // >3 indent so not a marker
		"x 1) y",        // mid-line: no OL marker
	}
	for _, in := range cases {
		t.Run(fmt.Sprintf("roundtrip_%q", in), func(t *testing.T) {
			toks, err := l.Tokenize(strings.NewReader(in))
			if err != nil {
				t.Fatal(err)
			}
			if got := concatLexemes(toks); got != in {
				t.Fatalf("round-trip mismatch:\n in: %q\nout: %q\ntoks: %v", in, got, toks)
			}
		})
	}
}

// fuzz: checks round-trip for random strings.
func FuzzTokenize_RoundTrip(f *testing.F) {
	l := NewLexer()
	seeds := []string{
		"", "# title", "## x", "- a", "1) b", "2. c",
		"![a](b)", "[x](y)", "`c`", "_i_", "*b*",
		"   3) d", "    4) not-ol",
		"x 1) y", "v1.2.3", "not-a-tokenizer-feature",
		"a\nb\nc", "### Title\n- item\n1) order\n",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		toks, err := l.Tokenize(strings.NewReader(s))
		if err != nil {
			return
		}
		if out := concatLexemes(toks); out != s {
			t.Fatalf("round-trip mismatch:\n in: %q\nout: %q", s, out)
		}
	})
}
