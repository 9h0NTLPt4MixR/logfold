package aggregator

import (
	"testing"
	"time"
)

func entry(level, msg string) LogEntry {
	return LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Raw:       level + " " + msg,
	}
}

func TestStore_IngestCreatesOneBucketPerFingerprint(t *testing.T) {
	s := NewStore()
	s.Ingest(entry("ERROR", "connection refused to 192.168.1.1"))
	s.Ingest(entry("ERROR", "connection refused to 10.0.0.2"))

	if s.Len() != 1 {
		t.Fatalf("expected 1 bucket, got %d", s.Len())
	}
}

func TestStore_IngestCountsCorrectly(t *testing.T) {
	s := NewStore()
	for i := 0; i < 5; i++ {
		s.Ingest(entry("WARN", "disk usage at 90%"))
	}

	snaps := s.Snapshots()
	if len(snaps) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(snaps))
	}
	if snaps[0].Count != 5 {
		t.Errorf("expected count 5, got %d", snaps[0].Count)
	}
}

func TestStore_DifferentLevelsDifferentBuckets(t *testing.T) {
	s := NewStore()
	s.Ingest(entry("INFO", "user logged in"))
	s.Ingest(entry("ERROR", "user logged in"))

	if s.Len() != 2 {
		t.Fatalf("expected 2 buckets, got %d", s.Len())
	}
}

func TestBucketSnapshot_Rate(t *testing.T) {
	now := time.Now()
	snap := BucketSnapshot{
		Count:     100,
		FirstSeen: now,
		LastSeen:  now.Add(10 * time.Second),
	}

	if snap.Rate() != 10.0 {
		t.Errorf("expected rate 10.0, got %f", snap.Rate())
	}
}

func TestBucketSnapshot_RateZeroDuration(t *testing.T) {
	now := time.Now()
	snap := BucketSnapshot{
		Count:     3,
		FirstSeen: now,
		LastSeen:  now,
	}

	if snap.Rate() != 3.0 {
		t.Errorf("expected rate 3.0 for zero duration, got %f", snap.Rate())
	}
}
