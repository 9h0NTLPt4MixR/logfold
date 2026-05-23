package aggregator

import (
	"sync"
	"time"
)

// Bucket accumulates log entries that share the same fingerprint.
type Bucket struct {
	mu          sync.Mutex
	Fingerprint string
	Sample      LogEntry
	Count       int64
	FirstSeen   time.Time
	LastSeen    time.Time
	window      *Window
}

// LogEntry is the minimal log record consumed by the aggregator.
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
	Raw       string
}

// newBucket initialises a Bucket for the given fingerprint and seed entry.
func newBucket(fingerprint string, entry LogEntry) *Bucket {
	b := &Bucket{
		Fingerprint: fingerprint,
		Sample:      entry,
		Count:       1,
		FirstSeen:   entry.Timestamp,
		LastSeen:    entry.Timestamp,
		window:      NewWindow(60*time.Second, 60),
	}
	b.window.Record(entry.Timestamp)
	return b
}

// Ingest adds a new occurrence to the bucket.
func (b *Bucket) Ingest(entry LogEntry) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Count++
	b.LastSeen = entry.Timestamp
	b.window.Record(entry.Timestamp)
}

// BucketSnapshot is a point-in-time read of a bucket's state.
type BucketSnapshot struct {
	Fingerprint string
	Sample      LogEntry
	Count       int64
	FirstSeen   time.Time
	LastSeen    time.Time
	Rate        float64 // events per second over the last 60 s
}

// Snapshot returns an immutable view of the bucket at the current moment.
func (b *Bucket) Snapshot() BucketSnapshot {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	return BucketSnapshot{
		Fingerprint: b.Fingerprint,
		Sample:      b.Sample,
		Count:       b.Count,
		FirstSeen:   b.FirstSeen,
		LastSeen:    b.LastSeen,
		Rate:        b.window.Rate(now),
	}
}
