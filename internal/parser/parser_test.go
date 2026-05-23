package parser

import (
	"testing"
	"time"
)

func TestParse_ISO8601WithLevel(t *testing.T) {
	p := New()
	line := "2024-01-15T12:34:56Z INFO user logged in"
	e, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Level != LevelInfo {
		t.Errorf("expected INFO, got %s", e.Level)
	}
	if e.Message != "user logged in" {
		t.Errorf("unexpected message: %q", e.Message)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestParse_BracketedFormat(t *testing.T) {
	p := New()
	line := "[2024-01-15 12:34:56] [ERROR] database connection failed"
	e, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Level != LevelError {
		t.Errorf("expected ERROR, got %s", e.Level)
	}
	if e.Message != "database connection failed" {
		t.Errorf("unexpected message: %q", e.Message)
	}
}

func TestParse_LevelColonFormat(t *testing.T) {
	p := New()
	line := "WARN: disk usage above 90%"
	e, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Level != LevelWarn {
		t.Errorf("expected WARN, got %s", e.Level)
	}
}

func TestParse_UnknownFormat_FallsBack(t *testing.T) {
	p := New()
	line := "something completely unstructured"
	e, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Level != LevelUnknown {
		t.Errorf("expected UNKNOWN, got %s", e.Level)
	}
	if e.Message != line {
		t.Errorf("expected raw line as message")
	}
}

func TestParse_EmptyLine_ReturnsError(t *testing.T) {
	p := New()
	_, err := p.Parse("   ")
	if err == nil {
		t.Error("expected error for empty line")
	}
}

func TestParse_RawPreserved(t *testing.T) {
	p := New()
	line := "2024-06-01T00:00:00Z DEBUG hello world"
	e, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Raw != line {
		t.Errorf("raw not preserved: got %q", e.Raw)
	}
}

func TestNormalizeLevel(t *testing.T) {
	cases := []struct{ in, want Level }{
		{"debug", LevelDebug},
		{"TRACE", LevelDebug},
		{"WARNING", LevelWarn},
		{"CRITICAL", LevelFatal},
		{"ERR", LevelError},
		{"nope", LevelUnknown},
	}
	for _, c := range cases {
		got := normalizeLevel(string(c.in))
		if got != c.want {
			t.Errorf("normalizeLevel(%q) = %q, want %q", c.in, got, c.want)
		}
	}
	_ = time.Now() // keep import used
}
