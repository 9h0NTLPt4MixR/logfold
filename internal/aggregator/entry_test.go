package aggregator

import (
	"testing"
	"time"
)

func makeEntry(level Level, msg string) *Entry {
	return &Entry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    map[string]string{},
		Raw:       msg,
	}
}

func TestFingerprint_SameStructure(t *testing.T) {
	e1 := makeEntry(LevelError, "connection failed after 3 retries")
	e2 := makeEntry(LevelError, "connection failed after 7 retries")

	if e1.Fingerprint() != e2.Fingerprint() {
		t.Errorf("expected same fingerprint for structurally identical messages, got %s vs %s",
			e1.Fingerprint(), e2.Fingerprint())
	}
}

func TestFingerprint_DifferentLevel(t *testing.T) {
	e1 := makeEntry(LevelInfo, "user logged in")
	e2 := makeEntry(LevelWarn, "user logged in")

	if e1.Fingerprint() == e2.Fingerprint() {
		t.Errorf("expected different fingerprints for different log levels")
	}
}

func TestFingerprint_DifferentMessage(t *testing.T) {
	e1 := makeEntry(LevelInfo, "disk usage at 80 percent")
	e2 := makeEntry(LevelInfo, "memory usage at 80 percent")

	if e1.Fingerprint() == e2.Fingerprint() {
		t.Errorf("expected different fingerprints for structurally different messages")
	}
}

func TestFingerprint_UUIDNormalized(t *testing.T) {
	e1 := makeEntry(LevelInfo, "request 550e8400-e29b-41d4-a716-446655440000 completed")
	e2 := makeEntry(LevelInfo, "request 123e4567-e89b-12d3-a456-426614174000 completed")

	if e1.Fingerprint() != e2.Fingerprint() {
		t.Errorf("expected same fingerprint when only UUID differs, got %s vs %s",
			e1.Fingerprint(), e2.Fingerprint())
	}
}

func TestNormalizeMessage(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"retried 5 times", "retried <X> times"},
		{"/var/log/app.log not found", "<X> not found"},
		{"no dynamic tokens here", "no dynamic tokens here"},
	}

	for _, tc := range cases {
		got := normalizeMessage(tc.input)
		if got != tc.expected {
			t.Errorf("normalizeMessage(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}
