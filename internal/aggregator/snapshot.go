package aggregator

import "time"

// BucketSnapshot is a point-in-time view of a single aggregation bucket.
type BucketSnapshot struct {
	Fingerprint string
	Level       string
	Sample      string
	Count       int64
	Rate        float64
	FirstSeen   time.Time
	LastSeen    time.Time
}

// Snapshot returns a BucketSnapshot for the current state of the bucket.
func (b *Bucket) Snapshot() BucketSnapshot {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return BucketSnapshot{
		Fingerprint: b.Fingerprint,
		Level:       b.Level,
		Sample:      b.Sample,
		Count:       b.Count,
		Rate:        b.window.Rate(),
		FirstSeen:   b.FirstSeen,
		LastSeen:    b.LastSeen,
	}
}

// StoreSnapshot holds an ordered slice of BucketSnapshots from the store.
type StoreSnapshot struct {
	Taken   time.Time
	Buckets []BucketSnapshot
}

// Snapshot captures the current state of all buckets in the store.
func (s *Store) Snapshot() StoreSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snap := StoreSnapshot{
		Taken:   time.Now(),
		Buckets: make([]BucketSnapshot, 0, len(s.buckets)),
	}
	for _, b := range s.buckets {
		snap.Buckets = append(snap.Buckets, b.Snapshot())
	}
	return snap
}
