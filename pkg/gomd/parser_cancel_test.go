package gomd

import (
	"context"
	"strings"
	"testing"
	"time"
)

func makeManyLines(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteString("some text with **bold** and [link](x)\n")
	}
	return b.String()
}

func TestParseCtx_Canceled_Immediate(t *testing.T) {
	p := NewOnePassParser()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := p.ParseCtx(ctx, makeManyLines(1000))
	if err == nil {
		t.Fatalf("expected cancellation error, got nil")
	}
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestParseCtx_Canceled_Midway(t *testing.T) {
	p := NewOnePassParser()
	src := makeManyLines(50000)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(3 * time.Millisecond)
		cancel()
	}()

	_, err := p.ParseCtx(ctx, src)
	if err == nil {
		t.Fatalf("expected cancellation error, got nil")
	}
	if err != context.Canceled && err != context.DeadlineExceeded {
		t.Fatalf("expected context cancellation, got %v", err)
	}
}

func TestParseCtx_Timeout_Immediate(t *testing.T) {
	p := NewOnePassParser()
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	_, err := p.ParseCtx(ctx, "a\nb\n")
	if err != context.DeadlineExceeded {
		t.Fatalf("want DeadlineExceeded, got %v", err)
	}
}

func TestParseCtx_Timeout_Midway(t *testing.T) {
	p := NewOnePassParser()
	src := makeManyLines(50000)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
	defer cancel()

	_, err := p.ParseCtx(ctx, src)
	if err != context.DeadlineExceeded {
		t.Fatalf("want DeadlineExceeded, got %v", err)
	}
}
