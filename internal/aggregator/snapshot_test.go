package aggregator

import (
	"testing"
	"time"
)

func TestBucketSnapshot_Fields(t *testing.T) {
	s := NewStore(10*time.Second, 10)

	e := entry("info", "user 123 logged in")
	s.Ingest(e)

	snap := s.Snapshot()
	if len(snap.Buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(snap.Buckets))
	}

	b := snap.Buckets[0]
	if b.Level != "info" {
		t.Errorf("expected level 'info', got %q", b.Level)
	}
	if b.Count != 1 {
		t.Errorf("expected count 1, got %d", b.Count)
	}
	if b.Fingerprint == "" {
		t.Error("fingerprint should not be empty")
	}
	if b.Sample == "" {
		t.Error("sample should not be empty")
	}
	if b.FirstSeen.IsZero() {
		t.Error("FirstSeen should not be zero")
	}
	if b.LastSeen.IsZero() {
		t.Error("LastSeen should not be zero")
	}
}

func TestStoreSnapshot_MultipleBuckets(t *testing.T) {
	s := NewStore(10*time.Second, 10)

	s.Ingest(entry("info", "server started"))
	s.Ingest(entry("error", "connection refused"))
	s.Ingest(entry("info", "server started"))

	snap := s.Snapshot()
	if len(snap.Buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(snap.Buckets))
	}
	if snap.Taken.IsZero() {
		t.Error("Taken timestamp should not be zero")
	}
}

func TestStoreSnapshot_CountAccumulates(t *testing.T) {
	s := NewStore(10*time.Second, 10)

	for i := 0; i < 5; i++ {
		s.Ingest(entry("warn", "disk usage high"))
	}

	snap := s.Snapshot()
	if len(snap.Buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(snap.Buckets))
	}
	if snap.Buckets[0].Count != 5 {
		t.Errorf("expected count 5, got %d", snap.Buckets[0].Count)
	}
}
