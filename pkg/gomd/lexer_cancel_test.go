package gomd

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

// slowReader dribbles out 1 byte per Read with a delay, so ctx can cancel mid-way.
type slowReader struct {
	s     string
	i     int
	delay time.Duration
}

func (sr *slowReader) Read(p []byte) (int, error) {
	if sr.i >= len(sr.s) {
		return 0, io.EOF
	}
	time.Sleep(sr.delay)
	p[0] = sr.s[sr.i]
	sr.i++
	return 1, nil
}

func TestTokenizeCtx_Canceled_Immediate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before starting

	_, err := TokenizeCtx(ctx, strings.NewReader(strings.Repeat("hello\n", 1000)))
	if err == nil {
		t.Fatalf("expected cancellation error, got nil")
	}
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestTokenizeCtx_Canceled_Midway(t *testing.T) {
	// Long-ish input + slow reader so we can cancel in-flight.
	input := strings.Repeat("line with some text\n", 2000)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	_, err := TokenizeCtx(ctx, &slowReader{s: input, delay: 50 * time.Microsecond})
	if err == nil {
		t.Fatalf("expected cancellation error, got nil")
	}
	// Accept either Canceled or DeadlineExceeded depending on test timing.
	if err != context.Canceled && err != context.DeadlineExceeded {
		t.Fatalf("expected context cancellation, got %v", err)
	}
}

func TestTokenizeCtx_Timeout_Immediate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 0) // already expired
	defer cancel()
	_, err := TokenizeCtx(ctx, strings.NewReader(strings.Repeat("x\n", 1000)))
	if err != context.DeadlineExceeded {
		t.Fatalf("want DeadlineExceeded, got %v", err)
	}
}

func TestTokenizeCtx_Timeout_Midway(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Millisecond)
	defer cancel()

	input := strings.Repeat("line with some text\n", 2000)
	_, err := TokenizeCtx(ctx, &slowReader{s: input, delay: 50 * time.Microsecond})
	if err != context.DeadlineExceeded {
		t.Fatalf("want DeadlineExceeded, got %v", err)
	}
}
