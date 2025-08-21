package gomd

import (
	"context"
	"strings"
	"testing"
	"time"
)

func makeManyTokens(lines int) []Token {
	tks := make([]Token, 0, lines*3+1)
	for i := 0; i < lines; i++ {
		tks = append(tks, Token{Kind: TText, Lexeme: "x"})
		tks = append(tks, Token{Kind: TText, Lexeme: " "})
		tks = append(tks, Token{Kind: TText, Lexeme: "y"})
		tks = append(tks, Token{Kind: TNewline, Lexeme: "\n"})
	}
	tks = append(tks, Token{Kind: TEOF})
	return tks
}

func TestParseTokensCtx_Canceled_Immediate(t *testing.T) {
	l := NewLexer()
	tp := NewTokenParser()
	toks, err := l.Tokenize(strings.NewReader(strings.Repeat("a b\n", 100)))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = tp.ParseTokensCtx(ctx, toks)
	if err == nil {
		t.Fatalf("expected cancellation error, got nil")
	}
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestParseTokensCtx_Canceled_Midway(t *testing.T) {
	tp := NewTokenParser()
	toks := makeManyTokens(50000)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	_, err := tp.ParseTokensCtx(ctx, toks)
	if err == nil {
		t.Fatalf("expected cancellation error, got nil")
	}
	if err != context.Canceled && err != context.DeadlineExceeded {
		t.Fatalf("expected context cancellation, got %v", err)
	}
}

func TestParseTokensCtx_Timeout_Immediate(t *testing.T) {
	l := NewLexer()
	tp := NewTokenParser()
	toks, err := l.Tokenize(strings.NewReader("a b\na b\n"))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	_, err = tp.ParseTokensCtx(ctx, toks)
	if err != context.DeadlineExceeded {
		t.Fatalf("want DeadlineExceeded, got %v", err)
	}
}

func TestParseTokensCtx_Timeout_Midway(t *testing.T) {
	tp := NewTokenParser()
	toks := makeManyTokens(50000)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Millisecond)
	defer cancel()

	_, err := tp.ParseTokensCtx(ctx, toks)
	if err != context.DeadlineExceeded {
		t.Fatalf("want DeadlineExceeded, got %v", err)
	}
}
