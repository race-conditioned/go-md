package gomd

import (
	"bufio"
	"context"
	"errors"
	"io"
	"strings"
	"unicode"
)

// Tokenize lexes from the input reader and returns a slice of tokens.
func Tokenize(r io.Reader) ([]Token, error) {
	return TokenizeCtx(context.Background(), r)
}

// TokenizeCtx lexes from the input reader and returns a slice of tokens while respecting context.
func TokenizeCtx(ctx context.Context, r io.Reader) ([]Token, error) {
	var tokens []Token
	br := bufio.NewReader(r)

	line, col := 1, 0
	var buf strings.Builder

	// Track "start-of-line" and indent width (so we can allow up to 3 spaces)
	atLineStart := true
	indent := 0 // count spaces (and tabs as 2, tweak if you prefer 4)

	emitText := func() {
		if buf.Len() == 0 {
			return
		}
		tokens = append(tokens, Token{
			Kind:   TText,
			Lexeme: buf.String(),
			Pos:    Pos{Line: line, Col: col - len([]rune(buf.String()))},
		})
		buf.Reset()
	}

	// helper to read a single rune (with col tracking)
	readRune := func() (r rune, ok bool, err error) {
		if err := ctx.Err(); err != nil {
			return 0, false, err
		}
		ch, _, e := br.ReadRune()
		if e == io.EOF {
			return 0, false, e
		}
		if e != nil {
			return 0, false, e
		}
		col++
		return ch, true, nil
	}

	// unread one rune (and fix col)
	unreadRune := func() {
		_ = br.UnreadRune()
		col--
	}

	// cheap throttle: check ctx every 1024 iterations too
	iter := 0
	checkCtx := func() error {
		iter++
		if iter&0x3FF == 0 { // every 1024 loop steps
			if err := ctx.Err(); err != nil {
				return err
			}
		}
		return nil
	}

	for {
		if err := checkCtx(); err != nil {
			return tokens, err
		}

		ch, ok, err := readRune()
		if err == io.EOF {
			emitText()
			tokens = append(tokens, Token{Kind: TEOF, Pos: Pos{Line: line, Col: col}})
			break
		}
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New("not ok")
		}

		// maintain indent if still at line start and we see spaces/tabs.
		if atLineStart && (ch == ' ' || ch == '\t') {
			if ch == ' ' {
				indent += 1
			} else { // '\t'
				indent += 2 // WARN: may trip on 4 space tabs
			}
			buf.WriteRune(ch) // keep whitespace as text
			continue
		}

		// lex an ordered-list marker at BOL (after <=3 spaces).
		if atLineStart && indent <= 3 && unicode.IsDigit(ch) {
			startCol := col // col for first digit
			startLine := line
			digits := []rune{ch}

			// consume more digits
			for {
				next, ok2, _ := readRune()
				if !ok2 {
					break
				}
				if unicode.IsDigit(next) {
					digits = append(digits, next)
					continue
				}

				// not a digit: check for '.' or ')'
				if next == '.' || next == ')' {
					emitText()
					tokens = append(tokens, Token{
						Kind:   TOLMarker,
						Lexeme: string(digits) + string(next),
						Pos:    Pos{Line: startLine, Col: startCol},
					})
					atLineStart = false
					indent = 0
					goto continueMain // <-- jump to outer loop
				}

				// not a marker so unread the lookahead and treat digits as text.
				unreadRune()
				// put the digits into the text buffer
				buf.WriteString(string(digits))
				atLineStart = false
				indent = 0
				goto continueMain
			}

			// reached EOF while reading digits (handled by outer EOF on next loop)
			buf.WriteString(string(digits))
			atLineStart = false
			indent = 0
			continue
		}

		// from here, we’re not at line-start anymore.
		atLineStart = false
		indent = 0

		switch ch {
		case '#':
			emitText()
			tokens = append(tokens, Token{Kind: THash, Lexeme: "#", Pos: Pos{line, col}})
		case '*':
			emitText()
			tokens = append(tokens, Token{Kind: TStar, Lexeme: "*", Pos: Pos{line, col}})
		case '_':
			emitText()
			tokens = append(tokens, Token{Kind: TUnderscore, Lexeme: "_", Pos: Pos{line, col}})
		case '[':
			emitText()
			tokens = append(tokens, Token{Kind: TLBracket, Lexeme: "[", Pos: Pos{line, col}})
		case ']':
			emitText()
			tokens = append(tokens, Token{Kind: TRBracket, Lexeme: "]", Pos: Pos{line, col}})
		case '(':
			emitText()
			tokens = append(tokens, Token{Kind: TLParen, Lexeme: "(", Pos: Pos{line, col}})
		case ')':
			emitText()
			tokens = append(tokens, Token{Kind: TRParen, Lexeme: ")", Pos: Pos{line, col}})
		case '`':
			emitText()
			tokens = append(tokens, Token{Kind: TBacktick, Lexeme: "`", Pos: Pos{line, col}})
		case '!':
			emitText()
			tokens = append(tokens, Token{Kind: TBang, Lexeme: "!", Pos: Pos{line, col}})
		case '-':
			emitText()
			tokens = append(tokens, Token{Kind: TDash, Lexeme: "-", Pos: Pos{line, col}})
		case '\n':
			emitText()
			tokens = append(tokens, Token{Kind: TNewline, Lexeme: "\n", Pos: Pos{line, col}})
			line++
			col = 0
			atLineStart = true
			indent = 0
		default:
			buf.WriteRune(ch)
		}

	continueMain:
		// label to jump to after unread/rollback above
		_ = struct{}{}
	}

	return tokens, nil
}
