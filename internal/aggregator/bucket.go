package aggregator

import (
	"sync"
	"time"
)

// Bucket holds collapsed log entries sharing the same fingerprint.
type Bucket struct {
	mu          sync.Mutex
	Fingerprint string
	Sample      LogEntry
	Count       int
	FirstSeen   time.Time
	LastSeen    time.Time
}

// LogEntry represents a single parsed log line.
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
	Raw       string
}

// Record adds a new occurrence of a log entry to the bucket.
func (b *Bucket) Record(entry LogEntry) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.Count == 0 {
		b.Sample = entry
		b.FirstSeen = entry.Timestamp
	}
	b.Count++
	b.LastSeen = entry.Timestamp
}

// Snapshot returns a point-in-time copy of the bucket state.
func (b *Bucket) Snapshot() BucketSnapshot {
	b.mu.Lock()
	defer b.mu.Unlock()

	return BucketSnapshot{
		Fingerprint: b.Fingerprint,
		Sample:      b.Sample,
		Count:       b.Count,
		FirstSeen:   b.FirstSeen,
		LastSeen:    b.LastSeen,
	}
}

// BucketSnapshot is an immutable view of a Bucket at a point in time.
type BucketSnapshot struct {
	Fingerprint string
	Sample      LogEntry
	Count       int
	FirstSeen   time.Time
	LastSeen    time.Time
}

// Rate returns approximate events per second over the bucket's lifetime.
func (s BucketSnapshot) Rate() float64 {
	dur := s.LastSeen.Sub(s.FirstSeen).Seconds()
	if dur <= 0 {
		return float64(s.Count)
	}
	return float64(s.Count) / dur
}
